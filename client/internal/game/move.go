package game

import (
	"bolderMo-client/internal/model"
	"encoding/json"
	"log"
)

func (g *Game) sendMoveRequest(direction string, dx, dy float64) error {
	// g.conn을 통해 서버로 메시비 보내기
	msg := model.Message{
		Type:   model.MoveClient,
		Sender: g.localID,
		Data: model.MoveContent{
			Direction: direction,
			Dx:        dx,
			Dy:        dy,
		},
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

func (g *Game) MoveLeftRequest() {
	err := g.sendMoveRequest("left", -2, 0)
	if err != nil {
		return
	}
	g.UpdateFromServer(g.localID, g.characters[0].x-2, g.characters[0].y) // 시뮬레이션
}

func (g *Game) MoveRightRequest() {
	err := g.sendMoveRequest("right", 2, 0)
	if err != nil {
		return
	}
	g.UpdateFromServer(g.localID, g.characters[0].x+2, g.characters[0].y)
}

func (g *Game) MoveUpRequest() {
	err := g.sendMoveRequest("up", 0, -2)
	if err != nil {
		return
	}
	g.UpdateFromServer(g.localID, g.characters[0].x, g.characters[0].y-2)
}

func (g *Game) MoveDownRequest() {
	err := g.sendMoveRequest("down", 0, 2)
	if err != nil {
		return
	}
	g.UpdateFromServer(g.localID, g.characters[0].x, g.characters[0].y+2)
}
