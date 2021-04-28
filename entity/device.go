package entity

import (
	"net"
)

type Device struct {
	Id                  string
	Name                string
	Code                string
	TempotaryCodeHolder string
	Oldconnection       bool `json:"-"`
	IsAlive             bool `json:"-"`
	LastIp              string
	Conn                net.Conn `json:"-"`
}
