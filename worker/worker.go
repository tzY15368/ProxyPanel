package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tzY15368/proxypanel/config"
)

// gin server for internal user auth service
var internalG *gin.Engine

var authFilterMap map[string]struct{}

func StartWorker(cfg *config.WorkerCFG) {
	authFilterMap = make(map[string]struct{})
	internalG = gin.Default()
	internalG.GET("/auth/:token", authHandler)
	go internalG.Run(fmt.Sprintf("%s:%d", cfg.WorkerHost, cfg.WorkerPort))
}

func authHandler(c *gin.Context) {
	token := c.Param("token")
	if _, ok := authFilterMap[token]; ok {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusForbidden, "")
	}
}
