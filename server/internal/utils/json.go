package utils

import (
	"encoding/json"

	. "my_game_project/internal/model"
)

func MarshalMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}
