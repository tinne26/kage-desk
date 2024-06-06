package display

import "image/color"

// Predefined background colors, to make life easier.
var (
	BCBlack   = color.RGBA{  0,   0,   0, 255}
	BCDark    = color.RGBA{ 12,  20,  18, 255} // dark but not black
	BCWhite   = color.RGBA{255, 255, 255, 255}
	BCGray    = color.RGBA{128, 128, 128, 255}
	BCBronze  = color.RGBA{192, 128,  64, 255} // a bit of everything
	BCOrchid  = color.RGBA{128,  64, 192, 255} // a bit of everything
	BCJade    = color.RGBA{  8, 162, 116, 255} // a bit of everything
	BCRed     = color.RGBA{255,   0,   0, 255}
	BCGreen   = color.RGBA{  0, 255,   0, 255}
	BCMagenta = color.RGBA{255,   0, 255, 255}
	BCCyan    = color.RGBA{  0, 255, 255, 255}
)
