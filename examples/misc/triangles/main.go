package main

import "os"
import "fmt"
import "log"
import "math"
import "image"
import "image/color"
import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/etxt"
import "github.com/tinne26/fonts/liberation/lbrtmono"

//go:embed point.kage
var shaderProgram []byte

var colorBack  = color.RGBA{155, 192, 163, 255}
var colorText  = color.RGBA{ 35,  31,  32, 255}
var colorPoint = color.RGBA{ 35,  31,  32, 216}
var colorArea  = color.RGBA{229, 255, 222, 255}
var colorTexA  = color.RGBA{255, 186,  73, 255}
var colorTexB  = color.RGBA{ 34, 100, 109, 255}

func init() {
	err := os.Setenv("EBITENGINE_SCREENSHOT_KEY", "q")
	if err != nil { log.Fatal(err) }
	
	for _, arg := range os.Args {
		if arg == "--opengl" {
			err := os.Setenv("EBITENGINE_GRAPHICS_LIBRARY", "opengl")
			if err != nil { log.Fatal(err) }
		}
	}
}

func main() {
	// compile shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { log.Fatal(err) }
	
	// create text renderer
	textRenderer := etxt.NewStdRenderer()
	cache := etxt.NewDefaultCache(16*1024*1024) // 16MiB
	textRenderer.SetCacheHandler(cache.NewHandler())
	textRenderer.SetFont(lbrtmono.Font())

	// create texture (256x192)
	img := image.NewRGBA(image.Rect(0, 0, 256, 192))
	for col := 0; col < 256; col++ {
		colColor := lerpColors(colorTexA, colorTexB, float64(col)/255)
		for row := 0; row < 192; row++ {
			img.SetRGBA(col, row, colColor)
		}
	}
	texture := ebiten.NewImageFromImage(img)

	// setup window and 
	ebiten.SetWindowTitle("examples/misc/triangles")
	ebiten.SetScreenClearedEveryFrame(false)
	x, y, width, height := GetWindowSize(4, 3, 0.66)
	ebiten.SetWindowPosition(x, y)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err = ebiten.RunGame(&App{
		textRenderer: textRenderer, texture: texture,
		dstAH: 0.15, dstAV: 0.15,
		dstBH: 0.15, dstBV: 0.85,
		dstCH: 0.85, dstCV: 0.85,
		srcAH: 0.15, srcAV: 0.15,
		srcBH: 0.15, srcBV: 0.85,
		srcCH: 0.85, srcCV: 0.85,
		pointShader: shader,
		colorR: 1.0,
		colorG: 1.0,
		colorB: 1.0,
		colorA: 1.0,
	})
	if err != nil { log.Fatal(err) }
}

type App struct {
	textRenderer *etxt.Renderer
	triangleDstRect image.Rectangle
	dstAH, dstAV float64
	dstBH, dstBV float64
	dstCH, dstCV float64
	textureSrcRect  image.Rectangle
	srcAH, srcAV float64
	srcBH, srcBV float64
	srcCH, srcCV float64
	colorsRect image.Rectangle
	colorR float32
	colorG float32
	colorB float32
	colorA float32
	texture *ebiten.Image
	cent float64
	mouseHolding string // "" none, "DstA", "DstB", "SrcA", "ClrA", "ClrR", etc. big switch.
	mousePos image.Point
	fscrKeyPressed bool
	pointShader *ebiten.Shader
}

func lerpColors(colorA, colorB color.RGBA, step float64) color.RGBA {
	r := lerpU8(colorA.R, colorB.R, step)
	g := lerpU8(colorA.G, colorB.G, step)
	b := lerpU8(colorA.B, colorB.B, step)
	a := lerpU8(colorA.A, colorB.A, step)
	return color.RGBA{r, g, b, a}
}

func lerpU8(a, b uint8, step float64) uint8 {
	return uint8(float64(a) + (float64(b) - float64(a))*step)
}

func (self *App) LayoutF(logWinWidth, logWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth  := math.Ceil(float64(logWinWidth)*scale)
	canvasHeight := math.Ceil(float64(logWinHeight)*scale)
	return canvasWidth, canvasHeight
}
func (self *App) Layout(logWinWidth, logWinHeight int) (int, int) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth  := int(math.Ceil(float64(logWinWidth)*scale))
	canvasHeight := int(math.Ceil(float64(logWinHeight)*scale))
	return canvasWidth, canvasHeight
}

func (self *App) Update() error {
	if self.cent == 0 { return nil } // wait for first draw

	fscrKeyPressed := ebiten.IsKeyPressed(ebiten.KeyF)
	if self.fscrKeyPressed != fscrKeyPressed {
		self.fscrKeyPressed = fscrKeyPressed
		if fscrKeyPressed {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
	}

	// helper functions
	var dist = func(ox, oy, fx, fy float64) float64 {
		x2 := (ox - fx)*(ox - fx)
		y2 := (oy - fy)*(oy - fy)
		return math.Sqrt(x2 + y2)
	}

	// get cursor position
	isPointing := false
	x, y := ebiten.CursorPosition()
	xy := image.Pt(x, y)

	// check if dragging something
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if self.mouseHolding != "" {
		if mousePressed {
			switch self.mouseHolding { // pressing and dragging case
			case "DstA": self.applyPosDelta(&self.dstAH, &self.dstAV, self.triangleDstRect, xy)
			case "DstB": self.applyPosDelta(&self.dstBH, &self.dstBV, self.triangleDstRect, xy)
			case "DstC": self.applyPosDelta(&self.dstCH, &self.dstCV, self.triangleDstRect, xy)
			case "SrcA": self.applyPosDelta(&self.srcAH, &self.srcAV, self.textureSrcRect, xy)
			case "SrcB": self.applyPosDelta(&self.srcBH, &self.srcBV, self.textureSrcRect, xy)
			case "SrcC": self.applyPosDelta(&self.srcCH, &self.srcCV, self.textureSrcRect, xy)
			case "ClrR": self.applyColorDelta(&self.colorR, xy)
			case "ClrG": self.applyColorDelta(&self.colorG, xy)
			case "ClrB": self.applyColorDelta(&self.colorB, xy)
			case "ClrA": self.applyColorDelta(&self.colorA, xy)
			default:
				panic(self.mouseHolding)
			}
		} else { // click release
			self.mouseHolding = ""
		}
		return nil
	}

	// check if clicking on a target that may be dragged later
	self.mousePos = xy
	if xy.In(self.colorsRect) {
		isPointing = true
		if mousePressed {
			vertNormPos := float64(xy.Y - self.colorsRect.Min.Y)/float64(self.colorsRect.Dy())
			if vertNormPos <= 0.25 {
				self.mouseHolding = "ClrR"
			} else if vertNormPos <= 0.5 {
				self.mouseHolding = "ClrG"
			} else if vertNormPos <= 0.75 {
				self.mouseHolding = "ClrB"
			} else {
				self.mouseHolding = "ClrA"
			}
		}
	} else if xy.In(self.triangleDstRect) {
		rct := self.triangleDstRect
		ax, ay, bx, by, cx, cy := projectPoints(rct, self.dstAH, self.dstAV, self.dstBH, self.dstBV, self.dstCH, self.dstCV)
		fx, fy := float64(xy.X), float64(xy.Y)
		minDist := 9999.0
		holdTarget := ""
		nextDist := dist(ax, ay, fx, fy)
		if nextDist < minDist { minDist = nextDist ; holdTarget = "DstA" }
		nextDist = dist(bx, by, fx, fy)
		if nextDist < minDist { minDist = nextDist ; holdTarget = "DstB" }
		nextDist = dist(cx, cy, fx, fy)
		if nextDist < minDist { minDist = nextDist ; holdTarget = "DstC" }

		// apply new target
		if minDist <= self.cent/2 {
			isPointing = true
			if mousePressed {
				self.mouseHolding = holdTarget
			}
		}
	} else if xy.In(self.textureSrcRect) {
		rct := self.textureSrcRect
		ax, ay, bx, by, cx, cy := projectPoints(rct, self.srcAH, self.srcAV, self.srcBH, self.srcBV, self.srcCH, self.srcCV)
		fx, fy := float64(xy.X), float64(xy.Y)
		minDist := 9999.0
		holdTarget := ""
		nextDist := dist(ax, ay, fx, fy)
		if nextDist < minDist { minDist = nextDist ; holdTarget = "SrcA" }
		nextDist = dist(bx, by, fx, fy)
		if nextDist < minDist { minDist = nextDist ; holdTarget = "SrcB" }
		nextDist = dist(cx, cy, fx, fy)
		if nextDist < minDist { minDist = nextDist ; holdTarget = "SrcC" }

		// apply new target
		if minDist <= self.cent/2 {
			isPointing = true
			if mousePressed {
				self.mouseHolding = holdTarget
			}
		}
	}

	// adjust cursor shape
	if isPointing {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	return nil
}

func (self *App) applyPosDelta(horz *float64, vert *float64, rect image.Rectangle, xy image.Point) {
	newHorzValue := float64(xy.X - rect.Min.X)/float64(rect.Dx())
	if newHorzValue > 1.0 { newHorzValue = 1.0 }
	if newHorzValue < 0.0 { newHorzValue = 0.0 }
	*horz = newHorzValue
	newVertValue := float64(xy.Y - rect.Min.Y)/float64(rect.Dy())
	if newVertValue > 1.0 { newVertValue = 1.0 }
	if newVertValue < 0.0 { newVertValue = 0.0 }
	*vert = newVertValue
}
func (self *App) applyColorDelta(clr *float32, xy image.Point) {
	delta := (xy.Y - self.mousePos.Y)
	if delta == 0 { return }

	change := float32(int(float64(delta)*6.0/self.cent))*0.01
	if change == 0 { return }
	self.mousePos = xy
	var newClrValue float32 = *clr - change
	if newClrValue > 1.0 { newClrValue = 1.0 }
	if newClrValue < 0.0 { newClrValue = 0.0 }
	*clr = newClrValue
}

func (self *App) Draw(screen *ebiten.Image) {
	canvas := GetCanvas(screen, 4, 3)
	canvas.Fill(colorBack)
	self.textRenderer.SetTarget(canvas)

	self.drawTriangleDstBlock(canvas)
	self.drawTextureSrcBlock(canvas)
	self.drawResultBlock(canvas)
}

func (self *App) drawTriangleDstBlock(canvas *ebiten.Image) {
	bounds := canvas.Bounds()
	width  := bounds.Dx()
	height := bounds.Dy()
	origin := image.Pt(bounds.Min.X, bounds.Min.Y)

	heightFactor := 0.26
	widthFactor  := heightFactor*(4.0/3.0)
	areaWidth  := int(float64(height)* widthFactor)
	areaHeight := int(float64(height)*heightFactor)
	self.cent = float64(height)*0.01

	x :=  width/2 - areaWidth  - int(self.cent*6)
	y := height/2 - areaHeight - int(self.cent*8)
	self.triangleDstRect = image.Rect(x, y, x + areaWidth, y + areaHeight).Add(origin)
	sub := canvas.SubImage(self.triangleDstRect).(*ebiten.Image)
	sub.Fill(colorArea)

	// draw + and = symbols
	rect := image.Rect(width/2 - int(self.cent*3), y + areaHeight/2 - int(self.cent/2), width/2 + int(self.cent*3), y + areaHeight/2 + int(self.cent/2))
	canvas.SubImage(rect.Add(origin)).(*ebiten.Image).Fill(colorArea)
	rect  = image.Rect(width/2 - int(self.cent/2), y + areaHeight/2 - int(self.cent*3), width/2 + int(self.cent/2), y + areaHeight/2 + int(self.cent*3))
	canvas.SubImage(rect.Add(origin)).(*ebiten.Image).Fill(colorArea)

	rect = image.Rect(width/2 - int(self.cent*3), height/2 - int(self.cent*1.2), width/2 + int(self.cent*3), height/2 - int(self.cent*1.2) + int(self.cent/2)*2)
	canvas.SubImage(rect.Add(origin)).(*ebiten.Image).Fill(colorArea)
	rect = image.Rect(width/2 - int(self.cent*3), height/2 + int(self.cent*1.2), width/2 + int(self.cent*3), height/2 + int(self.cent*1.2) + int(self.cent/2)*2)
	canvas.SubImage(rect.Add(origin)).(*ebiten.Image).Fill(colorArea)

	// draw text
	self.textRenderer.SetAlign(etxt.Bottom, etxt.XCenter)
	self.textRenderer.SetColor(colorText)
	self.textRenderer.SetSizePx(int(self.cent*2.9))
	pt := image.Pt(x + areaWidth/2, y - int(self.cent*2.4))
	pt  = pt.Add(origin)
	self.textRenderer.Draw("vertex Dst positions", pt.X, pt.Y)

	// draw triangle Dst points
	ctrx, ctry := Centroid(self.dstAH, self.dstAV, self.dstBH, self.dstBV, self.dstCH, self.dstCV)
	self.textRenderer.SetColor(colorPoint)
	self.textRenderer.SetSizePx(int(self.cent*2.1))
	rct := self.triangleDstRect
	ax, ay, bx, by, cx, cy := projectPoints(rct, self.dstAH, self.dstAV, self.dstBH, self.dstBV, self.dstCH, self.dstCV)
	self.drawPoint(canvas, ax, ay)
	self.setAlignForPoint(self.dstAH, self.dstAV, ctrx, ctry)
	self.textRenderer.Draw(" A ", int(ax), int(ay))
	self.drawPoint(canvas, bx, by)
	self.setAlignForPoint(self.dstBH, self.dstBV, ctrx, ctry)
	self.textRenderer.Draw(" B ", int(bx), int(by))
	self.drawPoint(canvas, cx, cy)
	self.setAlignForPoint(self.dstCH, self.dstCV, ctrx, ctry)
	self.textRenderer.Draw(" C ", int(cx), int(cy))
}

func (self *App) setAlignForPoint(x, y, ctrx, ctry float64) {
	var horzAlign etxt.HorzAlign
	var vertAlign etxt.VertAlign
	if x <= ctrx { horzAlign = etxt.Right  } else { horzAlign = etxt.Left }
	if y <= ctry { vertAlign = etxt.Bottom } else { vertAlign = etxt.Top  }
	self.textRenderer.SetAlign(vertAlign, horzAlign)
}

func (self *App) drawPoint(screen *ebiten.Image, x, y float64) {
	radius := self.cent/2
	opts := &ebiten.DrawRectShaderOptions{}
	trX := float64(math.Floor(x - radius))
	trY := float64(math.Floor(y - radius))
	sizeX := math.Ceil(x + radius) - trX
	sizeY := math.Ceil(y + radius) - trY
	size := sizeX
	if sizeY > size { size = sizeY }
	opts.GeoM.Translate(trX, trY)
	opts.Uniforms = map[string]interface{}{
		"Center": []float32{float32(x), float32(y)},
		"Radius": float32(radius),
	}
	screen.DrawRectShader(int(size), int(size), self.pointShader, opts)
}

func (self *App) drawTextureSrcBlock(canvas *ebiten.Image) {
	bounds := canvas.Bounds()
	width  := bounds.Dx()
	
	areaWidth  := int(float64(self.triangleDstRect.Dx())*0.6)
	areaHeight := int(float64(self.triangleDstRect.Dy())*0.6)
	
	x := width/2 + int(self.cent*6) + bounds.Min.X
	y := self.triangleDstRect.Min.Y + (self.triangleDstRect.Dy() - areaHeight)/2
	self.textureSrcRect = image.Rect(x, y, x + areaWidth, y + areaHeight)
	
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(float64(self.textureSrcRect.Dx())/256.0, float64(self.textureSrcRect.Dy())/192.0)
	opts.GeoM.Translate(float64(x), float64(y))
	canvas.DrawImage(self.texture, opts)

	self.textRenderer.SetAlign(etxt.Bottom, etxt.XCenter)
	self.textRenderer.SetColor(colorText)
	self.textRenderer.SetSizePx(int(self.cent*2.9))
	self.textRenderer.Draw("texture and\nSrc positions", x + areaWidth/2, y - int(self.cent*2.4))

	self.textRenderer.SetAlign(etxt.YCenter, etxt.Left)
	self.textRenderer.SetSizePx(int(self.cent*2.4))
	clr := colorText
	clr.A = 216
	self.textRenderer.SetColor(clr)
	str := fmt.Sprintf("ColorR: %.2f\nColorG: %.2f\nColorB: %.2f\nColorA: %.2f", self.colorR, self.colorG, self.colorB, self.colorA)
	textX := x + areaWidth + int(self.cent*2.4)
	textY := y + areaHeight/2
	self.textRenderer.Draw(str, textX, textY)

	textX += self.textRenderer.SelectionRect("ColorA: ").Width.Ceil()
	textWidth  := self.textRenderer.SelectionRect("0.00").Width.Ceil()
	textHeight := self.textRenderer.GetLineAdvance().Ceil()*4
	self.colorsRect = image.Rect(textX, textY - textHeight/2, textX + textWidth, textY + textHeight/2)

	// draw texture Src points
	ctrx, ctry := Centroid(self.srcAH, self.srcAV, self.srcBH, self.srcBV, self.srcCH, self.srcCV)
	self.textRenderer.SetColor(colorPoint)
	self.textRenderer.SetSizePx(int(self.cent*2.1))
	rct := self.textureSrcRect
	ax, ay, bx, by, cx, cy := projectPoints(rct, self.srcAH, self.srcAV, self.srcBH, self.srcBV, self.srcCH, self.srcCV)
	self.drawPoint(canvas, ax, ay)
	self.setAlignForPoint(self.srcAH, self.srcAV, ctrx, ctry)
	self.textRenderer.Draw(" A ", int(ax), int(ay))
	self.drawPoint(canvas, bx, by)
	self.setAlignForPoint(self.srcBH, self.srcBV, ctrx, ctry)
	self.textRenderer.Draw(" B ", int(bx), int(by))
	self.drawPoint(canvas, cx, cy)
	self.setAlignForPoint(self.srcCH, self.srcCV, ctrx, ctry)
	self.textRenderer.Draw(" C ", int(cx), int(cy))
}

func (self *App) drawResultBlock(canvas *ebiten.Image) {
	bounds := canvas.Bounds()
	width  := bounds.Dx()
	height := bounds.Dy()
	origin := image.Pt(bounds.Min.X, bounds.Min.Y)

	heightFactor := 0.26
	widthFactor  := heightFactor*(4.0/3.0)
	areaWidth  := int(float64(height)* widthFactor)
	areaHeight := int(float64(height)*heightFactor)

	x := (width - areaWidth)/2
	y := height/2 + int(self.cent*13)
	rect := image.Rect(x, y, x + areaWidth, y + areaHeight).Add(origin)
	sub := canvas.SubImage(rect).(*ebiten.Image)
	sub.Fill(colorArea)

	self.textRenderer.SetAlign(etxt.Bottom, etxt.XCenter)
	self.textRenderer.SetColor(colorText)
	self.textRenderer.SetSizePx(int(self.cent*2.9))
	pt := image.Pt(x + areaWidth/2, y - int(self.cent*2.4))
	pt  = pt.Add(origin)
	self.textRenderer.Draw("triangle render", pt.X, pt.Y)

	// render actual triangle
	vertices := make([]ebiten.Vertex, 3)
	rct := rect
	dax, day, dbx, dby, dcx, dcy := projectPoints(rct, self.dstAH, self.dstAV, self.dstBH, self.dstBV, self.dstCH, self.dstCV)
	vertices[0].DstX = float32(dax); vertices[0].DstY = float32(day)
	vertices[1].DstX = float32(dbx); vertices[1].DstY = float32(dby)
	vertices[2].DstX = float32(dcx); vertices[2].DstY = float32(dcy)
	rct  = image.Rect(0, 0, 256, 192)
	sax, say, sbx, sby, scx, scy := projectPoints(rct, self.srcAH, self.srcAV, self.srcBH, self.srcBV, self.srcCH, self.srcCV)
	vertices[0].SrcX = float32(sax); vertices[0].SrcY = float32(say)
	vertices[1].SrcX = float32(sbx); vertices[1].SrcY = float32(sby)
	vertices[2].SrcX = float32(scx); vertices[2].SrcY = float32(scy)
	for i := 0; i < 3; i++ {
		vertices[i].ColorR = self.colorR
		vertices[i].ColorG = self.colorG
		vertices[i].ColorB = self.colorB
		vertices[i].ColorA = self.colorA
	}
	canvas.DrawTriangles(vertices, []uint16{0, 1, 2}, self.texture, nil)
}

func projectPoints(rect image.Rectangle, ax, ay, bx, by, cx, cy float64) (x1, y1, x2, y2, x3, y3 float64) {
	x1, y1 = projectPoint(rect, ax, ay)
	x2, y2 = projectPoint(rect, bx, by)
	x3, y3 = projectPoint(rect, cx, cy)
	return
}

func projectPoint(rect image.Rectangle, x, y float64) (float64, float64) {
	return float64(rect.Min.X) + float64(rect.Dx())*x, float64(rect.Min.Y) + float64(rect.Dy())*y
}

func Centroid(ax, ay, bx, by, cx, cy float64) (float64, float64) {
	return (ax + bx + cx)/3, (ay + by + cy)/3
}

func GetCanvas(screen *ebiten.Image, widthRatio, heightRatio int) *ebiten.Image {
	bounds := screen.Bounds()
	workWidth  := bounds.Dx()
	workHeight := bounds.Dy()

	var width, height int
	desiredRatio := float64(widthRatio)/float64(heightRatio)
	workingRatio := float64(workWidth)/float64(workHeight)
	if desiredRatio >= workingRatio {
		width  = workWidth
		fheight := (float64(width)*float64(heightRatio))/float64(widthRatio)
		height = int(math.Ceil(fheight))
	} else { // workingRatio > desiredRatio
		height  = workHeight
		fwidth := (float64(height)*float64(widthRatio))/float64(heightRatio)
		width = int(math.Ceil(fwidth))
	}

	xOffset := (workWidth  - width )/2
	yOffset := (workHeight - height)/2
	rect := image.Rect(xOffset, yOffset, xOffset + width, yOffset + height)
	return screen.SubImage(rect).(*ebiten.Image)
}

// Returns the top-left (x, y) and the width and height
func GetWindowSize(widthRatio, heightRatio int, sizeFactor float64) (int, int, int, int) {
	if sizeFactor > 1.0 { panic("sizeFactor > 1.0 (size factor must stay between 0.1 and 1.0)") }
	if sizeFactor < 0.1 { panic("sizeFactor < 0.1 (size factor must stay between 0.1 and 1.0)") }

	_, _, maxWidth, maxHeight := ebiten.WindowSizeLimits()
	if maxWidth == -1 || maxHeight == -1 {
		fscrWidth, fscrHeight := ebiten.ScreenSizeInFullscreen()
		if maxWidth  == -1 || fscrWidth  < maxWidth  { maxWidth  = fscrWidth  }
		if maxHeight == -1 || fscrHeight < maxHeight { maxHeight = fscrHeight }
	}

	desiredRatio := float64(widthRatio)/float64(heightRatio)
	workingRatio := float64(maxWidth)/float64(maxHeight)
	if desiredRatio >= workingRatio {
		width  := int(float64(maxWidth)*sizeFactor)
		width  -= width % widthRatio
		height := (width*heightRatio)/widthRatio
		return (maxWidth - width)/2, (maxHeight - height)/2, width, height
	} else { // workingRatio > desiredRatio
		height := int(float64(maxHeight)*sizeFactor)
		height -= height % heightRatio
		width  := (height*widthRatio)/heightRatio
		return (maxWidth - width)/2, (maxHeight - height)/2, width, height
	}
}
