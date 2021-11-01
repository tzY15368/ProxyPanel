package worker

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master/handlers/servers"
	"github.com/tzY15368/lazarus/worker/auth"
	"github.com/tzY15368/lazarus/worker/rpc"
)

// gin server for internal user auth service
var internalG *gin.Engine

func StartWorker() {
	err := rpc.InitRPCClient()
	if err != nil {
		logrus.Fatal(err)
	}

	err = rpc.Startup()
	if err != nil {
		// service should be unreachable if not registered
		logrus.Fatal(err)
	}

	internalG = gin.Default()
	internalG.GET("/auth/:token", authHandler)
	cfg := config.Cfg.Worker
	addr := fmt.Sprintf("%s:%d", cfg.WorkerHost, cfg.WorkerPort)
	go internalG.Run(addr)
	logrus.Infof("started worker internal service on %s", addr)

	go func() {
		logrus.Infof("started heartbeat report")
		errorCounter := 0
		for {
			time.Sleep(time.Duration(config.Cfg.HeartBeatRateIntervalSec) * time.Second)
			err := rpc.SendHeartBeat()
			if errors.Is(err, servers.ErrServerNotFound) {
				logrus.Fatal("desynced with master, restarting registry")
			}
			if err != nil {
				errorCounter++
				logrus.Error(err)
			}
			if errorCounter > config.Cfg.HeartBeatErrorThres {
				logrus.Fatalf("heartbeat rpc failed failed more than %d times", config.Cfg.HeartBeatErrorThres)
			}
		}
	}()
}

func authHandler(c *gin.Context) {
	token := c.Param("token")
	v := auth.Check(token)
	if v {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}
