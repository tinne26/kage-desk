package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It is the GPU version of drawing a 300x300 vertical gradient from
// green to blue. The CPU version can be found at:
// >> github.com/tinne26/kage-desk/examples/intro/gradient-cpu
//
// Notice that the actual shader is on the shader.kage file.
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/gradient@latest

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/gradient")
	display.SetSize(300, 300)
	display.Shader(shaderProgram)
}
