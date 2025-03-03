package model

import "encoding/json"

type MessageType int

const (
	ClientConnected MessageType = iota
	ClientConnectedSuccess
	ClientDisconnected
	NewChat
	MoveClient
	SERVER_MESSAGE
)

type Message struct {
	Type   MessageType     `json:"type"`   // "message", "join", "leave"
	Sender string          `json:"sender"` // Client.ID
	Data   json.RawMessage `json:"data"`
}

type MoveContent struct {
	Direction string  `json:"direction"`
	Dx        float64 `json:"dx"`
	Dy        float64 `json:"dy"`
}
