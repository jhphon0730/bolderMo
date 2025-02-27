package game

import (
	"bolderMo-client/internal/background"
	"bolderMo-client/internal/model"
	"encoding/json"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) sendNowPosition() {
	for _, c := range g.characters {
		if c.id == g.localID {
			move := model.MoveContent{
				Direction: "now",
				Dx:        c.x,
				Dy:        c.y,
			}
			move_byte, err := json.Marshal(move)
			if err != nil {
				log.Println(err)
				return
			}

			msg := model.Message{
				Type:   model.MoveClient,
				Sender: g.localID,
				Data:   move_byte,
			}

			msgByte, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = g.conn.Write(msgByte)
			if err != nil {
				log.Println(err)
				return
			}
			break
		}
	}
}

func (g *Game) receiveMessage() {
	defer g.conn.Close()

	decoder := json.NewDecoder(g.conn)
	for {
		var msg model.Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Failed to decode message: %v", err)
			return
		}
		g.msgChan <- msg
	}
}

func (g *Game) handleServerMessage() {
	for {
		select {
		case msg := <-g.msgChan:
			switch msg.Type {
			case model.ClientDisconnected:
				removedClient := msg.Sender
				for i, char := range g.characters {
					if char.id == removedClient {
						g.characters = append(g.characters[:i], g.characters[i+1:]...)
						break
					}
				}
			case model.MoveClient:
				var move model.MoveContent
				if err := json.Unmarshal(msg.Data, &move); err != nil {
					log.Println(err)
					return
				}
				g.MoveClients(msg, move)
			case model.ClientConnected:
				charImage, err := background.LoadCharImage()
				if err != nil {
					log.Fatal(err)
				}
				charImg := ebiten.NewImageFromImage(charImage)
				char := &Character{
					id:    msg.Sender,
					x:     WINDOW_WIDTH / 2,
					y:     WINDOW_HEIGHT / 2,
					image: charImg,
				}
				g.characters = append(g.characters, char)
			case model.ClientConnectedSuccess:
				var dataStr string
				if err := json.Unmarshal(msg.Data, &dataStr); err != nil {
					log.Println(err)
					return
				}
				log.Printf("Client connected successfully: ID=%s", dataStr)
				g.localID = dataStr
				charImage, err := background.LoadCharImage()
				if err != nil {
					log.Fatal(err)
				}
				charImg := ebiten.NewImageFromImage(charImage)
				char := &Character{
					id:    dataStr,
					x:     WINDOW_WIDTH / 2,
					y:     WINDOW_HEIGHT / 2,
					image: charImg,
				}
				g.characters = append(g.characters, char)
			}
		}
	}
}
