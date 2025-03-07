package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/mgo"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"school/internal/config"
)

// MaxObjectID max object id
const MaxObjectID = "fffffffffffffffffffffffe"

// ErrRecordNotFound no records found
var ErrRecordNotFound = mgo.ErrNoDocuments

// InitMongodb connect mongodb
// For more information on connecting to mongodb, see https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Connect
func InitMongodb() *mgo.Database {
	dsn := utils.AdaptiveMongodbDsn(config.Get().Database.Mongodb.Dsn)
	mdb, err := mgo.Init(dsn)
	if err != nil {
		panic("mgo.Init error: " + err.Error())
	}
	return mdb
}

// ToObjectID convert string to ObjectID
func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Warnf("ToObjectID error: %v, id: %s", err, id)
	}
	return oid
}
