package servers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
)

var Servers = make(map[string]*ServerData)

var ErrServerNotFound = errors.New("err server not found")
var ErrServerExists = errors.New("err server exists")

const PORT = 443

type ServerData struct {
	*ServerConfig
	*ServerMetric
}
type ServerConfig struct {
	Tls  string `json:"tls"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Id   string `json:"id"`
	Host string `json:"host"`
	V    string `json:"v"`
	Aid  int    `json:"aid"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Type string `json:"type"`
	Port int    `json:"port"`

	created       bool
	registered    bool
	ipv6          bool
	lastHeartBeat time.Time `json:"-"`
}

type ServerMetric struct {
	cpu       int `json:"-"`
	mem       int `json:"-"`
	tcp       int `json:"-"`
	dataQuota int `json:"-"`
	dataTotal int `json:"-"`
}

type CreateServerParams struct {
	Ip   string `josn:"ip"`
	Add  string `json:"add"`
	Host string `json:"host"`
	Ps   string `json:"ps"`
}

func CreateServer(params *CreateServerParams) error {
	if _, ok := Servers[params.Ip]; ok {
		return ErrServerExists
	}
	Servers[params.Ip] = newServer(params)
	return nil
}

func newServer(params *CreateServerParams) *ServerData {
	sc := ServerConfig{
		Tls:  "tls",
		Add:  params.Add,
		Host: params.Host,
		Ps:   params.Ps,
		// 此时通用ID似乎不合适
		Id:   "3a789def-7ed6-4df9-81c4-815252d8b79d",
		V:    "2",
		Aid:  0,
		Net:  "ws",
		Path: "/index.php",
		Type: "none",
		Port: PORT,

		lastHeartBeat: time.Now(),
		registered:    false,
		ipv6:          false,
	}
	sd := &ServerData{
		ServerConfig: &sc,
		ServerMetric: &ServerMetric{},
	}
	return sd
}

func GetInitializeParams(ip string) (*CreateServerParams, error) {
	if server, ok := Servers[ip]; ok {
		csp := CreateServerParams{}
		csp.Add = server.Add
		csp.Host = server.Host
		csp.Ps = server.Ps
		return &csp, nil
	}
	return nil, ErrServerNotFound
}

func RegisterServer(ip string) error {
	if server, ok := Servers[ip]; ok {
		server.registered = true
		logrus.Infof("server %s was registered", server.Host)
		return nil
	}
	return ErrServerNotFound
}

func RegisterHeartbeat(ip string) error {
	if server, ok := Servers[ip]; ok {
		server.lastHeartBeat = time.Now()
		return nil
	}
	logrus.Warn("non-existent ip", ip)
	return ErrServerNotFound
}

// 生成v2ray订阅格式的base64编码字符串
func GenSubscriptionData(uid string) (string, error) {
	serverKeys := make([]string, 0, len(Servers))
	for k := range Servers {
		serverKeys = append(serverKeys, k)
	}
	sort.Strings(serverKeys)
	result := ""
	for _, serverKey := range serverKeys {
		result += "vmess://"
		sv := Servers[serverKey]
		sv.Path += "?token=" + uid
		logrus.Info(sv.Path, Servers[serverKey].Path)
		b, err := json.Marshal(sv)
		if err != nil {
			return "", err
		}
		result += base64.StdEncoding.EncodeToString(b)
		result += "\n"
	}
	logrus.Infof("generating subscription for %d servers", len(Servers))
	return base64.StdEncoding.EncodeToString([]byte(result)), nil
}

func init() {
	go func() {
		for {
			for sessionid, server := range Servers {
				if time.Now().Sub(server.lastHeartBeat) > time.Second*time.Duration(config.Cfg.HeartBeatRateIntervalSec*config.Cfg.HeartBeatErrorThres) {
					logrus.Warnf("timeout on server %s", server.Host)
					delete(Servers, sessionid)
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
