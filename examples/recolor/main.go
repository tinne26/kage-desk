package main

import "os"
import "fmt"
import "log"
import _ "embed"
import "image/png"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

var darkColor  = []float32{0.1, 0.0, 0.0, 1.0}
var lightColor = []float32{0.2, 1.0, 1.0, 1.0}
var image *ebiten.Image
var horzResolution, vertResolution int

func main() {
	// load image from arguments
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run main.go path/to/source.png\n")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil { log.Fatal(err) }
	defer file.Close() // ignoring close error
	img, err := png.Decode(file)
	if err != nil { log.Fatal(err) }
	image = ebiten.NewImageFromImage(img)
	bounds := image.Bounds()
	horzResolution, vertResolution = bounds.Dx(), bounds.Dy()

	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { log.Fatal(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	ebiten.SetWindowTitle("examples/recolor")
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
	return horzResolution, vertResolution
}

// No logic to update.
func (self *Game) Update() error { return nil }

// Core drawing function from where we call DrawRectShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["DarkColor"]  = darkColor
	opts.Uniforms["LightColor"] = lightColor
	opts.Images[0] = image
	
	// draw shader
	screen.DrawRectShader(horzResolution, vertResolution, self.shader, opts)
}
