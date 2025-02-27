package model

import (
	"net"
)

type Client struct {
	Conn net.Conn
	ID   string
	Name string
	Dx   float64 `json:"dx"`
	Dy   float64 `json:"dy"`
}
