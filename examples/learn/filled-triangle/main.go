package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/filled-triangle")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetUniformInfo("MouseButtons", "0b%02b", "", " [hold for SDF viz]")
	display.Shader(shader)
}
