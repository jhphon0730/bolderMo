package utils

import (
	"encoding/json"

	. "my_game_project/internal/model"
)

func MarshalMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func UnmarshalMessage(data []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}
