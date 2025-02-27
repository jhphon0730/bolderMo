package game

import (
	"bolderMo-client/internal/background"
	"bolderMo-client/internal/client"
	"bolderMo-client/internal/model"
	_ "image/png"
	"log"
	"net"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WINDOW_WIDTH  = 840
	WINDOW_HEIGHT = 680
	CHAR_WIDTH    = 120 // 캐릭터 크기 고정
	CHAR_HEIGHT   = 120
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
		g.MoveRequest("left", -2, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.MoveRequest("right", 2, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.MoveRequest("up", 0, -2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.MoveRequest("down", 0, 2)
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

func NewGame() *Game {
	// 1. Load background image
	bgImage, err := background.LoadBackground()
	if err != nil {
		log.Fatal(err)
	}
	bg := ebiten.NewImageFromImage(bgImage)

	// 2. Connect Server
	conn := client.ConnectServerTCP(SERVER_ADDR)

	game := &Game{
		background: bg,
		characters: []*Character{},
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
	// ebiten.SetWindowResizable(true)
	ebiten.SetWindowSizeLimits(640, 480, -1, -1)
	// ebiten.MaximizeWindow()
	ebiten.SetWindowTitle("[MOA] 공식")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
