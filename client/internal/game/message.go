package game

import (
	"bolderMo-client/internal/background"
	"bolderMo-client/internal/model"
	"encoding/json"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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
				log.Printf("Client connected successfully: ID=%s", msg.Data.(string))
				g.localID = msg.Data.(string)
				charImage, err := background.LoadCharImage()
				if err != nil {
					log.Fatal(err)
				}
				charImg := ebiten.NewImageFromImage(charImage)
				char := &Character{
					id:    msg.Data.(string),
					x:     WINDOW_WIDTH / 2,
					y:     WINDOW_HEIGHT / 2,
					image: charImg,
				}
				g.characters = append(g.characters, char)
			}
		}
	}
}
