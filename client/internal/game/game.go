package game

import (
	"bolderMo-client/internal/background"
	"bolderMo-client/internal/client"
	"bolderMo-client/internal/model"
	"encoding/json"
	_ "image/png"
	"log"
	"net"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
	CHAR_WIDTH    = 64 // 캐릭터 크기 고정
	CHAR_HEIGHT   = 64
	SERVER_ADDR   = "192.168.0.5:8080"
)

type Character struct {
	id    string
	x, y  float64
	image *ebiten.Image
}

type Game struct {
	background *ebiten.Image
	characters []*Character
	localID    string

	conn    net.Conn
	msgChan chan model.Message
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.sendMoveRequest("left", -2, 0)
		g.UpdateFromServer(g.localID, g.characters[0].x-2, g.characters[0].y) // 시뮬레이션
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.sendMoveRequest("right", 2, 0)
		g.UpdateFromServer(g.localID, g.characters[0].x+2, g.characters[0].y)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.sendMoveRequest("up", 0, -2)
		g.UpdateFromServer(g.localID, g.characters[0].x, g.characters[0].y-2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.sendMoveRequest("down", 0, 2)
		g.UpdateFromServer(g.localID, g.characters[0].x, g.characters[0].y+2)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Scale(float64(w)/float64(g.background.Bounds().Dx()), float64(h)/float64(g.background.Bounds().Dy()))
	screen.DrawImage(g.background, bgOpts)

	for _, char := range g.characters {
		opts := &ebiten.DrawImageOptions{}
		// 캐릭터 크기 조정
		imgW := float64(char.image.Bounds().Dx()) // 원본 가로
		imgH := float64(char.image.Bounds().Dy()) // 원본 세로
		scaleX := float64(CHAR_WIDTH) / imgW
		scaleY := float64(CHAR_HEIGHT) / imgH
		opts.GeoM.Scale(scaleX, scaleY)
		// 위치에 맞게 이동 (스케일링 후 조정)
		opts.GeoM.Translate(char.x, char.y)
		screen.DrawImage(char.image, opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) sendMoveRequest(direction string, dx, dy float64) {
	log.Printf("Sending move request: ID=%s, Direction=%s, dx=%.2f, dy=%.2f", g.localID, direction, dx, dy)
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
			log.Println(msg)
		}
	}
}

func NewGame() *Game {
	// 1. Load background image
	bgImage, err := background.LoadBackground()
	if err != nil {
		log.Fatal(err)
	}
	bg := ebiten.NewImageFromImage(bgImage)

	// 2. Load character image
	charImage, err := background.LoadCharImage()
	if err != nil {
		log.Fatal(err)
	}
	charImg := ebiten.NewImageFromImage(charImage)
	char := &Character{
		id:    "player1",
		x:     WINDOW_WIDTH / 2,
		y:     WINDOW_HEIGHT / 2,
		image: charImg,
	}

	// 3. Connect Server
	conn := client.ConnectServerTCP(SERVER_ADDR)

	game := &Game{
		background: bg,
		characters: []*Character{char},
		localID:    "player1",
		conn:       conn,
		msgChan:    make(chan model.Message),
	}

	return game
}

func (g *Game) Run() {
	go g.handleServerMessage()
	go g.receiveMessage()

	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSizeLimits(640, 480, -1, -1)
	ebiten.MaximizeWindow()
	ebiten.SetWindowTitle("[MOA] 공식")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
