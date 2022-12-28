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

func SetTitle(title string) {
	winTitle = title
	ebiten.SetWindowTitle(title)
}

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

func SetBackColor(backColor color.RGBA) {
	winBackColor = backColor
}
