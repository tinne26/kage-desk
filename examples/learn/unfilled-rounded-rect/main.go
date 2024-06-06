package main

import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/unfilled-rounded-rect")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetBackColor(display.BCOrchid)
	display.Shader(shader)
}
