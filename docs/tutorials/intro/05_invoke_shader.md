# Invoking shaders manually

While writing the first shaders we have been using `kage-desk/display` to keep our `main.go` file really simple. This is great to get started, but if we want to keep growing we need to take off the training wheels. We have learned how to make shaders, but we still haven't seen how to invoke them from Ebitengine by *ourselves*.

In theory, the easiest way to go would be to use [`ebiten.DrawRectShader(...)`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawRectShader). We won't be using that, but [`ebiten.DrawTrianglesShader(...)`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawTrianglesShader) instead. The reason is that `DrawTrianglesShader(...)` is more general and has less limitations. It might be a bit more confusing at first, but that's why you have this tutorial.

Here's a quick explanation of the function:
```Golang
func (*Image) DrawTrianglesShader(vertices []Vertex, indices []uint16, shader *Shader, options *DrawTrianglesShaderOptions)
```
- The `*Image` on which we call the `DrawTrianglesShader(...)` method is our *target image*, the image that we will be modifying or drawing upon.
- If we want our shader to be applied on the whole image, we would need to pass four `vertices` with `DstX` and `DstY` coordinates (destination coordinates) matching the corners of the `*Image.Bounds()`. If we want to operate only on a subregion of the image or an arbitrary polygon, we need to set the destination coordinates of our vertices accordingly. We might use from three to many vertices.
- The `indices` will tell the GPU how to create triangles from the given vertices. For example, if we pass vertices for our full target region, say `{top-left, top-right, bottom-left, bottom-right}`, then we will need at least two triangles. There are multiple ways to define these two triangles, but we could pick `{top-left, top-right, bottom-left}` and `{bottom-left, top-right, bottom-right}` to say something. This would mean we have to pass 6 indices: `{0, 1, 2, 2, 1, 3}`. The first three correspond to the first triangle (`{0, 1, 2}`), and the remaining to the second one `{2, 1, 3}`. This is effectively a form of mapping vertices to triangles.
- The `shader` is our Kage shader program.
- The `options` allow us to change the drawing blend mode and pass additional information to the shader, but we will explore this on later chapters, so it's better that you ignore it for the moment.

This is only a brief contextualization, it's ok if you don't fully understand everything yet... but if you are completely lost or want a more in depth explanation, there's a separate article explaining [triangles](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/triangles.md) in Ebitengine in more detail.

With all that out of the way, let's rework our `main.go`:
```Golang
package main

import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { panic(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	ebiten.SetWindowTitle("intro/invoke-shader")
	ebiten.SetWindowSize(512, 512)
	err = ebiten.RunGame(game)
	if err != nil { panic(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader *ebiten.Shader
	vertices [4]ebiten.Vertex
}

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	return 512, 512
}

// No logic to update.
func (self *Game) Update() error { return nil }

// Core drawing function from where we call DrawTrianglesShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// map the vertices to the target image
	bounds := screen.Bounds()
	self.vertices[0].DstX = float32(bounds.Min.X) // top-left
	self.vertices[0].DstY = float32(bounds.Min.Y) // top-left
	self.vertices[1].DstX = float32(bounds.Max.X) // top-right
	self.vertices[1].DstY = float32(bounds.Min.Y) // top-right
	self.vertices[2].DstX = float32(bounds.Min.X) // bottom-left
	self.vertices[2].DstY = float32(bounds.Max.Y) // bottom-left
	self.vertices[3].DstX = float32(bounds.Max.X) // bottom-right
	self.vertices[3].DstY = float32(bounds.Max.Y) // bottom-right
	// [VERTEX-NOTE]
	// Other properties will be set on later examples. The full
	// configuration is quite verbose, but you will typically create
	// your own helper functions to do the heavy lifting, and in
	// some cases you can optimize and omit some settings on
	// successive passes.

	// triangle shader options
	// (we are not setting any properties here, but you can do
	// it when needed. it's also possible to reuse the options
	// in many cases, here we are keeping it straight-forward)
	var shaderOpts ebiten.DrawTrianglesShaderOptions

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```

It takes more code than before, but most of it is boilerplate that you should already be familiar with:
- We use the [`embed`](https://pkg.go.dev/embed) package in order to store the shader program data directly into the compiled executable. While you could also use `os.ReadFile()` instead of the `go:embed` macro to obtain the shader's source code, using `embed` is generally recommended because it also works on browsers and mobile.
- In the `main` function, we compile the shader with `NewShader()` and create a minimal `Game` struct that can display it.
- The `Game` struct itself only keeps a reference to the compiled shader program, a few vertices, and draws using `DrawTrianglesShader()` on each frame.
- The last point to keep in mind is that we are using a fixed layout of 512x512 at all times, which is kinda hacky and not representative of what you will be using in most real projects. If you are still confused about how `Game.Layout()` works in Ebitengine, make sure to [revise the basics](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/ebitengine_game.md). Shaders in non-toy applications almost always have to be resizable, so depending on hardcoded positions and resolutions is something we will have to change in future chapters.

> [!NOTE]
> *If you are having any trouble, the full code for this manual shader invocation (along with a cleaned-up version of the wave shader of the previous chapter) can be found at [`kage-desk/examples/intro/invoke-shader`](https://github.com/tinne26/kage-desk/blob/main/examples/intro/invoke-shader).*


### Table of Contents
Next up: [#6](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [**Manual shader invocation**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
