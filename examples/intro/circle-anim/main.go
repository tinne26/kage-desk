package main

import "math"
import "log"
import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// create shader object
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { log.Fatal(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	ebiten.SetWindowTitle("intro/circle-anim")
	ebiten.SetWindowSize(512, 512)
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	degrees int
	shader *ebiten.Shader
}

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	return 512, 512
}

// Increase angle degrees.
func (self *Game) Update() error {
	self.degrees += 1
	if self.degrees >= 360 { self.degrees = 0 }
	return nil
}

// Core drawing function from where we call DrawRectShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx())/2,
		float32(screen.Bounds().Dy())/2,
	}
	opts.Uniforms["Radius"] = float32(80 + 30*math.Sin(float64(self.degrees)*math.Pi/180.0))
	
	// draw shader
	screen.DrawRectShader(512, 512, self.shader, opts)
}
