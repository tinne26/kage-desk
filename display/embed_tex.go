package display

import _ "embed"
import "image/png"
import "bytes"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed spider_cat_dog.png
var imageSpiderCatDogBytes []byte
var imageSpiderCatDog *ebiten.Image

func ImageSpiderCatDog() *ebiten.Image {
	if imageSpiderCatDog == nil {
		reader := bytes.NewReader(imageSpiderCatDogBytes)
		img, err := png.Decode(reader)
		if err != nil { panic(err) }
		imageSpiderCatDog = ebiten.NewImageFromImage(img)
	}
	return imageSpiderCatDog
}

