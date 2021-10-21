package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/master/handlers/servers"
	"github.com/tzY15368/lazarus/master/models"
)

func ServeHomeHTML(c *gin.Context) {
	c.File("html/panel.html")
	c.Abort()
}

func LoginRequired(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println(session.Get("id"))
	if session.Get("id") == nil {
		c.Status(http.StatusForbidden)
		c.Abort()
	}
	c.Next()
}
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.AbortWithStatus(http.StatusOK)
}

func LoginHandler(c *gin.Context) {
	email := c.Query("email")
	var user models.User
	err := models.DB.First(&user, "email=?", strings.ReplaceAll(email, " ", "")).Error
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	session := sessions.Default(c)
	logrus.Info("login:", user.ID)
	session.Set("id", user.ID)
	err = session.Save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.AbortWithStatus(http.StatusOK)
}

func UserInfoHandler(c *gin.Context) {
	var user models.User
	session := sessions.Default(c)
	err := models.DB.First(&user, "id=?", session.Get("id")).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, user)
}

func UpdateSubscription(c *gin.Context) {
	session := sessions.Default(c)
	var user models.User
	err := models.DB.First(&user, "id=?", session.Get("id")).Error
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
	}
	now := time.Now()
	if now.After(user.ExpireAt) {
		user.ExpireAt = now
	}
	newExpireDate := user.ExpireAt.AddDate(0, 0, 30)
	models.DB.Model(&user).Update("expire_at", newExpireDate)
	logrus.WithFields(logrus.Fields{
		"email":    user.Email,
		"expireat": user.ExpireAt,
	})
	c.AbortWithStatus(http.StatusOK)
}

func HandleSubscription(c *gin.Context) {
	token := c.Param("token")
	data, err := servers.GenSubscriptionData(token)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	c.String(http.StatusOK, data)
	c.Abort()
}

func HandleSubscriptionJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, servers.Servers)
}
