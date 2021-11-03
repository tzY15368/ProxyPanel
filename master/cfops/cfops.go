package cfops

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
)

const TTL = 61

var api *cloudflare.API
var err error
var ctx = context.Background()
var zoneID string

func InitCFApi() error {

	api, err = cloudflare.NewWithAPIToken(config.Cfg.Master.CloudFlareAPIKey)
	if err != nil {
		return err
	}

	zoneID, err = api.ZoneIDByName(config.Cfg.Master.CloudFlareZoneName)
	return err
}

func RegisterIP(ip string) (string, error) {
	dType := "A"
	// v6
	if strings.Contains(ip, ":") {
		dType = "AAAA"
	}
	domain := hashIPToDomain(ip) + "." + config.Cfg.Master.DomainBase
	res, err := api.CreateDNSRecord(ctx, zoneID, cloudflare.DNSRecord{Type: dType, Name: domain, Content: ip, TTL: TTL})
	if err != nil {
		return "", err
	}
	//logrus.Infof("res: %v\n", res)
	logrus.WithField("success", res.Success).Info("domain creation result")
	return domain, nil
}
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func hashIPToDomain(ip string) string {
	return md5V(ip)[0:5]
}
