package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB(db string) error {

	_DB, err := gorm.Open(sqlite.Open(db), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = _DB

	err = DB.AutoMigrate(&User{}, &Servers{})
	if err != nil {
		return err
	}
	return nil
}
