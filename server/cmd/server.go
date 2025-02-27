package main

import (
	"fmt"
	"log"
	"net"
	"my_game_project/internal/network"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
			log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Game server running on :8080")

	go network.Init()
	for {
			conn, err := listener.Accept()
			if err != nil {
					log.Printf("Failed to accept connection: %v", err)
					continue
			}
			go network.HandleClient(conn)
	}
}
