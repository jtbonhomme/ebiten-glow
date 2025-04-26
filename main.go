package main

import (
	"fmt"
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
	//
	// Note that this is a fixed function implementation of a box blur - more
	// efficiency can be gained by using a separable blur
	// (blurring horizontally and vertically separately, or for large blurs,
	// even multiple horizontal or vertical passes), ideally combined with
	// doing the summing up in a fragment shader (Kage can be used here).
	//
	// So this implementation only serves to demonstrate use of alpha blending.
	blurBox := []int{
		9, 9, 9, 9, 9, 9, 9,
		9, 7, 7, 7, 7, 7, 9,
		9, 7, 5, 5, 5, 7, 9,
		9, 7, 5, 1, 5, 7, 9,
		9, 7, 5, 5, 5, 7, 9,
		9, 7, 7, 7, 7, 7, 9,
		9, 9, 9, 9, 9, 9, 9,
	}
	for j := -3; j <= 3; j++ {
		for i := -3; i <= 3; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x+float64(i), y+float64(j))
			// This is a blur based on the source-over blend mode,
			// which is basically (GL_ONE, GL_ONE_MINUS_SRC_ALPHA). ColorM acts
			// on unpremultiplied colors, but all Ebitengine internal colors are
			// premultiplied, meaning this mode is regular alpha blending,
			// computing each destination pixel as srcPix * alpha + dstPix * (1 - alpha).
			//
			// This means that the final color is affected by the destination color when BlendSourceOver is used.
			// This blend mode is the default mode. See how this is calculated at the doc:
			// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Blend
			//
			// So if using the same alpha every time, the end result will sure be biased towards the last layer.
			//
			// Correct averaging works based on
			//   Let A_n := (a_1 + ... + a_n) / n
			//   A_{n+1} = (a_1 + ... + a_{n+1}) / (n + 1)
			//   A_{n+1} = (n * A_n + a_{n+1)) / (n + 1)
			//   A_{n+1} = A_n * (1 - 1/(n+1)) + a_{n+1} * 1/(n+1)
			// which is precisely what an alpha blend with alpha 1/(n+1) does.

			// This is a box blur, so we need to set the color scale to the inverse of the blurBox value.
			idx := (j+3)*7 + (i + 3)
			blur := blurBox[idx]
			fmt.Printf("%d %d (%d): blur: %d\n", i, j, idx, blur)
			op.ColorScale.ScaleAlpha(1 / float32(blur))
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
