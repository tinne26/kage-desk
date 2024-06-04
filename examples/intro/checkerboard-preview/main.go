package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It draws a checkerboard pattern in black and white.
//
// Notice that the actual shader is on the shader.kage file.
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/checkerboard-preview@latest

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/checkerboard-preview")
	display.SetSize(256, 256)
	display.Shader(shaderProgram)
}
