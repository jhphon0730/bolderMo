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

func addClient(conn net.Conn) Client {
	clientID := conn.RemoteAddr().String()
	client := &Client{
		Conn: conn,
		ID: clientID,
	}
	clientsMux.Lock()
	clients[clientID] = *client
	clientsMux.Unlock()

	messageChan <- Message{
		Type: "join",
		Sender: clientID,
	}

	return *client
}

func removeClient(clientID string) {
	clientsMux.Lock()
	delete(clients, clientID)
	clientsMux.Unlock()

	messageChan <- Message{
		Type: "leave",
		Sender: clientID,
	}
}

func HandleClient(conn net.Conn) {
	defer conn.Close()

	// 클라이언트 등록
	client := addClient(conn)
	log.Printf("Connection from %v", conn.RemoteAddr())

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

