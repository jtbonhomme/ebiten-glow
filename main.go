package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	offscreen     *ebiten.Image
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.blurIntensity += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.blurIntensity -= 0.01
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.blurBase += 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.blurBase -= 0.5
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.blurRadius += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.blurRadius -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	return nil
}

func (g *Game) drawGlowLine(screen *ebiten.Image, x, y float64) {
	// create an image line to draw on
	line := ebiten.NewImage(102, 5)
	c := color.RGBA{
		R: uint8(255),
		G: uint8(255),
		B: uint8(50),
		A: uint8(255)}
	vector.StrokeLine(line, 0, 3, 100, 3, 3, c, true)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	// Draw the result on the passed coordinates.

	/*
		// Copy the original line image to offscreen so as not to modify it.
		g.offscreen.Clear()
		g.offscreen.DrawImage(line, nil)
		blurredLine := ebiten.NewImage(102, 5)
	*/

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
			screen.DrawImage(line, op)
		}
	}

	screen.DrawImage(line, op)

	// Select and apply blending mode.
	//op.Blend = ebiten.BlendSourceOver
	//screen.DrawImage(blurredLine, op)
	/*
	   // Draw the result on the passed coordinates.
	   op = &ebiten.DrawImageOptions{}
	   op.GeoM.Translate(x, y)
	   screen.DrawImage(g.offscreen, op)
	*/
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 1})
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nblurIntensity:%.2f (up/down)\nblurRadius (left/right):%d\nblurBase:%.2f (a/q)", ebiten.ActualTPS(), ebiten.ActualFPS(), g.blurIntensity, g.blurRadius, g.blurBase)
	ebitenutil.DebugPrint(screen, msg)

	g.drawGlowLine(screen, screenWidth/2-50, screenHeight/2-3)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Blur (Ebitengine Demo)")

	g := &Game{}
	g.blurIntensity = BlurDefaultIntensity
	g.blurRadius = BlurDefaultRadius
	g.blurBase = BlurDefaultBase
	g.offscreen = ebiten.NewImage(102, 5)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
