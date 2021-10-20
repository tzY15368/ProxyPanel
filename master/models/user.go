package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Password string
	ExpireAt time.Time
}
