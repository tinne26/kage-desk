package display

import "os"
import "fmt"
import "math"
import "io/fs"
import "strings"
import "path/filepath"

import "github.com/hajimehoshi/ebiten/v2"

func Shader(args ...any) {
	var programBytes []byte

	//var updateFunc func(*ebiten.DrawTrianglesShaderOptions)
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case string:
			if programBytes != nil {
				fail("unexpected string argument after shader has already been loaded")
			}
			bytes, err := os.ReadFile(typedArg)
			if err != nil { fail(err.Error()) }
			programBytes = bytes
			_ = preprocess(programBytes)
		case []byte:
			if len(typedArg) == 0 {
				fail("received empty []byte shader")
			}
			if programBytes != nil {
				fail("unexpected []byte argument after shader has already been loaded")
			}
			programBytes = typedArg
		default:
			fail(fmt.Sprintf("unexpected argument of type %T on display.Image()", arg))
		}
	}

	if programBytes == nil {
		// search in current directory? or working directory? or what?
		dir, err := os.Getwd()
		if err != nil { fail(err.Error()) }
		filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil { return err }
			if entry.IsDir() && path != dir { return fs.SkipDir }
			if programBytes != nil { return nil }
			if strings.HasSuffix(entry.Name(), ".kage") {
				// found a shader, use it
				bytes, err := os.ReadFile(path)
				if err != nil { return err }
				programBytes = bytes
				_ = preprocess(programBytes)
			}
			return nil
		})
	}

	if programBytes == nil {
		fail("no shader could be found in the working directory")
	}

	shader, err := ebiten.NewShader(programBytes)
	if err != nil {
		// TODO: improve debug info, show context, etc.
		fmt.Printf("Failed to load shader:\n%s\n\n", err.Error())
		os.Exit(1)
	}

	displayer := &shaderDisplayer{
		shader: shader,
		scale: 1.0,
	}
	err = ebiten.RunGame(displayer)
	if err != nil && err != errEscClose {
		fail(err.Error())
	}
}

type shaderDisplayer struct {
	shader *ebiten.Shader
	vertices [4]ebiten.Vertex
	options ebiten.DrawTrianglesShaderOptions
	scale float64
	fsKeyPressed bool // fullscreen key
}

func (self *shaderDisplayer) Layout(w, h int) (int, int) {
	if winDisplayScaling {
		self.scale = ebiten.DeviceScaleFactor()
		if winResizable {
			return int(math.Ceil(float64(w)*self.scale)), int(math.Ceil(float64(h)*self.scale))
		} else {
			return int(math.Ceil(float64(winWidth)*self.scale)), int(math.Ceil(float64(winHeight)*self.scale))
		}
	} else {
		if winResizable {
			return w, h
		} else {
			return winWidth, winHeight
		}
	}
}

func (self *shaderDisplayer) Update() error {
	if argMaxFPS {
		ebiten.SetWindowTitle(fmt.Sprintf("%s | %.2ffps", winTitle, ebiten.ActualFPS()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errEscClose
	}
	fsKeyPressed := ebiten.IsKeyPressed(ebiten.KeyF)
	if fsKeyPressed != self.fsKeyPressed {
		if !self.fsKeyPressed {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
		self.fsKeyPressed = fsKeyPressed
	}

	//if self.updateFunc != nil { self.updateFunc(&self.options) }
	return nil
}

func (self *shaderDisplayer) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	workWidth, workHeight := winWidth, winHeight
	if winResizable {
		workWidth, workHeight = bounds.Dx(), bounds.Dy()
	}

	if winDisplayScaling {
		workWidth  = int(math.Ceil(float64(workWidth) *self.scale))
		workHeight = int(math.Ceil(float64(workHeight)*self.scale))
	}

	x := (bounds.Dx() - workWidth)/2
	y := (bounds.Dy() - workHeight)/2
	sxl, sxr := float32(x), float32(x + workWidth)
	syt, syb := float32(y), float32(y + workHeight)
	PositionRectVertices(&self.vertices, sxl, sxr, syt, syb, sxl, sxr, syt, syb)
	indices := []uint16{0, 1, 2, 1, 2, 3}
	screen.DrawTrianglesShader(self.vertices[0 : 4], indices, self.shader, &self.options)
}
