package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/black-and-white")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetUniformInfo("Cursor", "%vec2[0]", "X", " [0: accurate | 1: fast]")
	display.SetUniformInfo("MouseButtons", "0b%02b")
	display.Shader(shader)
}
