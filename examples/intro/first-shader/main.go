package main

// This is an example shown on the "Introduction to Kage" tutorial.
// It is the first example that shows how to create a shader and
// run it with some help from kage-desk/display. Running the program
// should be as simple as:
// >> go run main.go
//
// Notice that the actual shader is on the shader.kage file, and
// it should be automatically found and loaded by kage-desk/display.

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/first-shader")
	display.SetSize(512, 512)
	display.Shader("shader.kage")
}
