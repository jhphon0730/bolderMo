package network

import (
	"log"
	"net"
	"sync"
	"bufio"
	"encoding/json"

	. "my_game_project/internal/model"
)

var (
	clients = make(map[string]Client)
	clientsMux sync.Mutex

	// Message channel
	messageChan = make(chan Message)
)

func Init() {
	go broadcastHandler()
}

func HandleClient(conn net.Conn) {
	defer conn.Close()

	// 클라이언트 등록
	client := addClient(conn)

	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)

	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Connection from %v closed", client.ID)
			removeClient(client.ID)
			return
		}

		msg.Sender = client.ID
		messageChan <- msg
	}
}

func broadcastHandler() {
	for {
			select {
			case msg := <-messageChan:
					switch msg.Type {
					case "chat":
						handleChat(msg.Sender, msg)
					case "join":
						handleJoin(msg.Sender, msg)
					case "leave":
						handleLeave(msg.Sender, msg)
					default:
						log.Printf("Unknown message type: %s from %s", msg.Type, msg.Sender)
					}
			}
	}
}

