package display

import "image"
import "github.com/hajimehoshi/ebiten/v2"

// Converts an image.Rectangle to a []float32 for use as a shader uniform.
// 
// The []float32 always has 4 values: min x, min y, max x, max y.
func RectToUniform(rect image.Rectangle) []float32 {
	return []float32{
		float32(rect.Min.X),
		float32(rect.Min.Y),
		float32(rect.Max.X),
		float32(rect.Max.Y),
	}
}

// Similar to ebiten.DrawRectShader, but the options are ebiten.DrawTrianglesShaderOptions
// and the size of the destination rect can be different from the images from the options.
// When this happens, the textures will be stretched to fill the whole shader rectangle.
//
// In general, though, if the shader's target rectangle and the image sizes are different,
// you should compute the texture coordinates on your own.
func DrawShader(screen *ebiten.Image, rect image.Rectangle, shader *ebiten.Shader, opts *ebiten.DrawTrianglesShaderOptions) {
	// vertices will be top-left, top-right,
	// bottom-right, bottom-left (clockwise order)
	var vertices [4]ebiten.Vertex
	vertices[0].DstX = float32(rect.Min.X) // top-left
	vertices[0].DstY = float32(rect.Min.Y)
	vertices[1].DstX = float32(rect.Max.X) // top-right
	vertices[1].DstY = float32(rect.Min.Y)
	vertices[2].DstX = float32(rect.Max.X) // bottom-right
	vertices[2].DstY = float32(rect.Max.Y)
	vertices[3].DstX = float32(rect.Min.X) // bottom-left
	vertices[3].DstY = float32(rect.Max.Y)

	// setup texture coordinates if any image is used
	for _, img := range opts.Images {
		if img != nil {
			bounds := img.Bounds()
			vertices[0].SrcX = float32(bounds.Min.X) // top-left
			vertices[0].SrcY = float32(bounds.Min.Y)
			vertices[1].SrcX = float32(bounds.Max.X) // top-right
			vertices[1].SrcY = float32(bounds.Min.Y)
			vertices[2].SrcX = float32(bounds.Max.X) // bottom-right
			vertices[2].SrcY = float32(bounds.Max.Y)
			vertices[3].SrcX = float32(bounds.Min.X) // bottom-left
			vertices[3].SrcY = float32(bounds.Max.Y)
			break
		}
	}

	// actual draw call
	screen.DrawTrianglesShader(vertices[0 : 4], []uint16{0, 1, 2, 0, 2, 3}, shader, opts)
}
