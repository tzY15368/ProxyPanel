package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Cfg *Config

func InitConfig(path string) (*Config, error) {
	file, err := os.Open(path)
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

	logf, err := os.OpenFile(Cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	mw := io.MultiWriter(os.Stdout, logf)
	logrus.SetOutput(mw)
	level := logrus.DebugLevel
	switch Cfg.LogLevel {
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
	logrus.Info("config is read ", "logrus level is ", Cfg.LogLevel)
	return Cfg, nil
}

type Config struct {
	HeartBeatRateIntervalSec int
	HeartBeatErrorThres      int
	LogPath                  string
	LogLevel                 string
	Worker                   *WorkerCFG
	Master                   *MasterCFG
}

type WorkerCFG struct {
	Enabled     bool
	NicName     string
	WorkerHost  string
	WorkerPort  int
	MasterAddr  string
	TotalDataMB int32
}

type MasterCFG struct {
	Enabled            bool
	Db                 string
	Host               string
	Secret             string
	Port               int
	RpcHost            string
	RpcPort            int
	CloudFlareAPIKey   string
	CloudFlareZoneName string
	TelegramAPIKey     string
	TelegramGroupID    int64
}
