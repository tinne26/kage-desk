package main

import "log"
import _ "embed"
import "image"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/kage-desk/display"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// create shader object
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { log.Fatal(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	bounds := display.ImageSpiderCatDog().Bounds()
	ebiten.SetWindowTitle("intro/mirror")
	ebiten.SetWindowSize(bounds.Dx(), bounds.Dy()*2)
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct { shader *ebiten.Shader }

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	bounds := display.ImageSpiderCatDog().Bounds()
	return bounds.Dx(), bounds.Dy()*2
}

// No logic to update.
func (self *Game) Update() error { return nil }

// Core drawing function from where we call display.DrawShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawTrianglesShaderOptions{}

	// set images for the shader
	opts.Images[0] = display.ImageSpiderCatDog()

	// set uniforms for the shader
	bounds := display.ImageSpiderCatDog().Bounds()
	rect := image.Rect(0, 0, bounds.Dx(), bounds.Dy()*2)
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["TargetRect"] = display.RectToUniform(rect)
	opts.Uniforms["MirrorAlphaMult"] = float32(0.2)
	opts.Uniforms["VertDisplacement"] = 28
	
	// draw shader
	display.DrawShader(screen, rect, self.shader, opts)
}
