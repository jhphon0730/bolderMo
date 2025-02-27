package game

import (
	"bolderMo-client/internal/model"
	"encoding/json"
	"log"
)

func (g *Game) sendMoveRequest(direction string, dx, dy float64) error {
	move := model.MoveContent{
		Direction: direction,
		Dx:        dx,
		Dy:        dy,
	}
	move_byte, err := json.Marshal(move)
	if err != nil {
		log.Println(err)
		return err
	}

	// g.conn을 통해 서버로 메시비 보내기
	msg := model.Message{
		Type:   model.MoveClient,
		Sender: g.localID,
		Data:   move_byte,
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = g.conn.Write(msgByte)
	return nil
}

func (g *Game) UpdateFromServer(charID string, x, y float64) {
	for _, char := range g.characters {
		if char.id == charID {
			char.x = x
			char.y = y
			break
		}
	}
}

func (g *Game) MoveRequest(direction string, dx, dy float64) {
	err := g.sendMoveRequest(direction, dx, dy)
	if err != nil {
		return
	}

	for _, c := range g.characters {
		if c.id == g.localID {
			g.UpdateFromServer(g.localID, c.x+dx, c.y+dy)
			break
		}
	}
}

func (g *Game) MoveClients(msg model.Message, move model.MoveContent) {
	for _, c := range g.characters {
		if c.id == msg.Sender {
			g.UpdateFromServer(msg.Sender, c.x+move.Dx, c.y+move.Dy)
			break
		}
	}
}
