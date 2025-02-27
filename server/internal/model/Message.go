package model

type Message struct {
	Type 		string `json:"type"` // "message", "join", "leave"
	Sender 	string `json:"sender"` // Client.ID
	Data 		interface{} `json:"data"`
}
