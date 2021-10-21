package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB(db string) (err error) {

	DB, err = gorm.Open(sqlite.Open(db), &gorm.Config{})
	if err != nil {

		return err
	}

	DB.AutoMigrate(&User{})
	return nil
}
