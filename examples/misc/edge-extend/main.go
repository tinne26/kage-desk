package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("misc/edge-extend")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.LinkShaderImage(0, display.ImageWaterfall())
	display.SetUniformInfo("Cursor", "%vec2", "", " move to drag image")
	display.Shader(shader)
}
