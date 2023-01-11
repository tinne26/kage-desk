package main

import "log"
import _ "embed"

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
	ebiten.SetWindowTitle("intro/pixelize")
	ebiten.SetWindowSize(bounds.Dx(), bounds.Dy())
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct { shader *ebiten.Shader }

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	bounds := display.ImageSpiderCatDog().Bounds()
	return bounds.Dx(), bounds.Dy()
}

// No logic to update.
func (self *Game) Update() error { return nil }

// Core drawing function from where we call DrawRectShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = display.ImageSpiderCatDog()
	
	// draw shader
	bounds := display.ImageSpiderCatDog().Bounds()
	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), self.shader, opts)
}
