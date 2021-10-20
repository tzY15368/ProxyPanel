package main

import (
	"fmt"
	"os"

	"github.com/tzY15368/proxypanel/config"
	"github.com/tzY15368/proxypanel/master"
	"github.com/tzY15368/proxypanel/worker"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if cfg.Master.Enabled {
		master.StartMaster(cfg.Master)
	}
	if cfg.Worker.Enabled {
		worker.StartWorker(cfg.Worker)
	}

	select {}
}
