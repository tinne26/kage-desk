# Access textures

> [!WARNING]
> **_Since the addition of pixels mode in Ebitengine v2.6.0, the content of this page is no longer relevant; it's preserved only for historical context_**.


For explanations of texels and these functions, see [docs/tutorials/texels.md](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/texels.md).

```Golang
// Helper function to access an image's color at the given pixel coordinates.
func imageColorAtPixel(pixelCoords vec2) vec4 {
	sizeInPixels := imageSrcTextureSize()
	offsetInTexels, _ := imageSrcRegionOnTexture()
	adjustedTexelCoords := pixelCoords/sizeInPixels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}
```

```Golang
// Helper function to access an image's color at the given coordinates
// from the unit interval (e.g. top-left is (0, 0), center is (0.5, 0.5),
// bottom-right is (1.0, 1.0)).
func imageColorAtUnit(unitCoords vec2) vec4 {
	offsetInTexels, sizeInTexels := imageSrcRegionOnTexture()
	adjustedTexelCoords := unitCoords*sizeInTexels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}
```
