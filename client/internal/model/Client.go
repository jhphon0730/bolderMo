package model

import (
	"net"
)

type Client struct {
	Conn net.Conn `json:"-"`
	ID   string   `json:"id"`
	Name string   `json:"name"`

	Dx float64 `json:"dx"`
	Dy float64 `json:"dy"`
}
