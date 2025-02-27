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
			case model.SERVER_MESSAGE:
				var clienttBodyData map[string]*model.Client
				if err := json.Unmarshal(msg.Data, &clienttBodyData); err != nil {
					log.Println(err)
					return
				}
				if clienttBodyData == nil {
					return
				}

				g.syncMutex.Lock()
				for id, client := range clienttBodyData {
					if g.characters[id] == nil {
						continue
					}

					g.characters[id].x = client.Dx
					g.characters[id].y = client.Dy

				}
				g.syncMutex.Unlock()
			case model.ClientDisconnected:
				removedClient := msg.Sender
				for _, char := range g.characters {
					if char.id == removedClient {
						g.syncMutex.Lock()
						delete(g.characters, removedClient)
						g.syncMutex.Unlock()
						break
					}
				}
			case model.MoveClient:
				var move model.MoveContent
				if err := json.Unmarshal(msg.Data, &move); err != nil {
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
				g.syncMutex.Lock()
				g.characters[msg.Sender] = char
				g.syncMutex.Unlock()
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
				g.syncMutex.Lock()
				g.characters[msg.Sender] = char
				g.syncMutex.Unlock()
			}
		}
	}
}
