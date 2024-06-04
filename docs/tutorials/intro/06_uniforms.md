# More input: uniforms

Shaders, being run on the GPU, are rather limited in regards to the information they can access. Up to this point, the only real input we have seen for them is the `targetCoords vec4` input parameter that we receive on the `Fragment()` entry function... Which begs the question: can we pass additional parameters to the shaders?

The answer is yes: **uniforms** are variables whose values can be set from your CPU-side code and sent to the GPU for use with your shader.

In order to show how to use uniforms, we will try to draw a filled circle using a shader. Our draw function was something like this in the previous chapter:
```Golang

func (self *Game) Draw(screen *ebiten.Image) {
	// ... (some stuff)

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```

Now we will use the `shaderOpts` to add a new parameter for our shader: the center position of our circle.
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// ... (some stuff)

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Uniforms = make(map[string]interface{})
	shaderOpts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx())/2,
		float32(screen.Bounds().Dy())/2,
	}

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```

The `Uniforms` map on the `DrawTrianglesShaderOptions` allows us to send the `Center` variable to our shader. The main types that we can send are `int`, `float32` and `[]float32` values, which correspond to the shader `int`, `float` and `vec*` types. Since we passed a slice of two values, our shader has to look like this:
```Golang
//kage:unit pixels
package main

var Center vec2 // uniform: circle center coords

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	// ...
}
```

Uniform variables must always be capitalized (like exported variables) and appear at the start of our shader code.

Try to complete the shader so it draws a filled circle. Use a radius of 80px and any color you want.

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

var Center vec2 // uniform: circle center coords
const Radius = 80.0

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	distToCenter := distance(Center, targetCoords.xy)
	distToEdge   := distToCenter - Radius

	// dist to edge will be negative if we are inside the
	// circle and positive if we are outside, but we want to 
	// preserve the circle color if we are inside (multiply
	// by one), and discard it if we are outside (multiply
	// by zero), so we need to change the sign and clamp
	factor := clamp(-distToEdge, 0, 1)
	return vec4(1, 0, 0, 1)*factor
}
```
*(Full program available at [examples/intro/circle](https://github.com/tinne26/kage-desk/blob/main/examples/intro/circle))*

If you used `if` statements instead of `clamp()`, don't worry. The reason `clamp()` (or some combinations of `min()`/`max()`) are preferred to conditionals is that shaders are executed by many GPU processors in parallel, and typically they are all executing the same instruction at the same time. When there are branches, all branches may have to be executed for all processors anyway. The topic is deep and complex and it's not something you have to worry about right now, but it's good to start seeing ways to avoid conditionals. Here an actual conditional wouldn't be much worse, but you definitely don't want big conditionals doing completely different things, because you may end up having to execute all those big branches on all processors anyway.
</details>

With the circle shader working, the next step is to modify the `main.go` and the `shader.kage` programs so the `Radius` also becomes a uniform. Pass the value of 80 from the `Draw()` function in `main.go` instead of hardcoding it in the shader.

Now the draw function in `main.go` should look like this:
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// ... (some stuff)

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Uniforms = make(map[string]interface{})
	shaderOpts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx())/2,
		float32(screen.Bounds().Dy())/2,
	}
	shaderOpts.Uniforms["Radius"] = float32(80.0)

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```

And the `shader.kage` should have `const Radius = 80.0` replaced by `var Radius float`.

Now we are ready for the final challenge:

https://user-images.githubusercontent.com/95440833/212187119-f4445602-a534-47c0-822a-1c3a4bd5de14.mp4

Your goal is to use uniforms in order to create an animation by changing the radius of the circle at each frame. Try to write this shader yourself!

<details>
<summary>Click to show hints and link to the solution</summary>

To solve this problem you can add an `angle int` variable to the `Game` struct. You want its value to go from 0 to 359 and back again to zero at a rate of 1 degree per tick. Then, on `Draw()`, you can use a radius of `80 + 30*someOscillatingFactor`, where the factor oscillates between `[-1, 1]` and is derived from the `angle`.

The full code of a working solution can be found at [`kage-desk/examples/intro/circle-anim`](https://github.com/tinne26/kage-desk/blob/main/examples/intro/circle-anim).
</details>

You are now able to pass your own parameters to the shaders and we have even seen how to use this to create an animated effect. If you start having ideas to try on your own, now is the time! Feel free to experiment by yourself, you will only learn shaders by writing shaders.


### Table of Contents
Next up: [#7](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [**More input: uniforms**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
