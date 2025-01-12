package dao

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/singleflight"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/mgo"
	"github.com/go-dev-frame/sponge/pkg/mgo/query"

	"school/internal/cache"
	"school/internal/database"
	"school/internal/model"
)

var _ TeacherDao = (*teacherDao)(nil)

// TeacherDao defining the dao interface
type TeacherDao interface {
	Create(ctx context.Context, record *model.Teacher) error
	DeleteByID(ctx context.Context, id string) error
	UpdateByID(ctx context.Context, record *model.Teacher) error
	GetByID(ctx context.Context, id string) (*model.Teacher, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Teacher, int64, error)

	DeleteByIDs(ctx context.Context, ids []string) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.Teacher, error)
	GetByIDs(ctx context.Context, ids []string) (map[string]*model.Teacher, error)
	GetByLastID(ctx context.Context, lastID string, limit int, sort string) ([]*model.Teacher, error)
}

type teacherDao struct {
	collection *mongo.Collection
	cache      cache.TeacherCache  // if nil, the cache is not used.
	sfg        *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewTeacherDao creating the dao interface
func NewTeacherDao(collection *mongo.Collection, xCache cache.TeacherCache) TeacherDao {
	if xCache == nil {
		return &teacherDao{collection: collection}
	}
	return &teacherDao{
		collection: collection,
		cache:      xCache,
		sfg:        new(singleflight.Group),
	}
}

func (d *teacherDao) deleteCache(ctx context.Context, id string) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *teacherDao) Create(ctx context.Context, record *model.Teacher) error {
	if record.ID.IsZero() {
		record.ID = primitive.NewObjectID()
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
		record.UpdatedAt = time.Now()
	}
	_, err := d.collection.InsertOne(ctx, record)

	_ = d.deleteCache(ctx, record.ID.Hex())
	return err
}

// DeleteByID soft delete a record by id
func (d *teacherDao) DeleteByID(ctx context.Context, id string) error {
	filter := bson.M{"_id": database.ToObjectID(id)}
	_, err := d.collection.UpdateOne(ctx, mgo.ExcludeDeleted(filter), mgo.EmbedDeletedAt(bson.M{}))
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *teacherDao) UpdateByID(ctx context.Context, record *model.Teacher) error {
	err := d.updateDataByID(ctx, d.collection, record)

	// delete cache
	_ = d.deleteCache(ctx, record.ID.Hex())

	return err
}

func (d *teacherDao) updateDataByID(ctx context.Context, collection *mongo.Collection, table *model.Teacher) error {
	if table.ID.IsZero() {
		return errors.New("id is empty or invalid")
	}

	update := bson.M{}

	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Age != 0 {
		update["age"] = table.Age
	}
	if table.BornAt.IsZero() == false {
		update["born_at"] = table.BornAt
	}
	if table.Job != nil {
		update["job"] = table.Job
	}
	if table.Books != nil {
		update["books"] = table.Books
	}

	filter := bson.M{"_id": table.ID}
	_, err := collection.UpdateOne(ctx, mgo.ExcludeDeleted(filter), mgo.EmbedUpdatedAt(update))
	return err
}

// GetByID get a record by id
func (d *teacherDao) GetByID(ctx context.Context, id string) (*model.Teacher, error) {
	oid := database.ToObjectID(id)
	if oid.IsZero() {
		return nil, database.ErrRecordNotFound
	}
	filter := bson.M{"_id": oid}
	// no cache
	if d.cache == nil {
		record := &model.Teacher{}
		err := d.collection.FindOne(ctx, mgo.ExcludeDeleted(filter)).Decode(record)
		return record, err
	}

	// get from cache
	cacheRecord, err := d.cache.Get(ctx, id)
	if err == nil {
		return cacheRecord, nil
	}

	// get from mongodb
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to mongodb
		val, err, _ := d.sfg.Do(id, func() (interface{}, error) {
			record := &model.Teacher{}
			err = d.collection.FindOne(ctx, mgo.ExcludeDeleted(filter)).Decode(record)
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, record, cache.TeacherExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return record, nil
		})
		if err != nil {
			return nil, err
		}
		record, ok := val.(*model.Teacher)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return record, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get paging records by column information,
// Note: query performance degrades when table rows are very large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	limit: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name, if value is of type objectId, the suffix :oid must be added, e.g. order_id:oid
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, default value is "and", support &, and, ||, or
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Limit: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *teacherDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Teacher, int64, error) {
	filter, err := params.ConvertToMongoFilter()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	total, err := d.collection.CountDocuments(ctx, mgo.ExcludeDeleted(filter))
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, total, nil
	}

	records := []*model.Teacher{}
	sort, limit, skip := params.ConvertToPage()
	findOpts := new(options.FindOptions)
	findOpts.SetLimit(int64(limit)).SetSkip(int64(skip))
	findOpts.Sort = sort

	cursor, err := d.collection.Find(ctx, mgo.ExcludeDeleted(filter), findOpts)
	if err != nil {
		return nil, 0, err
	}
	err = cursor.All(ctx, &records)
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs soft delete records by batch id
func (d *teacherDao) DeleteByIDs(ctx context.Context, ids []string) error {
	oids := mgo.ConvertToObjectIDs(ids)
	filter := bson.M{"_id": bson.M{"$in": oids}}
	_, err := d.collection.UpdateMany(ctx, mgo.ExcludeDeleted(filter), mgo.EmbedDeletedAt(bson.M{}))
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a record by condition
// query conditions:
//
//	name: column name, if value is of type objectId, the suffix :oid must be added, e.g. post_id:oid
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, default value is "and", support &, and, ||, or
//
// example: query the id of the post under the user James
//
//	condition = &query.Conditions{
//	    Columns: []query.Column{
//		{
//			Name: "user_name",
//			Value: "James",
//		},
//		{
//			Name:  "post_id:oid",
//			Value: "65ce48483f11aff697e30d6d",
//		},
//	}
func (d *teacherDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.Teacher, error) {
	filter, err := c.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	record := &model.Teacher{}
	err = d.collection.FindOne(ctx, mgo.ExcludeDeleted(filter)).Decode(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetByIDs get records by batch id
func (d *teacherDao) GetByIDs(ctx context.Context, ids []string) (map[string]*model.Teacher, error) {
	// no cache
	if d.cache == nil {
		records := []*model.Teacher{}
		oids := mgo.ConvertToObjectIDs(ids)
		filter := bson.M{"_id": bson.M{"$in": oids}}
		cursor, err := d.collection.Find(ctx, mgo.ExcludeDeleted(filter))
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, &records)
		if err != nil {
			return nil, err
		}
		itemMap := make(map[string]*model.Teacher)
		for _, record := range records {
			itemMap[record.ID.Hex()] = record
		}
		return itemMap, nil
	}

	// get form cache
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []string
	for _, id := range ids {
		if _, ok := itemMap[id]; !ok {
			missedIDs = append(missedIDs, id)
		}
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in mongodb
		var realMissedIDs []string
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if d.cache.IsPlaceholderErr(err) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		// get missed id from database
		if len(realMissedIDs) > 0 {
			records := []*model.Teacher{}
			recordIDMap := make(map[string]struct{})
			oids := mgo.ConvertToObjectIDs(realMissedIDs)
			filter := bson.M{"_id": bson.M{"$in": oids}}
			cursor, err := d.collection.Find(ctx, mgo.ExcludeDeleted(filter))
			if err != nil {
				return nil, err
			}
			err = cursor.All(ctx, &records)
			if err != nil {
				return nil, err
			}
			if len(records) > 0 {
				for _, data := range records {
					itemMap[data.ID.Hex()] = data
					recordIDMap[data.ID.Hex()] = struct{}{}
				}
				if err = d.cache.MultiSet(ctx, records, cache.TeacherExpireTime); err != nil {
					logger.Warn("cache.MultiSet error", logger.Err(err), logger.Any("ids", records))
				}
				if len(records) == len(realMissedIDs) {
					return itemMap, nil
				}
			}
			for _, id := range realMissedIDs {
				if _, ok := recordIDMap[id]; !ok {
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
				}
			}
		}
	}

	return itemMap, nil
}

// GetByLastID get paging records by last id and limit
func (d *teacherDao) GetByLastID(ctx context.Context, lastID string, limit int, sort string) ([]*model.Teacher, error) {
	page := query.NewPage(0, limit, sort)

	findOpts := new(options.FindOptions)
	findOpts.SetLimit(int64(page.Limit())).SetSkip(int64(page.Skip()))
	findOpts.Sort = page.Sort()

	records := []*model.Teacher{}
	filter := bson.M{"_id": bson.M{"$lt": database.ToObjectID(lastID)}}

	cursor, err := d.collection.Find(ctx, mgo.ExcludeDeleted(filter), findOpts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}
