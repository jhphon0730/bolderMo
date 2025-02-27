package client

import (
	"log"
	"net"
)

func ConnectServerTCP() net.Conn {
	conn, err := net.Dial("tcp", "192.168.0.5:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the server")
	return conn
}
