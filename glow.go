// ebiten-glow is a Go library built on top of the Ebiten (https://ebiten.org/) game library.
// It provides utilities for creating glowing visual effects in 2D games or graphical applications. The library simplifies the process of adding dynamic glow effects to sprites, shapes, and other graphical elements.
package ebitenglow

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BlurDefaultIntensity float32 = 0.3
	BlurDefaultRadius    int     = 10
	BlurDefaultBase      float64 = 10.0
)

// Glow represents a glow effect that can be applied to images.
type Glow struct {
	BlurIntensity float32 // The intensity of the blur effect, ranging from 0.0 to 1.0.
	BlurRadius    int     // The radius of the blur effect, which determines how far the blur extends from the center.
	BlurBase      float64 // The base value for the blur effect, which can be adjusted to change the appearance of the blur.
	Active        bool    // A boolean flag indicating whether the glow effect is active or not.
}

// New creates a new Glow instance with default values.
// It initializes the blur intensity, radius, base, and glow active state.
// The default values are:
//   - Blur intensity: 0.3
//   - Blur radius: 10
//   - Blur base: 10.0
//   - Glow active: true
func New() *Glow {
	return &Glow{
		BlurIntensity: BlurDefaultIntensity,
		BlurRadius:    BlurDefaultRadius,
		BlurBase:      BlurDefaultBase,
		Active:        true,
	}
}

// DrawImageAt draws an image at the specified coordinates (x, y) on the given screen.
// If the glow effect is active, it applies a box blur effect around the image.
// The blur effect is achieved by drawing the image multiple times at different offsets.
// The blur radius determines how far the blur extends from the center of the image.
// The blur intensity and base values control the appearance of the blur effect.
// The function uses the ebiten library to handle image drawing and transformations.
// The function takes the following parameters:
//   - screen: The ebiten image representing the screen where the image will be drawn.
//   - img: The ebiten image to be drawn.
//   - x: The x-coordinate where the image will be drawn.
//   - y: The y-coordinate where the image will be drawn.
func (g *Glow) DrawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	if g.Active {
		// Box blur (7x7)
		// https://en.wikipedia.org/wiki/Box_blur
		for j := -g.BlurRadius; j <= g.BlurRadius; j++ {
			for i := -g.BlurRadius; i <= g.BlurRadius; i++ {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(x+float64(i), y+float64(j))
				// This is a box blur, so we need to set the color scale to the inverse of the blurBox value.
				blur := float64(i*i+j*j) + g.BlurBase
				coef := 1.0 / float32(blur)
				op.ColorScale.ScaleAlpha(coef * g.BlurIntensity)
				screen.DrawImage(img, op)
			}
		}
	}

	// Draw the result on the passed coordinates.
	screen.DrawImage(img, op)
}
