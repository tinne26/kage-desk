# More input: uniforms

Shaders, being run in the GPU as they are, are rather limited in regards to the information they can access. Until now, the only real input we have seen for them is the `position vec4` input parameter that we receive on the `Fragment()` entry function.

So... can we pass additional parameters to the shaders?

The answer to this question is yes: **uniforms** are variables whose values can be set from your code running on the CPU.

In order to show how to use uniforms, we will try to draw a filled circle using a shader now. Our draw function looked like this in the last chapter:
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.GeoM.Translate(0, 0) // you could adjust the drawing position here
	
	// draw shader
	screen.DrawRectShader(512, 512, self.shader, opts)
}
```

Now we will change it a bit to add a new parameter for our shader: the center position of our circle.
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx())/2,
		float32(screen.Bounds().Dy())/2,
	}
	
	// draw shader
	screen.DrawRectShader(512, 512, self.shader, opts)
}
```

This `Uniforms` map allows us to send the `Center` variable to our shader. The main types that we can send are `int`, `float32` and `[]float32` values, which correspond to the shader `int`, `float` and `vec*` types. Since we passed a slice of two values, our shader has to look like this:
```Golang
package main

var Center vec2 // uniform: circumference center coords

func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	// ...
}
```

Uniform variables must always be capitalized (exported variables) and appear at the start of our shader code.

Try to complete the shader so it draws a filled circle. Use a radius of 80px and any color you want.

<details>
<summary>Click to show the solution</summary>

```Golang
package main

var Center vec2 // uniform: circumference center coords
const Radius = 80.0

func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	distToCenter := distance(Center, position.xy)
	distToEdge   := distToCenter - Radius

	// dist to edge will be negative if we are inside the
	// circle and positive if we are outside, but we want to 
	// preserve the circle color if we are inside (multiply
	// by one), and discard it if we are outside (multiply
	// by zero), so we need to change the sign and clamp
	factor := clamp(-distToEdge, 0, 1)
	factor  = pow(factor, 1.0/2.2) // gamma correction
	return vec4(1, 0, 0, 1)*factor
}
```

If you used `if` statements instead of clamp and don't know what gamma correction is, don't worry. The reason clamp (or some combinations of min/max) are preferred to conditionals is that shaders are executed by many GPU processors in parallel, and typically they are all executing the same instruction at the same time. When there are branches, all branches may have to be executed for all processors anyway. The topic is deep and complex and it's not something you have to worry about right now, but it's good to start seeing ways to avoid conditionals. Here an actual conditional wouldn't be much worse, but you definitely don't want big conditionals doing completely different things, because you may end up having to execute all those big branches on all processors anyway.

On the other topic of gamma correction, the issue is that lightness is not perceived linearly by humans, but follows a power function instead. Therefore, using a linear fall-off for the opacity at the edge of the circumference is not ideal, so... we can use a simple formula to correct it. Again, this doesn't matter much here, but it's a concept you may want to know about for your future adventures. In more complex shaders it can have a significant effect.
</details>

Modify the `main.go` and the `shader.kage` programs now so the `Radius` also becomes a uniform. Pass the value of 80 from the `Draw()` function in `main.go` instead of hardcoding it in the shader.

Now the draw function in `main.go` should look like this:
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx())/2,
		float32(screen.Bounds().Dy())/2,
	}
	opts.Uniforms["Radius"] = float32(80.0)
	
	// draw shader
	screen.DrawRectShader(512, 512, self.shader, opts)
}
```

And the `shader.kage` should have `const Radius = 80.0` replaced by `var Radius float`.

Again we made it too easy? Ok, ok... then, here's a challenge. Try to do this:



Basically, we will have the same we had in the latest shader, but animating the circumference radius.

<details>
<summary>Click to show hints and link to the solution</summary>

Spoilers on how to solve this problem: add an `angle int` variable to the `Game` struct. You want its value to go from 0 to 359 and back again to zero at a rate of 1 degree per tick. Then, on `Draw()`, you can use a radius of `80 + 30*someOscillatingFactor`, where the factor oscillates between `[-1, 1]` and is derived from the `angle`.

The actual code can be found at [`kage-desk/examples/intro/circle-anim`](https://github.com/tinne26/kage-desk/blob/main/examples/intro/circle-anim).
</details>

You are now able to pass your own parameters to the shaders and we have even seen how to use this to create an animated effect. If you start having ideas to try on your own, now is the time! Feel free to experiment by yourself, you will only learn shaders by writing shaders.


### Table of Contents
Next up: [#7](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/07_images.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md)
6. [**More input: uniforms**](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/07_images.md)
8. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/08_loops.md)
