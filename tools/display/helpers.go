package display

import "image"
import "github.com/hajimehoshi/ebiten/v2"

// Given 4 vertices and source x left, right, y top, bottom and dest x left, right, y top and bottom
// coordinates, assigns them to vertices 0 (top-left), 1 (top-right), 2 (bottom-left), 3 (bottom-right).
func PositionRectVertices(vertices *[4]ebiten.Vertex, sxl, sxr, syt, syb, dxl, dxr, dyt, dyb float32) {
	vertices[0].SrcX = sxl
	vertices[0].SrcY = syt
	vertices[0].DstX = dxl
	vertices[0].DstY = dyt
	vertices[1].SrcX = sxr
	vertices[1].SrcY = syt
	vertices[1].DstX = dxr
	vertices[1].DstY = dyt
	vertices[2].SrcX = sxl
	vertices[2].SrcY = syb
	vertices[2].DstX = dxl
	vertices[2].DstY = dyb
	vertices[3].SrcX = sxr
	vertices[3].SrcY = syb
	vertices[3].DstX = dxr
	vertices[3].DstY = dyb
}

// Given an image.Rectangle, returns the min and max x and min and max y coordinates
// as float32 values, in this order.
func RectToF32(rect image.Rectangle) (xl, xr, yt, yb float32) {
	return float32(rect.Min.X), float32(rect.Max.X), float32(rect.Min.Y), float32(rect.Max.Y)
}
