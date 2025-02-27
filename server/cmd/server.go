package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	// Graceful shutdown을 위한 시그널 처리
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		<-sigChan
		fmt.Println("\nShutting down server...")
		
		// 서버 종료 시에 클라이언트들에게 종료 메시지 전송
		close(network.QuitChan)
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-network.QuitChan:
				log.Println("Server is shutting down...")
				return
			default:
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
		}
		go network.HandleClient(conn)
	}
}
