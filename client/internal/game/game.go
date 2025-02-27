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
	w, h := screen.Size()
	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Scale(float64(w)/float64(g.background.Bounds().Dx()), float64(h)/float64(g.background.Bounds().Dy()))
	screen.DrawImage(g.background, bgOpts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// NewGame은 Game 객체를 초기화하고 반환만 함
func NewGame() *Game {
	bg_image, err := background.LoadBackground()
	if err != nil {
		log.Fatal(err)
	}

	bg := ebiten.NewImageFromImage(bg_image)

	return &Game{
		background: bg,
	}
}

// Run은 창 설정과 게임 실행을 담당
func (g *Game) Run() {
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSizeLimits(640, 480, -1, -1)
	ebiten.MaximizeWindow()
	ebiten.SetWindowTitle("[MOA] 공식")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
