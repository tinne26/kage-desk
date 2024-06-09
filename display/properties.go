package display

import "os"
import "fmt"
import "image/color"
import "errors"
import "github.com/hajimehoshi/ebiten/v2"

var errEscClose = errors.New("closed game with ESC")

func fail(str string) {
	fmt.Printf("kage-desk/display error: %s\n", str)
	os.Exit(1)
}

func warn(str string) {
	fmt.Printf("kage-desk/display warning: %s\n", str)
}

type WindowOption uint8
const (
	Windowed       WindowOption = 0b0001
	Fullscreen     WindowOption = 0b0010
	Resizable      WindowOption = 0b0100
	DisplayScaling WindowOption = 0b1000
	HiRes          WindowOption = 0b1000 // better name than DisplayScaling
)

type debugMode uint8
const (
	MaxFPS debugMode = 0b0001
)

var argWindowed bool = false
var argFullscreen bool = false
var argMaxFPS bool = false
var winSizeSet bool = false
var winDisplayScaling bool = false
var winResizable bool = false
var winWidth  int = 640
var winHeight int = 480
var winTitle string
var winBackColor = color.RGBA{0, 0, 0, 255}

func init() {
	var argOpenGL bool = false
	for _, arg := range os.Args {
		switch arg {
		case "--maxfps":
			if argMaxFPS { warn("repeated --maxfps program flag") }
			argMaxFPS = true
		case "--windowed", "--fullscreen":
			if arg == "--windowed" {
				if argWindowed { warn("repeated --windowed program flag") }
				argWindowed = true
			}
			if arg == "--fullscreen" {
				if argFullscreen { warn("repeated --fullscreen program flag") }
				argFullscreen = true
			}
			if argFullscreen && argWindowed {
				fail("can't invoke program with both --fullscreen and --windowed flags")
			}
		case "--opengl":
			if argOpenGL {
				warn("repeated --opengl program flag")
			} else {
				argOpenGL = true
				err := os.Setenv("EBITENGINE_GRAPHICS_LIBRARY", "opengl")
				if err != nil { panic(err) }
			}
		default:
			// (allow other program flags or not?)
			//fail("unexpected '" + arg + "' program flag")
		}
	}

	// actually apply most flags
	if argMaxFPS {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	}
	if argFullscreen {
		ebiten.SetFullscreen(true)
	}
}

// Sets a window title. Fairly equivalent to [ebiten.SetWindowTitle](),
// but also stores the title internally in case you pass the --maxfps
// tag, which will display the FPS on the title bar in addition to
// the title.
func SetTitle(title string) {
	winTitle = title
	ebiten.SetWindowTitle(title)
}

// Sets the logical layout size you want to work with.
// Common options include [Resizable] and [HiRes].
func SetSize(width, height int, options ...WindowOption) {
	// safety asserts
	if width < 32 || height < 32 {
		fail(fmt.Sprintf("can't set window size below 32x32 (got %dx%d)", width, height))
	}

	// set properties
	winSizeSet = true
	winWidth, winHeight = width, height
	ebiten.SetWindowSize(width, height)
	actualWidth, actualHeight := ebiten.WindowSize()
	if actualWidth != width || actualHeight != height {
		warn(fmt.Sprintf("requested size %dx%d, but could only get %dx%d", width, height, actualWidth, actualHeight))
	}
	
	// apply window options
	var optsSeen WindowOption
	for _, opt := range options {
		switch opt {
		case Windowed:
			if optsSeen & opt != 0 {
				warn("repeated display.Windowed option in SetSize()")
				continue
			}
			optsSeen = optsSeen | opt
			if argWindowed || argFullscreen { continue }
			ebiten.SetFullscreen(false)
		case Fullscreen:
			if optsSeen & opt != 0 {
				warn("repeated display.Windowed option in SetSize()")
				continue
			}
			optsSeen = optsSeen | opt
			if argWindowed || argFullscreen { continue }
			ebiten.SetFullscreen(true)
		case DisplayScaling:
			if optsSeen & opt != 0 {
				warn("repeated display.DisplayScaling option in SetSize()")
				continue
			}
			optsSeen = optsSeen | opt
			winDisplayScaling = true
		case Resizable:
			if optsSeen & opt != 0 {
				warn("repeated display.Resizable option in SetSize()")
				continue
			}
			optsSeen = optsSeen | opt
			winResizable = true
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		default:
			panic(fmt.Sprintf("invalid window option %d", opt))
		}
	}
}

// Sets a background color to use on the canvas.
func SetBackColor(backColor color.RGBA) {
	winBackColor = backColor
}

var shaderImage0, shaderImage1, shaderImage2, shaderImage3 *ebiten.Image

// Links a specific image for use with shaders. The given n can
// only be 0, 1, 2 or 3.
// 
// By default, two sample textures are already linked for images
// 0 and 1 (see [ImageSpiderCatDog]() and [ImageWaterfall]()).
// You can override them with your own or restore them by setting
// their values back to nil.
func LinkShaderImage(n int, image *ebiten.Image) {
	switch n {
	case 0: shaderImage0 = image
	case 1: shaderImage1 = image
	case 2: shaderImage2 = image
	case 3: shaderImage3 = image
	default:
		panic("n must be between 0 and 3")
	}
}

var uniformValues map[string]any
var keyValueUniforms []keyValueUniform
type keyValueUniform struct {
	name string
	value any
	keys []ebiten.Key
}

// Links a uniform in the following manner: when 'key' is
// triggered, 'value' is set for uniform 'name'. If the
// uniform hadn't been linked yet, the value will also be
// set as the starting value.
func LinkUniformKey(name string, value any, keys ...ebiten.Key) {
	// detect forbidden overrides
	if name == "Time" { panic("can't override 'Time' uniform") }
	if name == "Cursor" { panic("can't override 'Cursor' uniform") }
	if name == "MouseButtons" { panic("can't override 'MouseButtons' uniform") }
	
	// see if key already present for the given name
	for i, _ := range keyValueUniforms {
		if keyValueUniforms[i].name == name {
			var allEqual bool = true
			for i, key := range keyValueUniforms[i].keys {
				if key != keys[i] {
					allEqual = false
					break
				}
			}
			if allEqual {
				keyValueUniforms[i].value = value
				return // early
			}
		}
	}

	// otherwise, append new entry
	keyValueUniforms = append(keyValueUniforms, keyValueUniform{ name, value, keys })
	_, found := uniformValues[name]
	if !found {
		if uniformValues == nil {
			uniformValues = make(map[string]any, 4)
		}
		uniformValues[name] = value
	}
}

var extraUniformInfos map[string]extraUniformInfo
type extraUniformInfo struct {
	verb string
	pre string
	post string
	replace string
	formatter func(value any) string
}
func (self *extraUniformInfo) FormatValue(value any) string {
	switch self.verb {
	case "%hide":
		return ""
	case "%vec2":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(%.02f, %.02f)", f32slice[0], f32slice[1]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(%.02f, %.02f)", f64slice[0], f64slice[1]) }
		f32array, ok := value.([2]float32)
		if ok { return fmt.Sprintf("(%.02f, %.02f)", f32array[0], f32array[1]) }
		f64array, ok := value.([2]float64)
		if ok { return fmt.Sprintf("(%.02f, %.02f)", f64array[0], f64array[1]) }
		panic("special format '%vec2' can only be used with []float32, []float64, [2]float32 or [2]float64")
	case "%vec2[0]":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("%.02f", f32slice[0]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("%.02f", f64slice[0]) }
		f32array, ok := value.([2]float32)
		if ok { return fmt.Sprintf("%.02f", f32array[0]) }
		f64array, ok := value.([2]float64)
		if ok { return fmt.Sprintf("%.02f", f64array[0]) }
		panic("special format '%vec2[0]' can only be used with []float32, []float64, [2]float32 or [2]float64")
	case "%vec2[1]":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("%.02f", f32slice[0]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("%.02f", f64slice[0]) }
		f32array, ok := value.([2]float32)
		if ok { return fmt.Sprintf("%.02f", f32array[0]) }
		f64array, ok := value.([2]float64)
		if ok { return fmt.Sprintf("%.02f", f64array[0]) }
		panic("special format '%vec2[1]' can only be used with []float32, []float64, [2]float32 or [2]float64")
	case "%vec3":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f)", f32slice[0], f32slice[1], f32slice[2]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f)", f64slice[0], f64slice[1], f64slice[2]) }
		f32array, ok := value.([3]float32)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f)", f32array[0], f32array[1], f32array[2]) }
		f64array, ok := value.([3]float64)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f)", f64array[0], f64array[1], f64array[2]) }
		panic("special format '%vec3' can only be used with []float32, []float64, [3]float32 or [3]float64")
	case "%vec4":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f, %.02f)", f32slice[0], f32slice[1], f32slice[2], f32slice[3]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f, %.02f)", f64slice[0], f64slice[1], f64slice[2], f64slice[3]) }
		f32array, ok := value.([4]float32)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f, %.02f)", f32array[0], f32array[1], f32array[2], f32array[3]) }
		f64array, ok := value.([4]float64)
		if ok { return fmt.Sprintf("(%.02f, %.02f, %.02f, %.02f)", f64array[0], f64array[1], f64array[2], f64array[3]) }
		panic("special format '%vec4' can only be used with []float32, []float64, [4]float32 or [4]float64")
	case "%ivec2":
		intSlice, ok := value.([]int)
		if ok { return fmt.Sprintf("(%d, %d)", intSlice[0], intSlice[1]) }
		intArray, ok := value.([2]int)
		if ok { return fmt.Sprintf("(%d, %d)", intArray[0], intArray[1]) }
		panic("special format '%ivec2' can only be used with []int or [2]int")
	case "%ivec4":
		intSlice, ok := value.([]int)
		if ok { return fmt.Sprintf("(%d, %d, %d, %d)", intSlice[0], intSlice[1], intSlice[2], intSlice[3]) }
		intArray, ok := value.([4]int)
		if ok { return fmt.Sprintf("(%d, %d, %d, %d)", intArray[0], intArray[1], intArray[2], intArray[3]) }
		panic("special format '%ivec4' can only be used with []int or [4]int")
	case "%RGB8":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(R%d G%d B%d)", int(f32slice[0]*255), int(f32slice[1]*255), int(f32slice[2]*255)) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(R%d G%d B%d)", int(f64slice[0]*255), int(f64slice[1]*255), int(f64slice[2]*255)) }
		f32array, ok := value.([3]float32)
		if ok { return fmt.Sprintf("(R%d G%d B%d)", int(f32array[0]*255), int(f32array[1]*255), int(f32array[2]*255)) }
		f32array4, ok := value.([4]float32)
		if ok { return fmt.Sprintf("(R%d G%d B%d)", int(f32array4[0]*255), int(f32array4[1]*255), int(f32array4[2]*255)) }
		panic("special format '%RGB8' can only be used with []float32, []float64, [3]float32 or [4]float32")
	case "%RGBA8":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(R%d G%d B%d A%d)", int(f32slice[0]*255), int(f32slice[1]*255), int(f32slice[2]*255), int(f32slice[3]*255)) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(R%d G%d B%d A%d)", int(f64slice[0]*255), int(f64slice[1]*255), int(f64slice[2]*255), int(f64slice[3]*255)) }
		f32array, ok := value.([4]float32)
		if ok { return fmt.Sprintf("(R%d G%d B%d A%d)", int(f32array[0]*255), int(f32array[1]*255), int(f32array[2]*255), int(f32array[3]*255)) }
		panic("special format '%RGBA8' can only be used with []float32, []float64 or [4]float32")
	case "%rgb":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f)", f32slice[0], f32slice[1], f32slice[2]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f)", f64slice[0], f64slice[1], f64slice[2]) }
		f32array, ok := value.([3]float32)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f)", f32array[0], f32array[1], f32array[2]) }
		f32array4, ok := value.([4]float32)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f)", f32array4[0], f32array4[1], f32array4[2]) }
		panic("special format '%rgb' can only be used with []float32, []float64, [3]float32 or [4]float32")
	case "%rgba":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f A%.02f)", f32slice[0], f32slice[1], f32slice[2], f32slice[3]) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f A%.02f)", f64slice[0], f64slice[1], f64slice[2], f64slice[3]) }
		f32array, ok := value.([4]float32)
		if ok { return fmt.Sprintf("(R%.02f G%.02f B%.02f A%.02f)", f32array[0], f32array[1], f32array[2], f32array[3]) }
		panic("special format '%rgba' can only be used with []float32, []float64 or [4]float32")
	case "%percent":
		f32, ok := value.(float32)
		if ok { return fmt.Sprintf("%d%%", int(f32*100.0)) }
		f64, ok := value.(float64)
		if ok { return fmt.Sprintf("%d%%", int(f64*100.0)) }
		panic("special format '%percent' can only be used with float32 or float64")
	case "%vec2-percent":
		f32slice, ok := value.([]float32)
		if ok { return fmt.Sprintf("(%d%%, %d%%)", int(f32slice[0]*100.0), int(f32slice[1]*100.0)) }
		f64slice, ok := value.([]float64)
		if ok { return fmt.Sprintf("(%d%%, %d%%)", int(f64slice[0]*100.0), int(f64slice[1]*100.0)) }
		f32array, ok := value.([2]float32)
		if ok { return fmt.Sprintf("(%d%%, %d%%)", int(f32array[0]*100.0), int(f32array[1]*100.0)) }
		f64array, ok := value.([2]float64)
		if ok { return fmt.Sprintf("(%d%%, %d%%)", int(f64array[0]*100.0), int(f64array[1]*100.0)) }
		panic("special format '%vec2-percent' can only be used with []float32, []float64, [2]float32 or [2]float64")
	default:
		return fmt.Sprintf(self.verb, value)
	}
}
var extraUniformInfoOrders []string

// Allows displaying additional information for a given uniform.
// You can pass up to three strings. If you don't pass any, the
// information will be displayed as "UniformName: {value}". The first
// string you can pass will be used as a prefix, the second as a
// suffix, and the third would replace the original "UniformName: "
// completely. For example, ("Mode", "%d", "[", "]") would result
// in "[Mode: 0]", while ("Cursor", "%vec2[0]", "", "(move cursor)",
// "X: ") would result in something like "X: 0.14 (move cursor)".
//
// By default, no info is shown.
// 
// Example setup in combination with [LinkUniformKey]():
//   display.LinkUniformKey("Mode", 1, ebiten.KeyDigit1, ebiten.KeyNumpad1)
//   display.LinkUniformKey("Mode", 2, ebiten.KeyDigit2, ebiten.KeyNumpad2)
//   display.LinkUniformKey("Mode", 3, ebiten.KeyDigit3, ebiten.KeyNumpad3)
//   display.SetUniformInfo("Mode", "%d", "", " [change with 1, 2, 3]")
//
// There are some additional special verbs for ease of use: "%vec2", "%vec3",
// "%vec4", "%ivec2", "%ivec4", "%RGB8", "%rgb", "%RGBA8", "%rgba", "%percent",
// "%vec2-percent", "%vec2[0]", "%vec2[1]", "%hide". If none of these are
// enough, see [SetUniformFmt]().
func SetUniformInfo(name, verb string, infos ...string) {
	if len(infos) > 3 {
		panic("SetUniformInfo() doesn't accept more than three info arguments: pre, post, replace")
	}

	var preInfo, postInfo, replaceInfo string
	replaceInfo = "NO-REPLACE"
	if len(infos) > 0 { preInfo  = infos[0] }
	if len(infos) > 1 { postInfo = infos[1] }
	if len(infos) > 2 { replaceInfo = infos[2] }
	if verb == "" && len(infos) == 0 {
		delete(extraUniformInfos, name)
	} else {
		if verb == "" { verb = "%v" }
		if extraUniformInfos == nil {
			extraUniformInfos = make(map[string]extraUniformInfo, 1)
		}
		extraUniformInfos[name] = extraUniformInfo{ verb, preInfo, postInfo, replaceInfo, nil }
		if len(extraUniformInfos) > cap(extraUniformInfoOrders) {
			extraUniformInfoOrders = extraUniformInfoOrders[ : cap(extraUniformInfoOrders)]
			extraUniformInfoOrders = append(extraUniformInfoOrders, "")
			extraUniformInfoOrders = extraUniformInfoOrders[ : 0]
		}
	}
}

// An adapter to return a static string from func(any). Intended
// for [SetUniformFmt](" ", display.StrFn("Some information")).
func StrFn(str string) func(any) string {
	return func(any) string { return str }
}

// Similar to [SetUniformInfo](), but with a fully customizable formatter.
// For fixed strings, you may rely on [StrFn]().
func SetUniformFmt(name string, formatter func(any) string) {
	info, found := extraUniformInfos[name]
	if extraUniformInfos == nil {
		extraUniformInfos = make(map[string]extraUniformInfo, 1)
	}

	if found {
		info.formatter = formatter
		extraUniformInfos[name] = info
	} else {
		var verb, preInfo, postInfo, replaceInfo string
		verb = "%v"
		extraUniformInfos[name] = extraUniformInfo{ verb, preInfo, postInfo, replaceInfo, formatter }
	}
}
