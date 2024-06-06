package display

import "os"
import "fmt"
import "image"
import "image/png"

import "github.com/hajimehoshi/ebiten/v2"

// Similar concept to [Shader](), but for single images.
// You can press 'E' at any time to export the image
// as a file too.
func Image(img image.Image) {
	if img == nil { panic("can't display nil image") }
	ebitenImage := ebiten.NewImageFromImage(img)

	if !winSizeSet {
		bounds := ebitenImage.Bounds()
		ebiten.SetWindowSize(bounds.Dx(), bounds.Dy())
	}

	ebiten.SetScreenClearedEveryFrame(false)
	err := ebiten.RunGame(&imageDisplayer{ original: img, img: ebitenImage })
	if err != nil && err != errEscClose {
		fail(fmt.Sprintf("display.Image() failure: %s", err.Error()))
	}
}

type imageDisplayer struct {
	original image.Image
	img *ebiten.Image
	needsRedraw bool
	lastWidth, lastHeight int
	exportKeyPressed bool
}

func (self *imageDisplayer) Layout(w, h int) (int, int) {
	if w != self.lastWidth {
		self.lastWidth = w
		self.needsRedraw = true
	}
	if h != self.lastHeight {
		self.lastHeight = h
		self.needsRedraw = true
	}
	return w, h
}

func (self *imageDisplayer) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errEscClose
	}

	exportKeyPressed := ebiten.IsKeyPressed(ebiten.KeyE)
	if exportKeyPressed != self.exportKeyPressed {
		if !self.exportKeyPressed {
			self.exportImageAsPNG()
		}
		self.exportKeyPressed = exportKeyPressed
	}
	return nil
}

func (self *imageDisplayer) exportImageAsPNG() {
	fmt.Printf("Exporting image as png...\n")
	file, err := os.Create("display_image_export.png")
	if err != nil {
		fmt.Printf("Aborted export: %s\n", err.Error())
		return
	}
	defer file.Close()
	
	err = png.Encode(file, self.original)
	if err != nil {
		fmt.Printf("Aborted export: %s\n", err.Error())
		return
	}

	fmt.Print("Successfully exported display_image_export.png\n")
}

func (self *imageDisplayer) Draw(screen *ebiten.Image) {
	if self.needsRedraw {
		self.needsRedraw = false

		// fill background and center image
		screen.Fill(winBackColor)
		scrBounds := screen.Bounds()
		imgBounds := self.img.Bounds()
		x := (scrBounds.Dx() - imgBounds.Dx())/2
		y := (scrBounds.Dy() - imgBounds.Dy())/2
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(self.img, opts)
	}
}
