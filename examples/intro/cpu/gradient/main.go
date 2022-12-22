package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It is the CPU version of drawing a 300x300 vertical gradient from
// green to magenta. The GPU version can be found at:
// >> github.com/tinne26/kage-desk/examples/intro/gpu/gradient
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/cpu/gradient@latest

import "image"
import "image/color"

import "github.com/tinne26/kage-desk/tools/display"

func main() {
	display.Image(Gradient(), "intro/cpu/gradient")
}

// Create the vertical gradient image.
func Gradient() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	for y := 0; y < 300; y++ {
		// compute the gradient color for this row
		floatGreen   := (300.0 - float64(y + 1))/300.0
		floatMagenta := (float64(y + 1)/300.0)

		// convert the color from float to rgba, 8-bits per channel
		green   := uint8(255*floatGreen)
		magenta := uint8(255*floatMagenta)
		clr := color.RGBA{magenta, green, magenta, 255}

		// apply the color to the whole row
		for x := 0; x < 300; x++ {
			img.SetRGBA(x, y, clr)
		}
	}	

	return img
}
