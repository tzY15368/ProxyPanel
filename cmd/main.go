package main

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master"
	"github.com/tzY15368/lazarus/worker"
)

func main() {
	// logrus.SetReportCaller(true)
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	TimestampFormat: "2006-01-02 15:03:04",
	// })
	path := "config.json"
	cfg, err := config.InitConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	logf, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logf.Close()
	logrus.SetOutput(logf)
	level := logrus.DebugLevel
	switch cfg.LogLevel {
	case logrus.DebugLevel.String():
		level = logrus.DebugLevel
	case logrus.InfoLevel.String():
		level = logrus.InfoLevel
	case logrus.WarnLevel.String():
		level = logrus.WarnLevel
	case logrus.ErrorLevel.String():
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
	logrus.Info("config is read", "logrus level is ", cfg.LogLevel)

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
