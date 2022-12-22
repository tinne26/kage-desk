package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It is the GPU version of drawing a 300x300 vertical gradient from
// green to magenta. The CPU version can be found at:
// >> github.com/tinne26/kage-desk/examples/intro/cpu/gradient
//
// Notice that the actual shader is on the shader.kage file.
// 
// The program can be run from your terminal with:
// >> go run github.com/tinne26/kage-desk/examples/intro/gpu/gradient@latest

import _ "embed"
import "github.com/tinne26/kage-desk/tools/display"

//go:embed shader.kage
var shaderSource []byte

func main() {
	display.Shader(shaderSource, "intro/gpu/gradient", 300, 300)
}
