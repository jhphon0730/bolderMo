package network

import (
	"log"
	"net"

	. "my_game_project/internal/model"
)

func addClient(conn net.Conn) Client {
	log.Printf("Connection from %v", conn.RemoteAddr())

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
