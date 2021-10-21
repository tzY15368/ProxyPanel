package auth

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/master/models"
	"gorm.io/gorm"
)

func TokenIsValid(token string) bool {
	var user models.User
	models.DB.Where("token = ?", token).Take(&user)
	return time.Now().Before(user.ExpireAt)
}

func GetUserMap() map[string]struct{} {
	var entries []models.User
	userMap := make(map[string]struct{})
	err := models.DB.Where("expire_at >= ?", time.Now()).Find(&entries).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("auth poll error:", err)
	} else {
		logrus.Infof("auth poll: got %d entries", len(entries))
	}
	for _, user := range entries {
		userMap[user.Token] = struct{}{}
	}
	return userMap
}
