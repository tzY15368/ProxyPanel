package servers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var Servers map[string]ServerData

const PORT = 443

type ServerData struct {
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

	lastHeartBeat time.Time `json:"-"`
	cpu           int       `json:"-"`
	mem           int       `json:"-"`
	tcp           int       `json:"-"`
	dataQuota     int       `json:"-"`
	dataTotal     int       `json:"-"`
}

func newServerData(add string, host string, ps string) ServerData {
	sd := ServerData{
		Tls:  "tls",
		Add:  add,
		Host: host,
		Ps:   ps,
		// 此时通用ID似乎不合适
		Id:   "7b796c05-6552-4764-87ce-c406641a04a2",
		V:    "2",
		Aid:  0,
		Net:  "ws",
		Path: "/index.php",
		Type: "none",
		Port: PORT,

		lastHeartBeat: time.Now(),
	}

	return sd
}

func RegisterServer(add string, host string, ps string) string {
	serverKey := fmt.Sprintf("%s-%s-%d", add, host, PORT)
	logrus.Infof("registered server: %s", serverKey)
	sessionID := uuid.New().String()
	Servers[sessionID] = newServerData(add, host, ps)
	return sessionID
}

func RegisterHeartbeat(sessionID string) {
	if data, ok := Servers[sessionID]; ok {
		data.lastHeartBeat = time.Now()
		Servers[sessionID] = data
	} else {
		logrus.Warn("inexistent session id ", sessionID)
	}
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
	Servers = make(map[string]ServerData)
}
