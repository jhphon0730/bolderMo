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
	QuitChan = make(chan struct{})
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

	messageChan <- Message{
		Type: ClientConnectedSuccess,
		Sender: client.ID,
		Data: client.ID,
	}

	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Failed to read message: %v", err)
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
					case NewChat:
						handleChat(msg.Sender, msg)
					case ClientConnected:
						handleJoin(msg.Sender, msg)
					case ClientConnectedSuccess:
						handleJoinSuccess(msg.Sender, msg)
					case ClientDisconnected:
						handleLeave(msg.Sender, msg)
					case MoveClient:
						handleMove(msg.Sender, msg)
					default:
						log.Printf("Unknown message type: %v from %s", msg.Type, msg.Sender)
					}
				case <-QuitChan:
					log.Println("Broadcast handler stopped")
					// 클라이언트들에게 종료 메시지 전송
					for _, client := range clients {
						client.Conn.Close()
					}
					return
			}
	}
}

