package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/jtbonhomme/ebitenglow"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var ()

type Game struct {
	glow *ebitenglow.Glow
	img  []*ebiten.Image
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.glow.Active = !g.glow.Active
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.glow.BlurIntensity += 0.01
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.glow.BlurIntensity -= 0.01
	}

	// this mapping is made from a qwerty keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.glow.BlurBase -= 0.5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.glow.BlurBase += 0.5
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.glow.BlurRadius += 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.glow.BlurRadius -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 1})
	msg := fmt.Sprintf(`
TPS:            %0.2f
FPS:            %0.2f
blurIntensity:  %.2f    (up/down)
blurRadius:     %d      (left/right)
blurBase:       %.2f   (a/q)
glow is active: %t    (space)
PRESS ESCAPE TO QUIT`,
		ebiten.ActualTPS(), ebiten.ActualFPS(), g.glow.BlurIntensity, g.glow.BlurRadius, g.glow.BlurBase, g.glow.Active)

	ebitenutil.DebugPrint(screen, msg)

	for i := 0; i < len(g.img); i++ {
		g.glow.DrawImageAt(screen, g.img[i], float64(i+1)*(screenWidth/3), screenHeight/2-3)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// main function initializes the game and starts the ebiten loop.
// It sets the window size, title, and creates a new Game instance.
// It also creates 2 simple images to be drawn on the screen and initializes the glow effect.
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Glow Demo (Ebitengine))")

	g := &Game{}

	g.img = make([]*ebiten.Image, 0)

	// create an image line to draw on
	g.img = append(g.img, ebiten.NewImage(102, 5))
	c0 := color.RGBA{
		R: uint8(255),
		G: uint8(255),
		B: uint8(50),
		A: uint8(255)}
	vector.StrokeLine(g.img[0], 0, 3, 100, 3, 3, c0, true)

	// create an image line to draw on
	g.img = append(g.img, ebiten.NewImage(200, 200))
	c1 := color.RGBA{
		R: uint8(50),
		G: uint8(50),
		B: uint8(255),
		A: uint8(255)}
	vector.StrokeLine(g.img[1], 3, 3, 103, 3, 3, c1, true)
	vector.StrokeLine(g.img[1], 3, 3, 3, 103, 3, c1, true)
	vector.StrokeLine(g.img[1], 103, 3, 103, 103, 3, c1, true)
	vector.StrokeLine(g.img[1], 3, 103, 103, 103, 3, c1, true)

	g.glow = ebitenglow.New()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
