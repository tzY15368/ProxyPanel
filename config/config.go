package config

import (
	"encoding/json"
	"os"
)

var Cfg *Config

func InitConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	Cfg = &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(Cfg)

	if err != nil {
		return nil, err
	}

	return Cfg, nil
}

type Config struct {
	HeartBeatRateIntervalSec int
	HeartBeatErrorThres      int
	LogPath                  string
	Worker                   *WorkerCFG
	Master                   *MasterCFG
}

type WorkerCFG struct {
	Enabled     bool
	NicName     string
	WorkerHost  string
	WorkerPort  int
	MasterIP    string
	MasterPort  int
	Address     string
	Host        string
	PS          string
	TotalDataMB int32
}

type MasterCFG struct {
	Enabled bool
	Db      string
	Host    string
	Secret  string
	Port    int
	RpcHost string
	RpcPort int
}
