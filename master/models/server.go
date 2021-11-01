package models

import (
	"time"

	"gorm.io/gorm"
)

type Servers struct {
	gorm.Model
	Ip   string `json:"-" gorm:"uniqueIndex"`
	Tls  string `json:"tls" `
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

	Registered    bool      `json:"-"`
	Ready         bool      `json:"-"`
	HasIpv6       bool      `json:"-"`
	LastHeartBeat time.Time `json:"-"`

	Cpu       int `json:"-"`
	Mem       int `json:"-"`
	Tcp       int `json:"-"`
	DataQuota int `json:"-"`
	DataTotal int `json:"-"`
}
