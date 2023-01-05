package main

import "log"
import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { log.Fatal(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	ebiten.SetWindowTitle("intro/invoke-shader")
	ebiten.SetWindowSize(512, 512)
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader *ebiten.Shader
}

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	return 512, 512
}

// No logic to update.
func (self *Game) Update() error { return nil }

// Core drawing function from where we call DrawRectShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.GeoM.Translate(0, 0) // you could adjust the drawing position here
	
	// draw shader
	screen.DrawRectShader(512, 512, self.shader, opts)
}
