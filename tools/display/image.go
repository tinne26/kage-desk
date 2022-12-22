package display

import "log"
import "image"

import "github.com/hajimehoshi/ebiten/v2"

func Image(img image.Image, opts ...any) {
	ebitenImage := ebiten.NewImageFromImage(img)
	bounds := ebitenImage.Bounds()
	ebiten.SetWindowSize(bounds.Dx(), bounds.Dy())

	titleSet := false
	for _, opt := range opts {
		switch typedOpt := opt.(type) {
		case string:
			if titleSet {
				log.Fatal("display.Image() option error: received more than one window title")
			}
			ebiten.SetWindowTitle(typedOpt)
			titleSet = true
		default:
			log.Fatalf("display.Image(): unexpected option of type %T", opt)
		}
	}
	if !titleSet {
		ebiten.SetWindowTitle("display/image")
	}

	err := ebiten.RunGame(&imageDisplayer{ ebitenImage })
	if err != nil {
		log.Fatalf("display.Image() failure: %s", err.Error())
	}
}

type imageDisplayer struct {
	img *ebiten.Image
}

func (self *imageDisplayer) Layout(_, _ int) (int, int) {
	bounds := self.img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func (self *imageDisplayer) Update() error {
	return nil
}

func (self *imageDisplayer) Draw(screen *ebiten.Image) {
	screen.DrawImage(self.img, nil)
}
