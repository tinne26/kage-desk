package main

import "fmt"
import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/filled-rounded-rect")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetUniformInfo("MouseButtons", "0b%02b")
	display.SetUniformFmt("Cursor", func(value any) string {
		percent := int((value.([]float32))[0]*100.0)
		var shape string
		if percent == 100 { shape = " (capsule)" }
		if percent == 0   { shape = " (pure rect)" }
		return fmt.Sprintf("Rounding: %03d%% [iff holding LMB]%s", percent, shape)
	})
	display.Shader(shader)
}
