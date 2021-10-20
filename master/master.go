package master

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tzY15368/proxypanel/config"
)

// gin server for handling business
var externalG *gin.Engine

func say(ctx *gin.Context) {
	ctx.String(http.StatusOK, "helo")
}

func StartMaster(cfg *config.MasterCFG) {
	externalG = gin.Default()
	externalG.GET("/", say)
	go externalG.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
}
