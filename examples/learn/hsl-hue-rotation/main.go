package main

import "strconv"
import _ "embed"
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/HSL-hue-rotation")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	display.SetBackColor(display.BCGray)
	
	display.LinkUniformKey("Image", 1, ebiten.KeyDigit1, ebiten.KeyNumpad1)
	display.LinkUniformKey("Image", 2, ebiten.KeyDigit2, ebiten.KeyNumpad2)
	display.SetUniformFmt(" ", display.StrFn("Press 1 | 2 to change image"))
	display.SetUniformFmt("Cursor", func(value any) string {
		degrees := int((value.([]float32))[0]*360.0)
		return "Hue shift: " + strconv.Itoa(degrees) + " degrees"
	})

	display.Shader(shader)
}
