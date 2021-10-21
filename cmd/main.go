package main

import (
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master"
	"github.com/tzY15368/lazarus/worker"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
	})

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("config is read")

	if cfg.Master.Enabled {
		master.StartMaster()
	}
	if cfg.Worker.Enabled {
		// make sure worker starts after master
		time.Sleep(2 * time.Second)
		worker.StartWorker()
	}

	select {}
}
