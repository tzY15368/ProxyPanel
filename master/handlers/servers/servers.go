package servers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/master/bot"
	"github.com/tzY15368/lazarus/master/cfops"
	"github.com/tzY15368/lazarus/master/models"
	"gorm.io/gorm"
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
	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
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
		Strict: true,

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
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		tx1 := tx.Where("ip=?", ip).First(sv)
		if tx1.Error != nil {
			return tx1.Error
		}

		tx2 := tx.Model(sv).Where("ip = ?", ip).Updates(models.Servers{Registered: true, LastHeartBeat: time.Unix(0, 0)})
		return tx2.Error
	})
	if err != nil {
		return nil, err
	}
	if sv.Add == "" {
		return nil, ErrServerNotFound
	}
	logrus.WithField("host", sv.Host).Info("designated host")
	csp := &ServerConfigParmas{
		Add:  sv.Add,
		Host: sv.Host,
	}
	return csp, nil
}

func RegisterHeartbeat(ip string) error {
	sv := &models.Servers{}
	serverBecameReady := false
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("ip=?", ip).First(sv).Error
		if err != nil {
			return err
		}
		serverBecameReady = !sv.Ready
		_tx := tx.Model(sv).Where("ip=?", ip).Updates(models.Servers{Ready: true, LastHeartBeat: time.Now()})
		if _tx.Error != nil {
			return _tx.Error
		}

		logrus.WithField("rows affected", _tx.RowsAffected).WithField("ip", ip).Debug("registered heartbeat")
		return nil
	})
	if err != nil {
		return err
	}
	if serverBecameReady {
		go bot.SendMessageToGroup(fmt.Sprintf("%s(%s) 已注册", sv.Ps, sv.Host))
	}
	return nil
}

func GetValidServers() []models.Servers {
	var servers = make([]models.Servers, 0)
	tx := models.DB.Where("ready = ?", true).Find(&servers)
	if tx.Error != nil {
		logrus.Error(tx.Error)
	} else {
		logrus.WithField("serverCount", len(servers)).Info("got ready server")
	}
	v6servers := make([]models.Servers, 0)
	for _, server := range servers {
		if server.HasIpv6 {
			v6servers = append(v6servers, server)
		}
	}
	for i := range v6servers {
		v6servers[i].Host = "v6." + v6servers[i].Host
		v6servers[i].Add = "v6." + v6servers[i].Add
	}
	servers = append(servers, v6servers...)
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
			svs := []models.Servers{}
			expireTime := time.Now().Add(-1 * time.Second * time.Duration(config.Cfg.HeartBeatErrorThres) * time.Duration(config.Cfg.HeartBeatRateIntervalSec))

			// 返回error是否会是ErrRecordNotFound https://juejin.cn/post/6979915555939713054
			models.DB.Transaction(func(tx *gorm.DB) error {

				tx2 := tx.Where("last_heart_beat < ? and ready = ? and strict = ?", expireTime, true, true).Find(&svs)
				if tx2.Error != nil && !errors.Is(gorm.ErrRecordNotFound, tx2.Error) {
					logrus.WithError(tx.Error).Info("error while timeout check")
					return tx2.Error
				}
				if len(svs) != 0 {
					logrus.WithField("count", len(svs)).Warn("timeout check detected change")

					tx3 := tx.Model(sv).Where("last_heart_beat < ? and ready = ? and strict = ?", expireTime, true, true).Update("ready", false)
					logrus.WithField("rows affected", tx3.RowsAffected).Info("removed timeout servers")
					for _, _sv := range svs {
						msg += fmt.Sprintf("%s(%s) 超时\n", _sv.Ps, _sv.Host)
					}
				}
				return nil
			})
			if msg != "" {
				go bot.SendMessageToGroup(msg)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
