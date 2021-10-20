package config

import (
	"encoding/json"
	"os"
)

func InitConfig() (*Config, error) {
	cfg := &Config{}
	file, err := os.Open("config.json")
	defer file.Close()

	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type Config struct {
	HeartBeatRatePerSec int
	HeartBeatErrorThres int
	LogPath             string
	Worker              *WorkerCFG
	Master              *MasterCFG
}

type WorkerCFG struct {
	Enabled    bool
	NicName    string
	WorkerHost string
	WorkerPort int
	MasterIP   string
	MasterPort int
}

type MasterCFG struct {
	Enabled bool
	Db      string
	Host    string
	Port    int
	RpcHost string
	RpcPort int
}
