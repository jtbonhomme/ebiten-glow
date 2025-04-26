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
)

const (
	screenWidth  = 640
	screenHeight = 480
)

const (
	BlurDefaultIntensity float32 = 0.3
	BlurDefaultRadius    int     = 10
	BlurDefaultBase      float64 = 10.0
)

var ()

type Game struct {
	blurIntensity float32
	blurRadius    int
	blurBase      float64
	glowActive    bool
	img           []*ebiten.Image
	offscreen     *ebiten.Image
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.glowActive = !g.glowActive
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.blurIntensity += 0.01
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.blurIntensity -= 0.01
	}

	// this mapping is made from a qwerty keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.blurBase -= 0.5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.blurBase += 0.5
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.blurRadius += 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.blurRadius -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	return nil
}

func (g *Game) drawGlowImage(screen *ebiten.Image, img *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	// Draw the result on the passed coordinates.

	if g.glowActive {
		// Box blur (7x7)
		// https://en.wikipedia.org/wiki/Box_blur
		for j := -g.blurRadius; j <= g.blurRadius; j++ {
			for i := -g.blurRadius; i <= g.blurRadius; i++ {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(x+float64(i), y+float64(j))
				// This is a box blur, so we need to set the color scale to the inverse of the blurBox value.
				blur := float64(i*i+j*j) + g.blurBase
				coef := 1.0 / float32(blur)
				op.ColorScale.ScaleAlpha(coef * g.blurIntensity)
				screen.DrawImage(img, op)
			}
		}
	}

	screen.DrawImage(img, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 1})
	msg := fmt.Sprintf(`
TPS:            %0.2f
FPS:            %0.2f
blurIntensity:  %.2f    (up/down)
blurRadius:     %d      (left/right)
blurBase:       %.2f   (a/q)
glow is active: %t    (space)`,
		ebiten.ActualTPS(), ebiten.ActualFPS(), g.blurIntensity, g.blurRadius, g.blurBase, g.glowActive)

	ebitenutil.DebugPrint(screen, msg)

	for i := 0; i < len(g.img); i++ {
		g.drawGlowImage(screen, g.img[i], float64(i+1)*(screenWidth/3), screenHeight/2-3)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Blur (Ebitengine Demo)")

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

	g.blurIntensity = BlurDefaultIntensity
	g.blurRadius = BlurDefaultRadius
	g.blurBase = BlurDefaultBase
	g.offscreen = ebiten.NewImage(102, 5)
	g.glowActive = true

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
