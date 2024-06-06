package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/simple-tint")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetBackColor(display.BCDark)
	display.SetUniformInfo("Cursor", "%vec2[0]", "X", " [move to adjust effect]")
	display.Shader(shader)
}
