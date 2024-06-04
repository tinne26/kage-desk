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
	ebiten.SetWindowTitle("intro/pixelize-anim")
	ebiten.SetWindowSize(bounds.Dx(), bounds.Dy())
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader *ebiten.Shader
	vertices [4]ebiten.Vertex
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
		if self.cellSize >= 32 {
			self.cellSize = 32
			self.cellGrowing = false
		}
	} else { // !self.cellGrowing
		self.cellSize -= CellSizeChange
		if self.cellSize <= 1 {
			self.cellSize = 1
			self.cellGrowing = true
		}
	}

	return nil
}

// Core drawing function from where we call DrawTrianglesShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// map the vertices to the target image
	bounds := screen.Bounds()
	self.vertices[0].DstX = float32(bounds.Min.X) // top-left
	self.vertices[0].DstY = float32(bounds.Min.Y) // top-left
	self.vertices[1].DstX = float32(bounds.Max.X) // top-right
	self.vertices[1].DstY = float32(bounds.Min.Y) // top-right
	self.vertices[2].DstX = float32(bounds.Min.X) // bottom-left
	self.vertices[2].DstY = float32(bounds.Max.Y) // bottom-left
	self.vertices[3].DstX = float32(bounds.Max.X) // bottom-right
	self.vertices[3].DstY = float32(bounds.Max.Y) // bottom-right

	// set the source image sampling coordinates
	srcBounds := display.ImageSpiderCatDog().Bounds()
	self.vertices[0].SrcX = float32(srcBounds.Min.X) // top-left
	self.vertices[0].SrcY = float32(srcBounds.Min.Y) // top-left
	self.vertices[1].SrcX = float32(srcBounds.Max.X) // top-right
	self.vertices[1].SrcY = float32(srcBounds.Min.Y) // top-right
	self.vertices[2].SrcX = float32(srcBounds.Min.X) // bottom-left
	self.vertices[2].SrcY = float32(srcBounds.Max.Y) // bottom-left
	self.vertices[3].SrcX = float32(srcBounds.Max.X) // bottom-right
	self.vertices[3].SrcY = float32(srcBounds.Max.Y) // bottom-right

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Images[0] = display.ImageSpiderCatDog()
	shaderOpts.Uniforms = make(map[string]any, 1)
	shaderOpts.Uniforms["CellSize"] = float32(int(self.cellSize))

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
