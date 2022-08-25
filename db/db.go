package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"streambot/config"
	"streambot/db/models"
)

var Conn *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open(config.Config.DbPath))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(models.All...)

	Conn = db
}
