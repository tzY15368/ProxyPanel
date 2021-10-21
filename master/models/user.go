package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string    `json:"email"`
	Token    string    `json:"token";sql:"index"`
	Password string    `json:"-"`
	ExpireAt time.Time `json:"expire_at"`
}
