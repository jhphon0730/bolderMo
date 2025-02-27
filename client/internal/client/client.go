package client

import (
	"log"
	"net"
)

func ConnectServerTCP(serverAddr string) net.Conn {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the server")
	return conn
}
