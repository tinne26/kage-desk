package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It paints the left half of the screen white, and the rest black.
//
// Notice that the actual shader is on the shader.kage file.
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/half-half@latest

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/half-half")
	display.SetSize(512, 512)
	display.Shader(shaderProgram)
}
