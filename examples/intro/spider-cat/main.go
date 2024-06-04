package main

import _ "embed"

import "github.com/tinne26/kage-desk/display"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { panic(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	ebiten.SetWindowTitle("intro/spider-cat")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(384, 384)
	err = ebiten.RunGame(game)
	if err != nil { panic(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader *ebiten.Shader
	vertices [4]ebiten.Vertex
}

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	return 384, 384
}

// No logic to update.
func (self *Game) Update() error { return nil }

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

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
