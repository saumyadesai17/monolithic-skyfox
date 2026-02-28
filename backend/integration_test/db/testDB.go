package db

import (
	"skyfox/bookings/database/common"
	"skyfox/bookings/database/connection"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	"skyfox/config"

	"sync"
)

type testDB struct {
	db *common.BaseDB
}

var once sync.Once
var instance *testDB

func InitDB(cfg config.DbConfig) *testDB {
	once.Do(func() {
		handler := connection.NewDBHandler(cfg)
		db := handler.Instance()
		instance = &testDB{db: db}
	})
	return instance
}

func (s *testDB) Seed() {

	err := s.db.GormDB().AutoMigrate(model.Booking{}, model.Show{}, model.Customer{}, model.Slot{}, model.User{})
	if err != nil {
		logger.Error("error occurred while migrating schema. error: %v", err)
		return
	}
}

func GetDB() *common.BaseDB {
	return instance.db
}
