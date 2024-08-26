# Built-in functions

The half-white half-black shader was too easy, so let's make it more challenging now. For example, if I asked you to make a shader that creates a checkerboard-like pattern by making a pixel white if it's on an even position, and black otherwise... would you be able to make it?

Well, with what you have been taught up to this point, you won't... because you don't know how to check if a number is even or odd in a Kage shader. To solve this, it's time to introduce the *built-in functions*. Although there's a [full list](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(mathematics)) at ebitengine.org, you don't need to go there yet; learning about just a few of the functions is enough for the moment:
- **Single-argument functions**: `abs`, `sign` (returns -1, 0 or 1), `sin`, `cos`, `sqrt`, `floor`, `ceil`, `fract` (returns the fractional part of a number), `length` (mathematical length of a vector) and `len` (same as in Golang, but applied to `vec` types, there are no slices or maps in Kage).
- **Two-argument functions**: `mod(x, m)` (a.k.a `%`), `min(a, b)`, `max(a, b)`, `pow(x, exp)`, `step(s, x)` (0 if `x < s`, 1 otherwise), `distance(pointA, pointB)`.
- **Three-argument functions**: `clamp(x, min, max)`, `mix(a, b, t)` (linear interpolation).

These are all built-in functions that you can use in the same way you would use `copy` or `make` on Golang. They also happen to work on multiple data types: `float`, `vec2`, `vec3`, `vec4`.

With this new superpower, implement the shader we mentioned earlier: even pixels white, odd pixels black. Go!

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	xy := floor(targetCoords.x) + floor(targetCoords.y)
	if mod(xy, 2) == 0 {
		return vec4(1) // white
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```
*(Full program available at [examples/intro/checkerboard-preview](https://github.com/tinne26/kage-desk/blob/main/examples/intro/checkerboard-preview))*

The key function for this program is `mod()`. Using `mod()` allows us to find whether the position corresponds to an even or odd pixel. Using `floor()` is not strictly necessary, but since positions are given at the center of the pixel, using floor gets rid of the 0.5 decimal part and gives us a more natural value to work with when we only want to check the "parity of the pixel index".
</details>

If your screen has any form of DPI scaling, you might see the result heavily aliased. What would you do to make the pattern bigger and easier to see? Like, alternating 8x8 pixels white and 8x8 pixels black? (You can also try with 2x2 if that makes it easier to think about the problem). This one is a bit trickier, but that's the kind of math problems you will often face with shaders. Think about it for a while!

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

const CellSize = 8

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	xy := floor(targetCoords.x/CellSize) + floor(targetCoords.y/CellSize)
	if mod(xy, 2) == 0 {
		return vec4(1) // white
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```
*(Full program available at [examples/intro/checkerboard](https://github.com/tinne26/kage-desk/blob/main/examples/intro/checkerboard))*

We have declared a constant here to make our life easier and show that you can totally use that, but if you hardcoded the values that's perfectly fine too. With this, it's easier to change the cell size and make it 32 or anything else.
</details>

If you got it, well done! If you didn't, don't worry. There's a common idea here that's worth explaining as it comes up all the time when writing shaders: conceptually, we wanted to do the same as in the previous example... but at a different scale. We wanted to project the original "canvas" to one that was a fraction of the size, and only then apply the same `mod()` function as we did earlier. This idea of scaling / projecting / deforming a space or surface is extremely common in shaders. It may be confusing at the beginning, but try to wrap your mind around it. In general, when you have to work at the context of a single pixel, mathematical transformations are a very powerful tool.

---

Let's close this lesson with a final challenge. Remember the half-white half-black screen from the previous chapter? Now try to make the split be wavy instead of a straight line, using `sin` or `cos`. Here's my take at a resolution of 512x512:

![](https://github.com/tinne26/kage-desk/blob/main/img/intro_gpu_wave.png?raw=true)

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	const Pi = 3.14159265
	const NumOscillations = 7.0
	const WaveWidth = 18.0

	waveFactor := sin((targetCoords.y/512.0)*2*Pi*NumOscillations)*(WaveWidth/2)
	if targetCoords.x < 256 + waveFactor {
		return vec4(1) // white
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```
*(Full program available at [examples/intro/wave-split](https://github.com/tinne26/kage-desk/blob/main/examples/intro/wave-split))*

Let's break it down a bit: `sin()` expects an angle in radians. There are `2*Pi` radians in a circumference. Therefore, if we can go from `0` to `2*Pi` through the 512 pixels of the screen's vertical axis, we will have completed a full sine oscillation. We want more oscillations? Just multiply the value going into `sin()` by `NumOscillations`! Finally, we can also control the width of the sine wave by multiplying the `sin` result by `WaveWidth/2`. Since the result is already oscillating between `[-N, +N]`, we only need to add this `waveFactor` to our previous cutoff point... and now we have a fancy sine wave splitting the screen instead of a boring straight line!
</details>

You might have noticed that the edge of the sine wave is jaggy, not smooth. We will see how to improve that in later examples, so don't get too hung up on it for the time being.

If this was a bit difficult don't worry. The most important part is that you get exposed to these ideas and slowly become used to them. There are other parts of shaders that don't revolve so much around maths, but it's important to get some practice and become more familiar with these techniques if you really want to get the most out of Kage.


### Table of Contents
Next up: [#5](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [**Built-in functions**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
