package model

import (
	"net"
)

type Client struct {
	Conn net.Conn
	ID   string
	Name string
}
