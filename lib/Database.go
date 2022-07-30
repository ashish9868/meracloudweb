package lib

import (
	"github.com/ashish9868/meracloud/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func DbInstance() *gorm.DB {
	if db == nil {
		var DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		DB.AutoMigrate(&models.User{})
		DB.AutoMigrate(&models.Upload{})

		if err != nil {
			panic("error: " + err.Error())
		}
		db = DB
	}
	return db
}

func Close() {
	if db != nil {
		sqlDB, err := db.DB()

		if err != nil {
			println("Error closing connection: " + err.Error())
		}
		sqlDB.Close()
	}
}
