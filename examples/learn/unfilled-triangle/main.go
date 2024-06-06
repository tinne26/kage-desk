package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/unfilled-triangle")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.Shader(shader)
}
