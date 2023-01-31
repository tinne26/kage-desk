# Invoking shaders manually

While writing the first shaders we have been using `kage-desk/display` to keep our `main.go` file really simple. This is great to get started, but if we want to keep growing we need to take off the training wheels. We have learned how to make shaders, but we still haven't seen how to invoke them from Ebitengine by *ourselves*.

Let's start by reworking our `main.go` with [`ebiten.DrawRectShader(...)`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawRectShader), which is the simplest call that can be used to invoke a shader in Ebitengine:
```Golang
package main

import "log"
import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { log.Fatal(err) }

	// create game struct
	game := &Game{ shader: shader }

	// configure window and run game
	ebiten.SetWindowTitle("intro/invoke-shader")
	ebiten.SetWindowSize(512, 512)
	err = ebiten.RunGame(game)
	if err != nil { log.Fatal(err) }
}

// Struct implementing the ebiten.Game interface.
type Game struct {
	shader *ebiten.Shader
}

// Assume a fixed layout.
func (self *Game) Layout(_, _ int) (int, int) {
	return 512, 512
}

// No logic to update.
func (self *Game) Update() error { return nil }

// Core drawing function from where we call DrawRectShader.
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.GeoM.Translate(0, 0) // you could adjust the drawing position here
	
	// draw shader
	screen.DrawRectShader(512, 512, self.shader, opts)
}
```

It takes more code than before, but most of it is boilerplate that you should already be familiar with:
- We use the [`embed`](https://pkg.go.dev/embed) package in order to store the shader program data directly into the compiled executable. While you could also use `os.ReadFile()` instead of the `go:embed` macro to obtain the shader's source code, using `embed` is recommended because it also works on browsers and mobile.
- In the `main` function, we compile the shader with `NewShader()` and create a minimal `Game` struct that can display it.
- The `Game` struct itself only keeps a reference to the compiled shader program and draws it with `DrawRectShader()` on each frame. The arguments of `DrawRectShader()` are the width and height of the area the shader will be drawn to, the compiled shader object and the draw options. These options include a `GeoM` matrix, as shown in the code, but also other options that will be explained in the next chapter.
- The last point to keep in mind is that we are using a fixed layout of 512x512 at all times. If you are still confused about how `Game.Layout()` works in Ebitengine, make sure to [revise the basics](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/ebitengine_game.md). Shaders in non-toy applications will often be resizable, so you will eventually have to learn to adjust the position or resolution of the shaders dynamically.

> If you are having any trouble, the full code for this manual shader invocation (along with a cleaned-up version of the wave shader of the previous chapter) can be found at [`kage-desk/examples/intro/invoke-shader`](https://github.com/tinne26/kage-desk/blob/main/examples/intro/invoke-shader).

There's another way to invoke a shader, using `DrawTrianglesShader()` instead of `DrawRectShader()`. We will explain this second method later in the tutorial... but only when it becomes necessary.


### Table of Contents
Next up: [#6](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [**Manual shader invocation**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [`DrawTrianglesShader()`](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_triangles.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
