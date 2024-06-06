package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/filled-circle")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetBackColor(display.BCJade)
	display.SetUniformInfo("Cursor", "%vec2", "", " [hold MLB to move circle]")
	display.SetUniformInfo("MouseButtons", "0b%02b")
	display.Shader(shader)
}
