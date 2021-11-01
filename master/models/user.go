package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string    `json:"email"`
	Token    string    `json:"token" gorm:"uniqueIndex"`
	Password string    `json:"-"`
	ExpireAt time.Time `json:"expire_at"`
}
