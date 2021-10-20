package master

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// gin server for handling business
var externalG *gin.Engine

// db conn
var db *gorm.DB

func say(ctx *gin.Context) {
	ctx.String(http.StatusOK, "helo")
}

func StartMaster(cfg *config.MasterCFG) {
	// seteup db
	var err error
	db, err = gorm.Open(sqlite.Open(cfg.Db), &gorm.Config{})
	if err != nil {
		logrus.Fatal("db connection error:", err)
	}

	db.AutoMigrate(&models.User{})

	externalG = gin.Default()
	externalG.GET("/", say)
	go externalG.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
