package main

import _ "embed"

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
	ebiten.SetWindowTitle("intro/circle")
	ebiten.SetWindowSize(512, 512)
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
	return 512, 512
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
	// [VERTEX-NOTE]
	// Other properties will be set on later examples. The full
	// configuration is quite verbose, but you will typically create
	// your own helper functions to do the heavy lifting, and in
	// some cases you can optimize and omit some settings on
	// successive passes.

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Uniforms = make(map[string]interface{})
	shaderOpts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx())/2,
		float32(screen.Bounds().Dy())/2,
	}

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
