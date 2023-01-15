package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It is the CPU version of drawing a 300x300 vertical gradient from
// green to blue. The GPU version can be found at:
// >> github.com/tinne26/kage-desk/examples/intro/gradient
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/gradient-cpu@latest

import "image"
import "image/color"

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/gradient-cpu")
	display.Image(Gradient())
}

// Create the vertical gradient image.
func Gradient() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	for y := 0; y < 300; y++ {
		// compute the gradient color for this row (amount of blue and green)
		greenLevel := (300.0 - float64(y) + 0.5)/300.0
		blueLevel  := ((float64(y) + 0.5)/300.0)

		// convert the color from float to rgba, 8-bits per channel
		greenValue := uint8(255*greenLevel)
		blueValue  := uint8(255*blueLevel)
		clr := color.RGBA{0, greenValue, blueValue, 255}

		// apply the color to the whole row
		for x := 0; x < 300; x++ {
			img.SetRGBA(x, y, clr)
		}
	}	

	return img
}
