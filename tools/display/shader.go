package display

import "os"
import "fmt"
import "log"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

// TODO: add ESC to handling.

type AutoUniform uint8
const (
	UTime AutoUniform = iota // current time in seconds. Uniform name: "Time" (float)
	UTick // current tick. Uniform name: "Tick" (float)
	UCursor // cursor position. Uniform name: "Cursor" (vec2)
	UClick // last click position or (0, 0). Uniform name: "LastClickPos" (vec2)
	UBackColor // background color. Uniform name: "BackColor" (vec4)
	USliderUnit // custom slider from 0 to 1
	USliderCent // custom slider from 0 to 100
	USliderSigned // custom slider from -1 to 1
	USlider2Unit // TODO: sliders for ints would also be cool. maybe have actual slider uniform. or a couple.
	// TODO: colormap for interpolation. with N positions?
	USlider2Cent
	USlider2Signed
	UPlayerPos // simulates a player position with gamepads or keyboard wasd/arrows
)

type changingUniformID uint8
const (
	cuRand changingUniformID = iota
	cuAngleDeg
	cuAngleRad
)

type AutoChangingUniform struct {
	key changingUniformID
	hz float32
}

func checkHz(hz float64) {
	if hz < 0 {
		panic("AutoChangingUniform doesn't accept negative hz values")
	}
}

// Hz indicate how many times we re-roll the rand per second.
// For example, 0.1 would mean re-roll once every 10 seconds.
func URand(hz float64) AutoChangingUniform {
	checkHz(hz)
	return AutoChangingUniform{ cuAngleRad, float32(hz) }
}

// Hz indicate how many loops we do per second (from 0 to 360 degrees).
func UAngleDeg(hz float64) AutoChangingUniform {
	checkHz(hz)
	return AutoChangingUniform{ cuAngleDeg, float32(hz) }
}

// Hz indicate how many loops we do per second (from 0 to 360 degrees).
func UAngleRad(hz float64) AutoChangingUniform {
	checkHz(hz)
	return AutoChangingUniform{ cuAngleRad, float32(hz) }
}

// Predefined textures that can be used with your shader.
// You can pass up to 4 to the shader options, and they will
// be assigned from image 0 to image 3, but if all the images
// need to be resizable or have the same fixed size.
//
// Fixed size images have a number at the end of their name.
// Resizable images will be resized to the canvas size.
type AutoTexture uint8
const (
	TexTriangleRed AutoTexture = iota
	TexTriangleCyan
	TexTriangleColor
	TexNoiseMono
	TexNoiseColor
	// TODO: something that has normals for it too. Like, TeapotMask, TeapotNormals

	texFixedAnchor // unexported for internal use (resizable above, fixed size below)
	TexNoiseMono32x32
	TexEbiten32x32
	TexSprite32x32
	TexAniSprite32x32
)

// Predefined background colors, to make life easier.
var (
	BCBlack   = color.RGBA{0, 0, 0, 255}
	BCWhite   = color.RGBA{255, 255, 255, 255}
	BCGray    = color.RGBA{128, 128, 128, 255}
	BCBronze  = color.RGBA{192, 128, 64, 255} // a bit of everything
	BCOrchid  = color.RGBA{128, 64, 192, 255} // a bit of everything
	BCRed     = color.RGBA{255, 0, 0, 255}
	BCGreen   = color.RGBA{0, 255, 0, 255}
	BCMagenta = color.RGBA{255, 0, 255, 255}
	BCCyan    = color.RGBA{0, 255, 255, 255}
)

// Possible options:
//  - 1 or 2 ints: window width and height. 640x480 by default.
//  - string: window title.
//  - func(*ebiten.DrawTrianglesShaderOptions): update function, at ebiten.TPS().
//  - color.RGBA: background color. Black by default.
func Shader(program []byte, opts ...any) {
	shader, err := ebiten.NewShader(program)
	if err != nil {
		// TODO: improve debug info, show context, etc.
		fmt.Printf("Failed to load shader:\n%s\n\n", err.Error())
		os.Exit(1)
	}
	
	var width, height int

	titleSet := false
	var updateFunc func(*ebiten.DrawTrianglesShaderOptions)
	for _, opt := range opts {
		switch typedOpt := opt.(type) {
		case string:
			if titleSet {
				log.Fatal("display.Shader() option error: received more than one window title")
			}
			ebiten.SetWindowTitle(typedOpt)
			titleSet = true
		case int:
			if typedOpt < 32 || typedOpt > 4096 {
				log.Fatal("display.Shader() option error: int window dimensions, when provided, must be between 32 and 4096")
			}
			if height != 0 {
				log.Fatal("display.Shader() option error: received more than two int window dimensions")
			}
			if width == 0 {
				width = typedOpt
			} else {
				height = typedOpt
			}
		case func(*ebiten.DrawTrianglesShaderOptions):
			if updateFunc != nil {
				log.Fatal("display.Shader() option error: received more than one update function")
			}
			if typedOpt == nil {
				log.Fatal("display.Shader() option error: update function can't be nil")
			}
			updateFunc = typedOpt
		default:
			log.Fatalf("display.Image(): unexpected option of type %T", opt)
		}
	}

	if !titleSet {
		ebiten.SetWindowTitle("display/shader")
	}
	if width == 0 {
		width = 640
		height = 480
	} else if height == 0 {
		height = width
	}
	ebiten.SetWindowSize(width, height)

	err = ebiten.RunGame(
		&shaderDisplayer{
			shader: shader,
			width: width,
			height: height,
			updateFunc: updateFunc,
		},
	)
	if err != nil {
		log.Fatalf("display.Shader() failure: %s", err.Error())
	}
}

type shaderDisplayer struct {
	shader *ebiten.Shader
	width, height int
	vertices [4]ebiten.Vertex
	options ebiten.DrawTrianglesShaderOptions
	updateFunc func(*ebiten.DrawTrianglesShaderOptions)
}

func (self *shaderDisplayer) Layout(_, _ int) (int, int) {
	return self.width, self.height
}

func (self *shaderDisplayer) Update() error {
	if self.updateFunc != nil { self.updateFunc(&self.options) }
	return nil
}

func (self *shaderDisplayer) Draw(screen *ebiten.Image) {
	xl, xr, yt, yb := RectToF32(screen.Bounds())
	PositionRectVertices(&self.vertices, xl, xr, yt, yb, xl, xr, yt, yb)
	indices := []uint16{0, 1, 2, 1, 2, 3}
	screen.DrawTrianglesShader(self.vertices[0 : 4], indices, self.shader, &self.options)
}
