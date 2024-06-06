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

func ImageSpiderCatDog() *ebiten.Image {
	if imageSpiderCatDog == nil {
		reader := bytes.NewReader(imageSpiderCatDogBytes)
		img, err := png.Decode(reader)
		if err != nil { panic(err) }
		imageSpiderCatDog = ebiten.NewImageFromImage(img)
	}
	return imageSpiderCatDog
}

func ImageWaterfall() *ebiten.Image {
	if imageWaterfall == nil {
		reader := bytes.NewReader(imageWaterfallBytes)
		img, err := png.Decode(reader)
		if err != nil { panic(err) }
		imageWaterfall = ebiten.NewImageFromImage(img)
	}
	return imageWaterfall
}

