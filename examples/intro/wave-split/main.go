package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It fills the left side of the screen white and the right side
// black, similarly to examples/intro/half-half, but the split
// between the sides follows a sinusoidal or wavy shape.
//
// Notice that the actual shader is on the shader.kage file.
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/wave-split@latest

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/wave-split")
	display.SetSize(512, 512)
	display.Shader(shaderProgram)
}
