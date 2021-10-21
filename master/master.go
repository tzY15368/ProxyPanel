package master

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master/handlers"
	"github.com/tzY15368/lazarus/master/models"

	// 导入session存储引擎
	"github.com/gin-contrib/sessions/cookie"
)

// gin server for handling business
var externalG *gin.Engine

func say(ctx *gin.Context) {
	ctx.String(http.StatusOK, "helo")
}

func StartMaster(cfg *config.MasterCFG) {
	// seteup db
	err := models.SetupDB(cfg.Db)
	if err != nil {
		log.Fatal("db conn error", err)
	}

	externalG = gin.Default()
	sessionStore := cookie.NewStore([]byte(cfg.Secret))

	externalG.Use(sessions.Sessions("masterSession", sessionStore))
	externalG.GET("/", say)
	externalG.GET("/login", handlers.LoginHandler)
	externalG.GET("/logout", handlers.LogoutHandler)
	externalG.GET("/update", handlers.LoginRequired, handlers.UpdateSubscription)
	externalG.GET("/user", handlers.LoginRequired, handlers.UserInfoHandler)
	externalG.GET("/s/:token", handlers.HandleSubscription)
	go externalG.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
