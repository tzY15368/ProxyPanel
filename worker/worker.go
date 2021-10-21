package worker

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/worker/rpc"
)

// gin server for internal user auth service
var internalG *gin.Engine

var authFilterMap map[string]struct{}

func StartWorker() {
	err := rpc.RegisterSelf()
	if err != nil {
		// service would be unreachable if not registered
		logrus.Fatal(err)
	}

	authFilterMap = make(map[string]struct{})
	internalG = gin.Default()
	internalG.GET("/auth/:token", authHandler)
	cfg := config.Cfg.Worker
	addr := fmt.Sprintf("%s:%d", cfg.WorkerHost, cfg.WorkerPort)
	go internalG.Run(addr)
	logrus.Infof("started worker internal service on %s", addr)
}

func authHandler(c *gin.Context) {
	token := c.Param("token")
	if _, ok := authFilterMap[token]; ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusForbidden)
	}
}
