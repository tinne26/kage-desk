# Triangles in Ebitengine

The only geometry that GPUs understand are triangles. Triangles are not only cool because they are the simplest type of polygon, but because any polygon can be decomposed into triangles. We can create any shapes we want using only triangles!

In 3D game engines, triangles are omnipresent. In 2D game engines, though, triangles might go unnoticed.

The reason is that with 2D game engines, we most often only care about sprites and textures, which are usually rectangular. Even if triangles are still used under the hood, they may not be exposed in the APIs.

Luckily, Ebitengine has two API calls that allow us to work directly with triangles:
- [Image.DrawTriangles(...)](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawTriangles), which allows us to draw triangles with colored vertices or some specific texture.
- [Image.DrawTrianglesShader(...)](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawTrianglesShader), which is similar to `DrawTriangles(...)`, but we also specify a shader to draw with.

## Understanding `DrawTriangles()`

Lots of people are confused about the arguments passed to these triangle-drawing functions when they first see them:
```Golang
DrawTriangles(vertices []Vertex, indices []uint16, img *Image, options *DrawTrianglesOptions)
```

The next sections explain the parameters one by one and how to use them.

### Vertices and indices

Let's start with the `vertices []Vertex` and the `indices []uint16`, which define the geometry of the shape to be drawn.

Say we want to draw a triangle with vertices at `(0, 0)`, `(10, 0)` and `(10, 10)`. We need to create three vertices and set these coordinates on the `DstX` and `DstY` fields of each [`Vertex`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Vertex) struct. Once we have this, we can put the vertices into a `[]Vertex` slice, which will be the first parameter of the `DrawTriangles()` function:
```Golang
vertices := make([]Vertex, 3)
vertices[0].DstX = float32( 0)
vertices[0].DstY = float32( 0)
vertices[1].DstX = float32(10)
vertices[1].DstY = float32( 0)
vertices[2].DstX = float32(10)
vertices[2].DstY = float32(10)
// ...
```

Next come the indices. The indices tell `DrawTriangles()` which vertices from the first argument must be used for each triangle. If we only have one three vertices and one triangle, the indices should be `[]uint16{0, 1, 2}`.

If you only draw one triangle, passing the indices can seem redundant. To really understand the point of the indices, we need to consider a more complex shape: a square, for example. Let's make a square of 10x10. The previous triangle we made covers half that, but we will need a second triangle with the vertices `(0, 0)`, `(10, 10)` and `(0, 10)` to cover the remaining bottom-left half. Here is where we notice that two of those vertices are already shared with the previous triangle! This means that we don't need to pass 6 vertices, only 4. Rectangles only have 4 vertices, so this makes sense. For our square, we only need to pass those 4 vertices and then adjust the indices so that they still point to 3 vertices per triangle. In this case, we would use the indices `[]uint16{0, 1, 2, 0, 2, 3}`.

The key idea is that having vertices and indices separated allows us to reuse vertices more efficiently. We can create strips of triangles that always share two of the three vertices. As a more extreme example, if we try to approximate a circle with triangles, the center of the circle will be a vertex shared by all the triangles in the polygon.

### Coloring the triangles

We have understood the geometry, but now we want to give color to the triangles. Should we use the `Vertex.Color*` fields, `DrawTrianglesOptions.ColorM`, the `*Image` passed as an argument..? Or all of them..? OR WHAT?!

It's normal to be confused:
- `DrawTrianglesOptions.ColorM` is a global color matrix affecting the final result. A global modifier. In general you don't need to use this.
- The `*Image` argument is a texture that can be used for the triangles. The problem is that the way this texture is applied to the triangles depends on the `Vertex.SrcX` and `Vertex.SrcY` fields. Each vertex can specify a different point to sample from the texture. In other words, the area from the texture used to "paint" a triangle is determined by the `SrcX` and `SrcY` fields of the three vertices that conform the triangle. This texture image must always be passed to `DrawTriangles()`, and it can't be nil. If you don't fully get this, don't get too hung up on it, there are visual examples later.
- The `Color*` fields in the `Vertex` structs are another color modifier. Be careful because the behavior is *different* for `DrawTriangles()` and `DrawTrianglesShader()`!
	- For `DrawTriangles()`, the color fields are *multipliers* for the color values sampled from the `*Image` texture. If `ColorA` is 0, the triangle will be fully transparent. If `ColorA` is 1 and all the other color fields are 0 (and the texture is opaque... because if it was transparent the triangle would still be transparent), then the triangle color will be black (`anyColorValue*0 = 0`).
	- For `DrawTrianglesShader()`, the vertex colors are only passed to the shader. The shader can decide what to do with them.

All this can be a bit difficult to grasp, so let's see it in practice:

![](https://github.com/tinne26/kage-desk/blob/main/img/triangle_A.png?raw=true)

In this image, we can see a single triangle with three vertices. The `Vertex.Dst*` coordinates are shown on the top-left. One point for each vertex. We also have a texture on the right, and the `Vertex.Src*` fields also define where to sample the colors of the texture for each vertex. Notice that the size of the texture can be different from the triangle or the canvas, and that the `Src*` positions must refer to the dimensions of the texture being used, not the game canvas.

The color values for the vertices are all set to 1. Since all the color multipliers are 1 (identity), the color of the texture will be respected. Let's see the next example to be able to compare the differences:

![](https://github.com/tinne26/kage-desk/blob/main/img/triangle_B.png?raw=true)

In this example, the position of the vertex A and the texture sampling point of the vertex B have both changed. This example should make it clearer that `Dst*` and `Src*` positions are independent. `Dst*` coordinates refer to the *geometry* of the triangles. `Src*` coordinates refer to the *texture sampling point* for a vertex. They work on different spaces. Notice that the color of the vertex B has changed in the result.

Let's look at a third example that uses the vertex color fields:

![](https://github.com/tinne26/kage-desk/blob/main/img/triangle_C.png?raw=true)

In this example all the texture sampling points are fairly close, which should result in a rather homogenous golden color for the triangle. But if you look at the vertex color values, you will see that the `ColorG` field (the green component) has been reduced to 0. This means that any green component in the texture will be cancelled. The golden in the texture has a lot of red, a large amount of green and then some blue. When the green channel is nullified, the remaining is mostly red and some blue. This results in the pink-red color you can see in the final render.

A very common and helpful trick is to use a single-pixel white texture and have all the `Src*` texture coordinates sampling the center of that pixel. Then, we can use the `Color*` attributes of the vertices to set a color manually. Since the texture is white, the vertex colors will be the only ones to determine the triangle's color.

The final example shows what happens if you nullify all colors. As you should expect, we get a black triangle:

![](https://github.com/tinne26/kage-desk/blob/main/img/triangle_D.png?raw=true)

You can try your own combinations interactively by running this program:
```
go run github.com/tinne26/kage-desk/examples/misc/triangles@latest
```

Summarizing:
- If you want to control the triangle colors exclusively with the texture image, set all `Color*` fields to 1 in the vertices. They will act as a neutral multiplier (identity).
- If you want to control the triangle colors manually with the `Color*` fields (useful for triangles with solid or simple colors), then use a 1x1 white texture instead and sample its center.
- You can still further transform the colors with `options.ColorM`.


## Differences for `DrawTrianglesShader()`

If you have understood `DrawTriangles()`, then `DrawTrianglesShader()` is almost the same, but:
- Instead of a texture, you use a shader program for the coloring.
- The vertex `Color*` fields do not directly affect anything, they are passed as the third argument of the shader's `Fragment(position vec4, _ vec2, color vec4)` entry function. At `position = vertex`, the input color will match the vertex color. Throughout the rest of the triangle, colors will be interpolated among the three vertices composing the triangle.




