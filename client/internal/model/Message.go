package model

type MessageType int

const (
	ClientConnected MessageType = iota
	ClientDisconnected
	NewChat
)

type Message struct {
	Type 		MessageType `json:"type"` // "message", "join", "leave"
	Sender 	string `json:"sender"` // Client.ID
	Data 		interface{} `json:"data"`
}
