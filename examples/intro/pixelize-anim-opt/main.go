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
	ebiten.SetWindowTitle("intro/pixelize-anim-opt")
	ebiten.SetWindowSize(bounds.Dx(), bounds.Dy())
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader *ebiten.Shader
	cellSize float64
	cellGrowing bool
}

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	bounds := display.ImageSpiderCatDog().Bounds()
	return bounds.Dx(), bounds.Dy()
}

// Update cell size.
func (self *Game) Update() error {
	const CellSizeChange = 8.0/60.0
	if self.cellGrowing {
		self.cellSize += CellSizeChange
		if self.cellSize >= 32.5 {
			self.cellSize = 32.5
			self.cellGrowing = false
		}
	} else { // !self.cellGrowing
		self.cellSize -= CellSizeChange
		if self.cellSize <= 1.5 {
			self.cellSize = 1.5
			self.cellGrowing = true
		}
	}

	return nil
}

// Core drawing function from where we call DrawRectShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["CellSize"] = int(self.cellSize)
	opts.Images[0] = display.ImageSpiderCatDog()
	
	// draw shader
	bounds := display.ImageSpiderCatDog().Bounds()
	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), self.shader, opts)
}
