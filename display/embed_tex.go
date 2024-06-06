package display

import _ "embed"
import "image/png"
import "bytes"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed spider_cat_dog.png
var imageSpiderCatDogBytes []byte

//go:embed waterfall.png
var imageWaterfallBytes []byte

var imageSpiderCatDog *ebiten.Image
var imageWaterfall *ebiten.Image

// The source 0 image used for [Shader]().
// You can access it manually too. 384x384.
func ImageSpiderCatDog() *ebiten.Image {
	if imageSpiderCatDog == nil {
		reader := bytes.NewReader(imageSpiderCatDogBytes)
		img, err := png.Decode(reader)
		if err != nil { panic(err) }
		imageSpiderCatDog = ebiten.NewImageFromImage(img)
	}
	return imageSpiderCatDog
}

// The source 1 image used for [Shader]().
// You can access it manually too. 512x512.
func ImageWaterfall() *ebiten.Image {
	if imageWaterfall == nil {
		reader := bytes.NewReader(imageWaterfallBytes)
		img, err := png.Decode(reader)
		if err != nil { panic(err) }
		imageWaterfall = ebiten.NewImageFromImage(img)
	}
	return imageWaterfall
}

