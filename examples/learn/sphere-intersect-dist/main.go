package main

import "fmt"
import _ "embed"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/sphere-intersect-dist")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetUniformFmt("Cursor", func(value any) string {
		x := (value.([]float32))[0]
		return fmt.Sprintf("Lightness: %d%% [move cursor horizontally]", int((1.0 - x)*100))
	})
	display.Shader(shader)
}
