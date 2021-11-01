package servers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master/bot"
	"github.com/tzY15368/lazarus/master/cfops"
	"github.com/tzY15368/lazarus/master/models"
)

var ErrServerNotFound = errors.New("err server not found")
var ErrServerExists = errors.New("err server exists")

const PORT = 443

type ServerConfigParmas struct {
	Add  string `json:"add"`
	Host string `json:"host"`
}

type CreateServerParams struct {
	Ip   string `json:"ip"  form:"ip" binding:"required"`
	Ps   string `json:"ps"  form:"ps" binding:"required"`
	Ipv6 string `json:"ipv6"  form:"ipv6" uri:"ipv6"`
}

func CreateServer(params *CreateServerParams) error {
	sv := &models.Servers{}
	tx := models.DB.Where("ip = ?", params.Ip).First(sv)
	if tx.RowsAffected != 0 {
		return ErrServerExists
	}
	// ssh到目标机上看看有没有docker和docker-compose

	// 注册cfdns
	host, err := cfops.RegisterIP(params.Ip)
	if err != nil {
		return err
	}
	if params.Ipv6 != "" {
		_, err = cfops.RegisterIP(params.Ipv6)
	}
	if err != nil {
		return err
	}

	sv = &models.Servers{
		Ip:   params.Ip,
		Tls:  "tls",
		Add:  host,
		Host: host,
		Ps:   params.Ps,
		// 此时通用ID似乎不合适
		Id:   "3a789def-7ed6-4df9-81c4-815252d8b79d",
		V:    "2",
		Aid:  0,
		Net:  "ws",
		Path: "/index.php",
		Type: "none",
		Port: PORT,

		HasIpv6: params.Ipv6 != "",
	}
	err = models.DB.Create(sv).Error
	if err != nil {
		return err
	}

	logrus.WithField("params", *params).Info("server is created")
	return nil
}

// RegisterServer tells master that the worker is up, and provides it with config
//
// The service wont be available until the first heartbeat is received
func RegisterServer(ip string) (*ServerConfigParmas, error) {
	sv := &models.Servers{}
	tx := models.DB.Model(sv).Where("ip = ?", ip).Updates(models.Servers{Registered: true, LastHeartBeat: time.Unix(0, 0)})
	logrus.WithField("rows affected", tx.RowsAffected).WithField("ip", ip).Info("registered server")
	if tx.Error != nil {
		return nil, tx.Error
	}

	csp := &ServerConfigParmas{
		Add:  sv.Add,
		Host: sv.Host,
	}
	return csp, nil
}

func RegisterHeartbeat(ip string) error {
	sv := &models.Servers{}
	tx := models.DB.Model(sv).Where("ip=?", ip).Updates(models.Servers{Ready: true, LastHeartBeat: time.Now()})
	logrus.WithField("rows affected", tx.RowsAffected).WithField("ip", ip).Debug("registered heartbeat")
	return tx.Error
}

func GetValidServers() []models.Servers {
	var servers = make([]models.Servers, 0)
	tx := models.DB.Where("ready = ?", true).Find(&servers)
	if tx.Error != nil {
		logrus.Error(tx.Error)
	} else {
		logrus.WithField("serverCount", len(servers)).Info("got ready server")
	}
	return servers
}

// 生成v2ray订阅格式的base64编码字符串
func GenSubscriptionString(uid string) (string, error) {
	result := ""
	servers := GetValidServers()
	for _, sv := range servers {
		result += "vmess://"
		sv.Path += ("?token=" + uid)
		b, err := json.Marshal(sv)
		if err != nil {
			return "", err
		}
		result += base64.StdEncoding.EncodeToString(b)
		result += "\n"
	}
	return base64.StdEncoding.EncodeToString([]byte(result)), nil
}

func StartTimeoutCheck() {
	go func() {
		for {
			msg := ""

			sv := &models.Servers{}
			expireTime := time.Now().Add(-1 * time.Second * time.Duration(config.Cfg.HeartBeatErrorThres) * time.Duration(config.Cfg.HeartBeatRateIntervalSec))
			tx := models.DB.Model(sv).Where("last_heart_beat < ?", expireTime).Where("ready = ?", true).Update("ready", false)
			if tx.Error != nil {
				logrus.WithError(tx.Error).Info("error while timeout check")
			}
			if tx.RowsAffected != 0 {
				logrus.WithField("rows affected", tx.RowsAffected).Info("timeout check detected change")
			}
			if msg != "" {
				go func() {
					err := bot.SendMessageToGroup(msg)
					if err != nil {
						logrus.WithError(err).Warn("tgbot api send server delta failure")
					}
				}()
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
