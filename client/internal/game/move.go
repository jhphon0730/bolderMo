package game

import (
	"bolderMo-client/internal/model"
	"encoding/json"
	"log"
)

func (g *Game) sendMoveRequest(direction string, dx, dy float64) error {
	move := model.MoveContent{
		Direction: direction,
		Dx:        g.characters[g.localID].x + dx,
		Dy:        g.characters[g.localID].y + dy,
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

func (g *Game) MoveRequest(direction string, dx, dy float64) {
	err := g.sendMoveRequest(direction, dx, dy)
	if err != nil {
		return
	}
}

// 다른 사용자의 움직임을 반영
func (g *Game) MoveClients(msg model.Message, move model.MoveContent) {
	if g.characters[msg.Sender] == nil {
		return
	}

	g.characters[msg.Sender].x += move.Dx
	g.characters[msg.Sender].y += move.Dy
}
