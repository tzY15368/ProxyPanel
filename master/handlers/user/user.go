package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tzY15368/lazarus/master/models"
)

func GetCurrentUser(c *gin.Context) *models.User {
	user := &models.User{}
	session := sessions.Default(c)
	id := session.Get("id")
	err := models.DB.First(user, "id=?", id).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	return user
}
