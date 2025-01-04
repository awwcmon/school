package handler

import (
	"context"
	"errors"
	"school/internal/utils"
	"strings"
	"time"

	"github.com/jinzhu/copier"

	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/mgo/query"

	schoolV1 "school/api/school/v1"
	"school/internal/cache"
	"school/internal/dao"
	"school/internal/database"
	"school/internal/ecode"
	"school/internal/model"
)

var _ schoolV1.TeacherLogicer = (*teacherHandler)(nil)
var _ time.Time

type teacherHandler struct {
	teacherDao dao.TeacherDao
}

// NewTeacherHandler create a handler
func NewTeacherHandler() schoolV1.TeacherLogicer {
	collectionName := new(model.Teacher).TableName()
	return &teacherHandler{
		teacherDao: dao.NewTeacherDao(
			database.GetDB().Collection(collectionName),
			cache.NewTeacherCache(database.GetCacheType()),
		),
	}
}

// Create a record
func (h *teacherHandler) Create(ctx context.Context, req *schoolV1.CreateTeacherRequest) (*schoolV1.CreateTeacherReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	teacher := &model.Teacher{}
	err = copier.CopyWithOption(teacher, req, utils.CopierOption)
	if err != nil {
		return nil, ecode.ErrCreateTeacher.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = h.teacherDao.Create(ctx, teacher)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("teacher", teacher), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.CreateTeacherReply{Id: teacher.ID.Hex()}, nil
}

// DeleteByID delete a record by id
func (h *teacherHandler) DeleteByID(ctx context.Context, req *schoolV1.DeleteTeacherByIDRequest) (*schoolV1.DeleteTeacherByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	err = h.teacherDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Warn("DeleteByID error", logger.Err(err), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.DeleteTeacherByIDReply{}, nil
}

// UpdateByID update a record by id
func (h *teacherHandler) UpdateByID(ctx context.Context, req *schoolV1.UpdateTeacherByIDRequest) (*schoolV1.UpdateTeacherByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	teacher := &model.Teacher{}
	err = copier.CopyWithOption(teacher, req, utils.CopierOption)
	if err != nil {
		return nil, ecode.ErrUpdateByIDTeacher.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	teacher.ID = database.ToObjectID(req.Id)

	err = h.teacherDao.UpdateByID(ctx, teacher)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("teacher", teacher), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.UpdateTeacherByIDReply{}, nil
}

// GetByID get a record by id
func (h *teacherHandler) GetByID(ctx context.Context, req *schoolV1.GetTeacherByIDRequest) (*schoolV1.GetTeacherByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	record, err := h.teacherDao.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("id", req.Id), middleware.CtxRequestIDField(ctx))
			return nil, ecode.NotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", req.Id), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	data, err := convertTeacher(record)
	if err != nil {
		logger.Warn("convertTeacher error", logger.Err(err), logger.Any("teacher", record), middleware.CtxRequestIDField(ctx))
		return nil, ecode.ErrGetByIDTeacher.Err()
	}

	return &schoolV1.GetTeacherByIDReply{
		Teacher: data,
	}, nil
}

// List of records by query parameters
func (h *teacherHandler) List(ctx context.Context, req *schoolV1.ListTeacherRequest) (*schoolV1.ListTeacherReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.CopyWithOption(params, req.Params, utils.CopierOption)
	if err != nil {
		return nil, ecode.ErrListTeacher.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	records, total, err := h.teacherDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), middleware.CtxRequestIDField(ctx))
			return nil, ecode.InvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	teachers := []*schoolV1.Teacher{}
	for _, record := range records {
		data, err := convertTeacher(record)
		if err != nil {
			logger.Warn("convertTeacher error", logger.Err(err), logger.Any("id", record.ID), middleware.CtxRequestIDField(ctx))
			continue
		}
		teachers = append(teachers, data)
	}

	return &schoolV1.ListTeacherReply{
		Total:    total,
		Teachers: teachers,
	}, nil
}

// DeleteByIDs delete records by batch id
func (h *teacherHandler) DeleteByIDs(ctx context.Context, req *schoolV1.DeleteTeacherByIDsRequest) (*schoolV1.DeleteTeacherByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	err = h.teacherDao.DeleteByIDs(ctx, req.Ids)
	if err != nil {
		logger.Warn("DeleteByIDs error", logger.Err(err), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &schoolV1.DeleteTeacherByIDsReply{}, nil
}

// GetByCondition get a record by condition
func (h *teacherHandler) GetByCondition(ctx context.Context, req *schoolV1.GetTeacherByConditionRequest) (*schoolV1.GetTeacherByConditionReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	conditions := &query.Conditions{}
	for _, v := range req.Conditions.GetColumns() {
		column := query.Column{}
		_ = copier.CopyWithOption(&column, v, utils.CopierOption)
		conditions.Columns = append(conditions.Columns, column)
	}
	err = conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error", logger.Err(err), logger.Any("conditions", conditions), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	record, err := h.teacherDao.GetByCondition(ctx, conditions)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
			return nil, ecode.NotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	data, err := convertTeacher(record)
	if err != nil {
		logger.Warn("convertTeacher error", logger.Err(err), logger.Any("teacher", record), middleware.CtxRequestIDField(ctx))
		return nil, ecode.ErrGetByIDTeacher.Err()
	}

	return &schoolV1.GetTeacherByConditionReply{
		Teacher: data,
	}, nil
}

// ListByIDs list of records by batch id
func (h *teacherHandler) ListByIDs(ctx context.Context, req *schoolV1.ListTeacherByIDsRequest) (*schoolV1.ListTeacherByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	teacherMap, err := h.teacherDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("ids", req.Ids), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	teachers := []*schoolV1.Teacher{}
	for _, id := range req.Ids {
		if v, ok := teacherMap[id]; ok {
			record, err := convertTeacher(v)
			if err != nil {
				logger.Warn("convertTeacher error", logger.Err(err), logger.Any("teacher", v), middleware.CtxRequestIDField(ctx))
				return nil, ecode.InternalServerError.Err()
			}
			teachers = append(teachers, record)
		}
	}

	return &schoolV1.ListTeacherByIDsReply{
		Teachers: teachers,
	}, nil
}

// ListByLastID get records by last id
func (h *teacherHandler) ListByLastID(ctx context.Context, req *schoolV1.ListTeacherByLastIDRequest) (*schoolV1.ListTeacherByLastIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}
	if req.LastID == "" {
		req.LastID = database.MaxObjectID
	}

	records, err := h.teacherDao.GetByLastID(ctx, req.LastID, int(req.Limit), req.Sort)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	teachers := []*schoolV1.Teacher{}
	for _, record := range records {
		data, err := convertTeacher(record)
		if err != nil {
			logger.Warn("convertTeacher error", logger.Err(err), logger.Any("id", record.ID), middleware.CtxRequestIDField(ctx))
			continue
		}
		teachers = append(teachers, data)
	}

	return &schoolV1.ListTeacherByLastIDReply{
		Teachers: teachers,
	}, nil
}

func convertTeacher(record *model.Teacher) (*schoolV1.Teacher, error) {
	value := &schoolV1.Teacher{}
	err := copier.CopyWithOption(value, record, utils.CopierOption)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here, e.g. CreatedAt, UpdatedAt
	value.Id = record.ID.Hex()

	return value, nil
}
