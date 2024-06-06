package main

import "fmt"
import _ "embed"
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shader []byte

func main() {
	display.SetTitle("learn/Oklab-chroma-shift")
	display.SetSize(512, 512, display.Resizable, display.HiRes)
	
	display.LinkUniformKey("Image", 1, ebiten.KeyDigit1, ebiten.KeyNumpad1)
	display.LinkUniformKey("Image", 2, ebiten.KeyDigit2, ebiten.KeyNumpad2)
	display.SetUniformFmt("Cursor", func(value any) string {
		x := (value.([]float32))[0]
		return fmt.Sprintf("Chroma shift: %+.02f [move cursor to change]", (x - 0.5)/2.0)
	})
	display.SetUniformFmt("_", display.StrFn("Hold LMB to compare to original"))
	display.SetUniformFmt("__", display.StrFn("Press 1 | 2 to change image"))
	display.Shader(shader)
}
