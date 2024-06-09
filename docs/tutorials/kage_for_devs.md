# Kage for game devs

Quick reference for developers with previous knowledge of how shaders work.

Kage is the language used on Ebitengine to write shaders. It has a Golang-like syntax and it's internally translated to GLSL, HSL or MSL at runtime as required. Only fragment shaders are supported at the moment.
- Supported types: `bool`, `int`, `float` (float32), `vec2`, `vec3`, `vec4` (float vectors), `mat2`, `mat3`, `mat4` and constants.
- Vectors support swizzling with `rgba`, `xyzw` and `stpq`. You can also index them `[N]` directly.
- You can write helper functions, but there are no slices, maps, strings, structs, import, switch, etc.
- Source textures limited to 4 RGBA images per shader invocation.

#### Table of contents
- [Built-in functions](#built-in-functions)
- [Basic example](#basic-shaderkage-example)
- [Load and invoke](#load-and-invoke)
- [Uniforms](#uniforms)
- [Textures](#textures)

## Built-in functions

Almost all of them can be applied both to numerical types like `int` and `float`, but also vectors. E.g.: `abs(vec2(-1, 0)) == vec2(1, 0)`.

Key single-argument functions:
```Golang
len(vec) // for vec2, vec3, vec4. same as in Golang
length(x) // mathematical length / magnitude of a vector
abs(x)
sign(x) // returns -1, 0 or 1
sin(x); cos(x); tan(x) // plus asin, acos, atan, etc.
sqrt(x)
floor(x); ceil(x); fract(x)
```

Key two-argument functions:
```Golang
mod(x, m) // %
min(a, b); max(a, b)
pow(x, exp)
step(s, x) // 0 if `x < s`, 1 otherwise
atan2(x, y) // classic `angle := atan2(x - ox, y - oy)`
dot(x, y); cross(x, y vec3) // dot and cross products
distance(pointA, pointB) // == length(pointA - pointB)
```

Key three-argument functions:
```Golang
clamp(x, min, max)
mix(a, b, t) // lerp, linear interpolation
```

Full [official reference](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(mathematics)).

## Basic `shader.kage` example

Example Kage shader that generates a checkerboard pattern.
```Golang
//kage:unit pixels
package main

func Fragment(targetCoords vec4, sourceCoords vec2, color vec4) vec4 {
	const CellSize = 32

	cellCoords := floor(targetCoords/CellSize)
	if mod(cellCoords.x + cellCoords.y, 2) == 0 {
		return vec4(1, 0, 1, 1) // magenta
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```
*(Make sure to [configure your editor](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/config_editor.md) so you get syntax highlight)*

For quick testing, you can put the code into a `shader.kage` file and run it with the following `main.go`:
```Golang
package main

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("kage/checkerboard")
	display.SetSize(512, 512)
	display.Shader("shader.kage")
}
```

The `//kage:unit pixels` is a special directive similar to Golang's compiler directives. Without it, you would be operating in texels mode. The general recommendation —and what all advanced users I know of are doing— is to use the new pixels mode. All the information in this tutorial is based on the pixel mode; if there's anything that behaves differently under the texels mode, I won't even bother telling you, so keep that in mind.

> [!TIP]
> *If you are too lazy, you can also use [www.kageland.com](https://www.kageland.com/), a very handy playground created by [@tomlister](https://github.com/tomlister) to write and share Kage shaders from your browser. Try copy pasting the shader code above and `Run` it!*

## Load and invoke

If you want to compile and invoke a shader manually, here is some reference code. Basically, use 4 vertices to create a quad and set the vertex target coordinates. While `DrawRectShader()` also exists, I recommend focusing on `DrawTrianglesShader()` instead[^1].

We will build upon the following template for the next examples:

[^1]: While `DrawRectShader()` can be handy in some situations, in many practical scenarios you will have to end up reaching for `DrawTrianglesShader()` anyway, so I'm of the opinion that you should just spare yourself the cognitive overhead and simply ignore the function. If you really must know, its biggest restriction is probably that all source images must match the explicit dimensions passed as the first arguments.

```Golang
package main

import "time"
import _ "embed"
import "github.com/hajimehoshi/ebiten/v2"

//go:embed shader.kage
var shaderProgram []byte

func main() {
	// compile the shader
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil { panic(err) }

	// create game struct
	game := &Game{ shader: shader, startTime: time.Now() }

	// configure window and run game
	ebiten.SetWindowTitle("kage/load-and-invoke")
	ebiten.SetWindowSize(512, 512)
	err = ebiten.RunGame(game)
	if err != nil { panic(err) }
}

// Struct implementing the ebiten.Game interface.
// Reusing the vertices and options is advisable.
type Game struct {
	shader *ebiten.Shader
	vertices [4]ebiten.Vertex
	shaderOpts ebiten.DrawTrianglesShaderOptions
	startTime time.Time
}

func (self *Game) Update() error { return nil }
func (self *Game) Layout(_, _ int) (int, int) {
	return 512, 512 // fixed layout
}

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

	// NOTE: here we will also map the vertices to 
	//       the source image in later examples.

	// triangle shader options
	if self.shaderOpts.Uniforms == nil {
		// initialize uniforms if necessary
		self.shaderOpts.Uniforms = make(map[string]any, 2)
		self.shaderOpts.Uniforms["Center"] = []float32{
			float32(screen.Bounds().Dx())/2,
			float32(screen.Bounds().Dy())/2,
		} // this will be passed as a vec2

		// link images if necessary (omit if nil)
		self.shaderOpts.Images[0] = nil
		self.shaderOpts.Images[1] = nil
		self.shaderOpts.Images[2] = nil
		self.shaderOpts.Images[3] = nil
	}

	// additional uniforms
	seconds := float32(time.Now().Sub(self.startTime).Seconds())
	self.shaderOpts.Uniforms["Time"] = seconds
	
	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &self.shaderOpts)
}
```

> [!TIP]
> *If you are too lazy to do all this when starting but you still prefer your editor to [kageland](https://www.kageland.com/), notice that the [`kage-desk/display`](https://pkg.go.dev/github.com/tinne26/kage-desk/display) package also provides many utilities: quick setup, background color, two default images, high resolution and resizability options, `Time float`, `Cursor vec2` and `MouseButtons int` uniforms, F (fullscreen) and ESC shortcuts... If you are interested, check out the [main.go](https://github.com/tinne26/kage-desk/blob/main/examples/learn/filled-circle/main.go) and [shade.kage](https://github.com/tinne26/kage-desk/blob/main/examples/learn/filled-circle/shader.kage) files of the `learn/filled-circle` example for a quick reference.*

## Uniforms

The code from the previous section already shows how to link uniforms from the CPU side. Now let's see how to use them in the actual shader. We are going to make a shader where a pixel orbits around the center of the screen, at a rate of one revolution per minute (so, a clock that tracks seconds):

```Golang
//kage:unit pixels
package main

var Center vec2 // technically this isn't necessary, will explain later
var Time float 

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	const MarkerDistance = 160
	const Pi = 3.14159265

	// compute the position for the seconds marker
	secAngle := (mod(Time, 60)/60)*2*Pi
	secPos   := Center + vec2(sin(secAngle)*MarkerDistance, -cos(secAngle)*MarkerDistance)

	// return the sum of contributions for the two dots in the screen
	centerMarker := vec4(1)*inDotMask(targetCoords.xy, Center, 2, 1.5)
	secondMarker := vec4(1)*inDotMask(targetCoords.xy, secPos, 2, 1.5)
	return centerMarker + secondMarker
}

// Returns 1 if the current position is within 'hardRadius' of 'target',
// between 1 and 0 if within 'hardRadius + softRadius', zero otherwise.
func inDotMask(current vec2, target vec2, hardRadius, softRadius float) float {
	return 1.0 - smoothstep(hardRadius, hardRadius + softRadius, distance(current, target))
}
```

As you can see, adding uniforms is as simple as declaring exported variables at the start of the file with the correct names. Vectorial types are inferred from `[]float32` slices. Implicit conversions from `float64` and `int` to the shader's `float` will also happen automatically, but using `float32` directly on the CPU side is probably better practice.

## Textures

To sample a texture, you will typically use the `sourceCoords` input argument and the `imageSrc0At()` function, which expects a coordinate in pixels. As showcased in the [load and invoke](#load-and-invoke) section, you can link up to 4 images in the shader options. The full collection of relevant image functions is the following:
- <code>imageSrc<b><i>N</i></b>At()</code> (replace N with {0, 1, 2, 3}): source texture sampling. Sampling is [always nearest](https://github.com/hajimehoshi/ebiten/issues/2962); if you want to perform linear interpolation you will have to do it manually. Until we make something better, you can find some interpolation implementations [here](https://github.com/tinne26/mipix/tree/main/filters)[^2].
- <code>imageSrc<b><i>N</i></b>UnsafeAt()</code> (replace N with {0, 1, 2, 3}): like <code>imageSrc<b><i>N</i></b>At()</code>, but doesn't check whether you go out of bounds. If you go out of bounds with <code>imageSrc<b><i>N</i></b>At()</code>, you will get back `vec4(0)`. With the unsafe function, you could actually peek at the whole internal atlas.
- <code>imageSrc<b><i>N</i></b>Size()</code> (replace N with {0, 1, 2, 3}): returns the size in pixels of the requested source texture.
- <code>imageSrc<b><i>N</i></b>Origin()</code> (replace N with {0, 1, 2, 3}): returns the origin of the requested source texture in pixels. This is relevant when working with subimages that might not start at (0, 0). Always keep those in mind!
- `imageDstSize()` and `imageDstOrigin()`: same idea as the two previous functions, but for the target texture instead of the sources. With these you could eliminate the `Center` uniform of the previous section, for example.

[^2]: Many of those shaders use `SourceRelativeTextureUnitX` and `SourceRelativeTextureUnitY` uniforms, but that can often be replaced with `units := fwidth(sourceCoords)`. Otherwise, the classic bilinear interpolation with +/-0.5 can be found at [src_bilinear.kage](https://github.com/tinne26/mipix/blob/main/filters/src_bilinear.kage) and doesn't require any uniforms. This is pretty much what Ebitengine does by default with `FilterLinear`, but with some extra clamping that you might or might not be interested in.

> [!NOTE]
> *Remember that colors are in RGBA format, with values between `[0, 1]`, and [premultiplied alpha](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/premult.md) (color channel values can't exceed the alpha value, or weird stuff will happen). It's easy to slip when you are often working with [0, 255] RGBA on the CPU side.*

If you are using `DrawTrianglesShader(...)`, you also need to *map the source texture to the target vertices*:
```Golang
// (typically done after the DstX/DstY setup)
srcBounds := yourImage.Bounds()
self.vertices[0].SrcX = float32(srcBounds.Min.X) // top-left
self.vertices[0].SrcY = float32(srcBounds.Min.Y) // top-left
self.vertices[1].SrcX = float32(srcBounds.Max.X) // top-right
self.vertices[1].SrcY = float32(srcBounds.Min.Y) // top-right
self.vertices[2].SrcX = float32(srcBounds.Min.X) // bottom-left
self.vertices[2].SrcY = float32(srcBounds.Max.Y) // bottom-left
self.vertices[3].SrcX = float32(srcBounds.Max.X) // bottom-right
self.vertices[3].SrcY = float32(srcBounds.Max.Y) // bottom-right
```

Don't forget to link your images too!
```Golang
self.shaderOpts.Images[0] = yourImage
```

If you need texture clamping / repeat, the [snippets](https://github.com/tinne26/kage-desk/blob/main/docs/snippets/snippets.md#texture-clamping-and-repeat) page has some code for it.
