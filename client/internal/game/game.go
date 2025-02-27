package game

import (
	"bolderMo-client/internal/background"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
)

type Game struct {
	background *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func NewGame() Game {
	bg_image, err := background.LoadBackground()
	if err != nil {
		log.Fatal(err)
	}

	bg := ebiten.NewImageFromImage(bg_image)

	gameEngine := &Game{
		background: bg,
	}

	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	ebiten.SetWindowTitle("[MOA] 공식")
	if err := ebiten.RunGame(gameEngine); err != nil {
		log.Fatal(err)
	}

	return *gameEngine
}
