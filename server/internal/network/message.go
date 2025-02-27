package network

import (
	"log"

	. "my_game_project/internal/model"
	"my_game_project/internal/utils"
)

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
			msgStr, err := utils.MarshalMessage(msg)
			if err != nil {
					log.Printf("Failed to marshal message: %v", err)
					continue
			}
			if _, err := client.Conn.Write(append(msgStr, '\n')); err != nil {
					log.Printf("Failed to send message to %v: %v", client.ID, err)
			}
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
			msgStr, err := utils.MarshalMessage(msg)
			if err != nil {
					log.Printf("Failed to marshal message: %v", err)
					continue
			}
			if _, err := client.Conn.Write(append(msgStr, '\n')); err != nil {
					log.Printf("Failed to send message to %v: %v", client.ID, err)
			}
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
			msgStr, err := utils.MarshalMessage(msg)
			if err != nil {
					log.Printf("Failed to marshal message: %v", err)
					continue
			}
			if _, err := client.Conn.Write(append(msgStr, '\n')); err != nil {
					log.Printf("Failed to send message to %v: %v", client.ID, err)
			}
	}
}
