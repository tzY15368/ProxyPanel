package main

import (
	"fmt"
	"os"

	"github.com/tzY15368/proxypanel/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if cfg.Master.Enabled {

	}
	if cfg.Worker.Enabled {

	}
}
