package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	offscreen *ebiten.Image
}

func (g *Game) Update() error {
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
	vector.StrokeLine(line, 1, 3, 101, 3, 3, c, true)

	line2 := ebiten.NewImage(102, 5)
	c2 := color.RGBA{
		R: uint8(0),
		G: uint8(0),
		B: uint8(255),
		A: uint8(255)}
	vector.StrokeLine(line2, 1, 3, 101, 3, 3, c2, true)

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
	blurBox := []int{
		13, 13, 13, 13, 13, 13, 13,
		13, 17, 17, 17, 17, 17, 13,
		13, 17, 21, 21, 21, 17, 13,
		13, 17, 21, 25, 21, 17, 13,
		13, 17, 21, 21, 21, 17, 13,
		13, 17, 17, 17, 17, 17, 13,
		13, 13, 13, 13, 13, 13, 13,
	}
	for j := -3; j <= 3; j++ {
		for i := -3; i <= 3; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x+float64(i), y+float64(j))
			// This is a box blur, so we need to set the color scale to the inverse of the blurBox value.
			idx := (j+3)*7 + (i + 3)
			blur := blurBox[idx]
			coef := 1.0 / float32(blur)
			op.ColorScale.ScaleAlpha(coef * 0.5)
			screen.DrawImage(line, op)
		}
	}

	screen.DrawImage(line2, op)

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

	g.drawGlowLine(screen, screenWidth/2-50, screenHeight/2-3)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Blur (Ebitengine Demo)")

	g := &Game{}
	g.offscreen = ebiten.NewImage(102, 5)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
