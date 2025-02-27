package model

type MessageType int

const (
	ClientConnected MessageType = iota
	ClientConnectedSuccess
	ClientDisconnected
	NewChat
	MoveClient
)

type Message struct {
	Type 		MessageType `json:"type"` // "message", "join", "leave"
	Sender 	string `json:"sender"` // Client.ID
	Data 		interface{} `json:"data"`
}
