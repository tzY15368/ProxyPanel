package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/master/handlers/servers"
	"github.com/tzY15368/lazarus/master/handlers/user"
	"github.com/tzY15368/lazarus/master/models"
)

func ServeHomeHTML(c *gin.Context) {
	c.File("html/panel.html")
	c.Abort()
}

func LoginRequired(c *gin.Context) {
	session := sessions.Default(c)
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
	u := user.GetCurrentUser(c)
	c.AbortWithStatusJSON(http.StatusOK, u)
}

func UpdateSubscription(c *gin.Context) {
	u := user.GetCurrentUser(c)
	now := time.Now()
	if now.After(u.ExpireAt) {
		u.ExpireAt = now
	}
	newExpireDate := u.ExpireAt.AddDate(0, 0, 30)
	models.DB.Model(&u).Update("expire_at", newExpireDate)
	logrus.WithFields(logrus.Fields{
		"email":    u.Email,
		"expireat": u.ExpireAt,
	})
	c.AbortWithStatus(http.StatusOK)
}

func HandleSubscription(c *gin.Context) {
	token := c.Param("token")
	data, err := servers.GenSubscriptionString(token)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	c.String(http.StatusOK, data)
	c.Abort()
}

func HandleSubscriptionJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, servers.GetValidServers())
}

func HandleTokenRefresh(c *gin.Context) {
	u := user.GetCurrentUser(c)
	err := models.DB.Model(&u).Update("token", uuid.New().String()).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
