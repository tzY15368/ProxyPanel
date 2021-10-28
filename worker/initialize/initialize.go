package initialize

import "errors"

var NginxMissing = errors.New("nginx binary is missing")

var V2RayMissing = errors.New("v2ray binary is missing")

func InitializeComponents(add string, host string, port int) error {
	return nil
}
