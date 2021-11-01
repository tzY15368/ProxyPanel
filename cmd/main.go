package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master"
	"github.com/tzY15368/lazarus/worker"
)

func main() {
	path := "config.json"
	cfg, err := config.InitConfig(path)
	if err != nil {
		logrus.Fatal(err)
	}

	if len(os.Args) > 1 {
		logrus.Info("has more than one args, only starting the worker anyways, overriding worker config")
		masterAddr := os.Args[1]
		cfg.Worker.Enabled = true
		cfg.Worker.MasterAddr = masterAddr
		cfg.Worker.WorkerHost = "127.0.0.1"
		cfg.Worker.WorkerPort = 1239
		logrus.Info("worker: got master addr", masterAddr)
		worker.StartWorker()
	} else {

		if cfg.Master.Enabled {
			master.StartMaster()
		}
		if cfg.Worker.Enabled {
			// make sure worker starts after master
			time.Sleep(2 * time.Second)
			worker.StartWorker()
		}
	}

	select {}
}
