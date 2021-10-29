package cfops

import (
	"context"
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

func RegisterIP(ip string) (domain string, err error) {
	dType := "A"
	// v6
	if strings.Contains(ip, ":") {
		dType = "AAAA"
	}
	res, err := api.CreateDNSRecord(ctx, zoneID, cloudflare.DNSRecord{Type: dType, Name: "test.fmagic.icu", Content: ip, TTL: TTL})
	if err != nil {
		return
	}
	logrus.Infof("res: %v\n", res)
	return
}
