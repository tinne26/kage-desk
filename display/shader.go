package display

import "os"
import "fmt"
import "math"
import "sort"
import "time"
import "image"
import "io/fs"
import "strings"
import "path/filepath"
import "unicode/utf8"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/ebitenutil"
import "github.com/hajimehoshi/ebiten/v2/inpututil"

// Loads and executes a shader. The shader might be explicitly
// given as a string or byte slice; otherwise, the method will
// try to search on the current directory for a relevant .kage
// file.
//
// The shaders have a few functionalities built-in:
//  - Multiple uniforms are given by default. This includes
//    'Time float', 'Cursor vec2' (in [0, 1] normalized
//    coordinates), 'MouseButtons int' (0b00 if none, 0b10
//    if left, 0b01 if right, 0b11 if both).
//  - Sample textures are linked to Images[0] and Images[1] if
//    image usage is detected. You can also [LinkShaderImage]()
//    on your own.
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
			// _ = preprocess(programBytes)
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
				// _ = preprocess(programBytes)
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
		startTime: time.Now(),
		usingImage0: containsOutsideComment(programBytes, []rune("imageSrc0")),
		usingImage1: containsOutsideComment(programBytes, []rune("imageSrc1")),
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
	usingImage0 bool
	usingImage1 bool
	startTime time.Time
}

func (self *shaderDisplayer) Layout(w, h int) (int, int) {
	if winDisplayScaling {
		scale := ebiten.DeviceScaleFactor()
		aspectRatio := float64(winWidth)/float64(winHeight)
		var w64, h64 float64 = float64(w), float64(h)
		proportionalWidth  := h64*aspectRatio
		proportionalHeight := w64/aspectRatio
		switch {
		case proportionalWidth  > w64:
			h64 = w64/aspectRatio
		case proportionalHeight > h64:
			w64 = h64*aspectRatio
		default:
			if w64/h64 != 1.0 { panic("unexpected") }
		}
		return int(math.Ceil(w64*scale)), int(math.Ceil(h64*scale))
	} else {
		return winWidth, winHeight
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

	// update key detection
	for _, kvu := range keyValueUniforms {
		for _, key := range kvu.keys {
			if inpututil.IsKeyJustPressed(key) {
				if uniformValues == nil {
					uniformValues = make(map[string]any, 4)
				}
				uniformValues[kvu.name] = kvu.value
				break
			}
		}
	}

	//if self.updateFunc != nil { self.updateFunc(&self.options) }
	return nil
}

func (self *shaderDisplayer) Draw(screen *ebiten.Image) {
	screen.Fill(winBackColor)

	bounds := screen.Bounds()
	width, height := float64(bounds.Dx()), float64(bounds.Dy())
	var x, y float32 = 0.0, 0.0
	sxl, sxr := x, x + float32(width)
	syt, syb := y, y + float32(height)
	PositionRectVertices(&self.vertices, sxl, sxr, syt, syb, sxl, sxr, syt, syb)
	indices := []uint16{0, 1, 2, 1, 2, 3}
	
	// uniforms
	if uniformValues == nil { uniformValues = make(map[string]any, 3) }
	uniformValues["Time"] = float32(time.Now().Sub(self.startTime).Seconds())
	cx, cy := ebiten.CursorPosition()
	
	// horzMargin, vertMargin := self.hackyGetMargins(screen)
	var horzMargin, vertMargin float64 = 0, 0
	cx64 := minf64(maxf64(float64(cx) - horzMargin, 0), width)
	cy64 := minf64(maxf64(float64(cy) - vertMargin, 0), height)
	uniformValues["Cursor"] = []float32{ float32(cx64/width), float32(cy64/height) }
	
	var mouseButtons int = 0b00
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft ) { mouseButtons += 0b10 }
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) { mouseButtons += 0b01 }
	uniformValues["MouseButtons"] = mouseButtons
	
	// link uniforms to shader options
	if self.options.Uniforms == nil {
		self.options.Uniforms = make(map[string]any, len(uniformValues))
	}
	for key, value := range uniformValues {
		self.options.Uniforms[key] = value
	}

	// source image linking
	var srcBounds image.Rectangle
	if self.usingImage0 || shaderImage0 != nil {
		if shaderImage0 != nil {
			self.options.Images[0] = shaderImage0
		} else {
			self.options.Images[0] = ImageSpiderCatDog()
		}
		srcBounds = self.options.Images[0].Bounds()
	}
	if self.usingImage1 || shaderImage1 != nil {
		self.options.Images[1] = ImageWaterfall()
		if srcBounds.Empty() { srcBounds = self.options.Images[1].Bounds() }
	}
	if shaderImage2 != nil {
		self.options.Images[2] = shaderImage2
		if srcBounds.Empty() { srcBounds = self.options.Images[2].Bounds() }
	}

	if srcBounds.Empty() { srcBounds = bounds }
	self.vertices[0].SrcX = float32(srcBounds.Min.X) // top-left
	self.vertices[0].SrcY = float32(srcBounds.Min.Y) // top-left
	self.vertices[1].SrcX = float32(srcBounds.Max.X) // top-right
	self.vertices[1].SrcY = float32(srcBounds.Min.Y) // top-right
	self.vertices[2].SrcX = float32(srcBounds.Min.X) // bottom-left
	self.vertices[2].SrcY = float32(srcBounds.Max.Y) // bottom-left
	self.vertices[3].SrcX = float32(srcBounds.Max.X) // bottom-right
	self.vertices[3].SrcY = float32(srcBounds.Max.Y) // bottom-right

	// colors, just in case, so they have some value
	self.vertices[0].ColorR = 1.0 // top-left (red)
	self.vertices[0].ColorG = 0.0 // top-left (red)
	self.vertices[0].ColorB = 0.0 // top-left (red)
	self.vertices[0].ColorA = 1.0 // top-left (red)
	self.vertices[1].ColorR = 0.0 // top-right (green)
	self.vertices[1].ColorG = 1.0 // top-right (green)
	self.vertices[1].ColorB = 0.0 // top-right (green)
	self.vertices[1].ColorA = 1.0 // top-right (green)
	self.vertices[2].ColorR = 0.0 // bottom-left (blue)
	self.vertices[2].ColorG = 0.0 // bottom-left (blue)
	self.vertices[2].ColorB = 1.0 // bottom-left (blue)
	self.vertices[2].ColorA = 1.0 // bottom-left (blue)
	self.vertices[3].ColorR = 1.0 // bottom-right (yellow)
	self.vertices[3].ColorG = 1.0 // bottom-right (yellow)
	self.vertices[3].ColorB = 0.0 // bottom-right (yellow)
	self.vertices[3].ColorA = 1.0 // bottom-right (yellow)

	// actual shader draw call
	screen.DrawTrianglesShader(self.vertices[0 : 4], indices, self.shader, &self.options)

	// handle uniform infos
	for key, _ := range extraUniformInfos {
		extraUniformInfoOrders = append(extraUniformInfoOrders, key)
	}
	sort.Strings(extraUniformInfoOrders)
	for i, name := range extraUniformInfoOrders {
		uniformInfo := extraUniformInfos[name]
		var details string
		if uniformInfo.formatter != nil {
			details = uniformInfo.formatter(uniformValues[name])
		} else {
			valueStr := uniformInfo.FormatValue(uniformValues[name])
			if uniformInfo.replace == "NO-REPLACE" {
				details = uniformInfo.pre + name + ": " + valueStr + uniformInfo.post
			} else {
				details = uniformInfo.pre + uniformInfo.replace + valueStr + uniformInfo.post
			}
		}
		
		ebitenutil.DebugPrintAt(screen, details, bounds.Min.X + 1, bounds.Min.Y + i*13)
	}
	extraUniformInfoOrders = extraUniformInfoOrders[ : 0]
}

func minf64(a, b float64) float64 {
	if a <= b { return a }
	return b
}

func maxf64(a, b float64) float64 {
	if a >= b { return a }
	return b
}

// TODO: unused
func (self *shaderDisplayer) hackyGetMargins(canvas *ebiten.Image) (left, top float64) {
	winWidth, winHeight := ebiten.WindowSize()
	canvasBounds := canvas.Bounds()
	windowAspectRatio := float64(winWidth)/float64(winHeight)
	canvasAspectRatio := float64(canvasBounds.Dx())/float64(canvasBounds.Dy())
	switch {
	case windowAspectRatio == canvasAspectRatio: // just scaling
		return 0, 0
	case windowAspectRatio  > canvasAspectRatio: // horz margins
		xMargin := int((float64(winWidth) - canvasAspectRatio*float64(winHeight))/2.0)
		return float64(xMargin), 0
	case canvasAspectRatio  > windowAspectRatio: // vert margins
		yMargin := int((float64(winHeight) - float64(winWidth)/canvasAspectRatio)/2.0)
		return 0, float64(yMargin)
	default:
		panic("unreachable")
	}
}

func containsOutsideComment(bytes []byte, seq []rune) bool {
	var inComment int = 0 // -1 if not in comment, 1 if in comment, 0 if unknown
	var matchLen int = 0

	var prevRune rune
	var index int = 0
	for index < len(bytes) {
		codePoint, codePointLen := utf8.DecodeRune(bytes[index : ])
		switch inComment {
		case 0: // undecided
			if codePoint == '/' {
				if prevRune == '/' { inComment = 1 }
			} else if codePoint == ' ' || codePoint == '\t' || codePoint == '\n' || codePoint == '\r' {
				// continue undecided
			} else {
				inComment = -1
				matchLen = 0
				if codePoint == seq[0] {
					if len(seq) == 1 { return true }
					matchLen = 1
				}
			}
		case -1: // not in comment
			if codePoint == seq[matchLen] {
				matchLen += 1
				if len(seq) == matchLen { return true }
			} else if codePoint == seq[0] {
				if len(seq) == 1 { return true }
				matchLen = 1
			} else {
				matchLen = 0
			}
			if codePoint == '\n' || codePoint == '\r' {
				inComment = 0
			}
		case 1: // in comment
			if codePoint == '\n' || codePoint == '\r' { inComment = 0 }
		default:
			panic("unreachable")
		}
		prevRune = codePoint
		index += codePointLen
	}

	return false
}
