package network

import (
	"log"
	"encoding/json"

	. "my_game_project/internal/model"
	"my_game_project/internal/utils"
)

func handleSend(client Client, msg Message) {
	msgStr, err := utils.MarshalMessage(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
	}
	if _, err := client.Conn.Write(append(msgStr, '\n')); err != nil {
		log.Printf("Failed to send message to %v: %v", client.ID, err)
	}
}

func handleClients() {
	if clients == nil {
		return
	}

	msg := Message{
		Type: SERVER_MESSAGE,
		Data: clients,
		Sender: "server",
	}

	msgStr, err := utils.MarshalMessage(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
	}

	for _, client := range clients {
		if _, err := client.Conn.Write(append(msgStr, '\n')); err != nil {
			log.Printf("Failed to send message to %v: %v", client.ID, err)
		}
	}

}

func handleChat(clientID string, msg Message) {
	clientsMux.Lock()
	clientList := make(map[string]Client, len(clients))
	for id, c := range clients {
		clientList[id] = c
	}
	clientsMux.Unlock()

	for _, client := range clientList {
		if client.ID == clientID { // 발신자 제외
			continue
		}
		handleSend(client, msg)
	}
}

func handleJoin(clientID string, msg Message) {
	log.Printf("Client %v joined", clientID)

	clientsMux.Lock()
	clientList := make(map[string]Client, len(clients))
	for id, c := range clients {
			clientList[id] = c
	}
	clientsMux.Unlock()

	for _, client := range clientList {
		if client.ID == clientID {
			continue
		}
		handleSend(client, msg)
	}
}

func handleLeave(clientID string, msg Message) {
	log.Printf("Client %v left", clientID)

	clientsMux.Lock()
	clientList := make(map[string]Client, len(clients))
	for id, c := range clients {
			clientList[id] = c
	}
	clientsMux.Unlock()

	for _, client := range clientList {
		if client.ID == clientID {
			continue
		}
		handleSend(client, msg)
	}
}

func handleMove(clientID string, msg Message) {
	// catch the only move message
	msgByte, err := utils.MarshalMessage(msg)
	if err != nil {
		return
	}
	var moveMessage MoveMessage
	if err := json.Unmarshal(msgByte, &moveMessage); err != nil {
		return
	}


	clientsMux.Lock()
	clientList := make(map[string]Client, len(clients))
	for id, c := range clients {
			clientList[id] = c
	}
	clientsMux.Unlock()

	for _, client := range clientList {
		if client.ID == clientID {
			find_client := clients[clientID]
			find_client.Dx = moveMessage.Data.Dx
			find_client.Dy = moveMessage.Data.Dy
			clients[clientID] = find_client
			continue
		}
		handleSend(client, msg)
	}
}

func handleJoinSuccess(clientID string, msg Message) {
	log.Printf("Client %v joined successfully", clientID)

	clientsMux.Lock()
	clientList := make(map[string]Client, len(clients))
	for id, c := range clients {
			clientList[id] = c
	}
	clientsMux.Unlock()

	for _, client := range clientList {
		if client.ID != clientID {
			continue
		}
		handleSend(client, msg)
	}
}
