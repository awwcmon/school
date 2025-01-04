package database

import (
	"strings"
	"sync"

	"github.com/zhufuyi/sponge/pkg/mgo"

	"school/internal/config"
)

var (
	mdb     *mgo.Database
	mdbOnce sync.Once
)

// InitDB connect database
func InitDB() {
	dbDriver := config.Get().Database.Driver
	switch strings.ToLower(dbDriver) {
	case mgo.DBDriverName:
		mdb = InitMongodb()
	default:
		panic("InitDB error, please modify the correct 'database' configuration at yaml file. " +
			"Refer to https://school/blob/main/configs/school.yml#L85")
	}
}

// GetDB get db
func GetDB() *mgo.Database {
	if mdb == nil {
		mdbOnce.Do(func() {
			InitDB()
		})
	}

	return mdb
}

// CloseDB close db
func CloseDB() error {
	return mgo.Close(mdb)
}
