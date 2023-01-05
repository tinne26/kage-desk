## Built-in functions

The half-white half-black shader was too easy, so let's complicate it a bit more now. For example, if I asked you to make a shader that creates a pattern by making a pixel white if it's on an even position, and black otherwise... would you be able to make it?

Well, with what you have been explained up to this point, you won't. Because you don't know how to check if a number is even or odd in a Kage shader: it's time to introduce the *built-in* functions. Although there's a [full list](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(mathematics)) at the official ebitengine.org page for shaders, I don't recommend going there yet. Just know it exists. Instead, I'll share the most common ones here:
- Single argument functions: `abs`, `sign` (returns -1, 0 or 1), `sin`, `cos`, `sqrt`, `floor`, `ceil`, `fract` (returns the fractional part of a number), `length` (length of a vector) and `len` (same as in Golang, but applied to `vec` types, there are no slices or maps in Kage).
- Two-argument functions: `mod(x, m)` (`%`), `min(a, b)`, `max(a, b)`, `pow(x, exp)`, `step(s, x)` (0 if `x < s`, 1 otherwise), `distance(pointA, pointB)`.
- Three-argument functions: `clamp(x, min, max)`, `mix(a, b, t)` (linear interpolation).

These are all built-in functions that you can use just like you would use `copy` or `make` on Golang. They also happen to work on multiple data types: `float`, `vec2`, `vec3`, `vec4`.

With this new superpower, implement the shader we mentioned earlier: even pixels white, odd pixels black. Go!

This is a possible answer:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	xy := position.x + position.y
	if mod(xy, 2) == 0 {
		return vec4(1) // white
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```

If your screen has any form of DPI scaling, you may see the result somewhat aliased. What would you do to make the pattern bigger and easier to see? Like, alternating 2x2 pixels white and 2x2 pixels black? This one is a bit trickier, but that's the kind of math problems you will often find with shaders. Think about it for a while!

Here's my solution:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	xy := floor(position.x/2) + floor(position.y/2)
	if mod(xy, 2) == 0 {
		return vec4(1) // white
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```

You can also divide by 4, for example, to see the result even more clearly. If you got it, well done! If you didn't, don't worry. There's a common idea here that's worth explaining as it comes up all the time when writing shaders: conceptually, we wanted to do the same as in the previous example... but at a different scale. We wanted to project the original "canvas" to one that was half the size, and then apply the same `mod` function as we had been doing. This idea of scaling / projecting / deforming a space or surface is extremely common in shaders. It may be confusing at the beginning, but try to wrap your mind around it. In general, when you have to work at the context of a single pixel, mathematical transformations are a very powerful tool.

Let's go with another challenge. Remember the half-white half-black screen? Try now to make the split be wavy instead of a straight line using `sin` or `cos`.

I kinda liked this version (you can use the environment variable `EBITENGINE_SCREENSHOT_KEY=q` to take screenshots):

![](https://github.com/tinne26/kage-desk/blob/main/img/intro_gpu_wave.png?raw=true)

My code looked like this:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	const Pi = 3.14159265
	const NumOscillations = 7.0
	const WaveWidth = 18.0

	waveFactor := sin((position.y/511.0)*2*Pi*NumOscillations)*(WaveWidth/2)
	if position.x < 256 + waveFactor {
		return vec4(1) // white
	} else {
		return vec4(0, 0, 0, 1) // black
	}
}
```

Let's break it down a bit: `sin` expects a value in radians. There are `2*Pi` radians in a circumference. Therefore, if we can go from `0` to `2*Pi` through the 512 pixels of height that the screen has, we will have completed a full sine oscillation. We want more oscillations? Just multiply the value going into `sin` by `NumOscillations`! Finally, we can also control the width of the sine wave by multiplying the `sin` result by `WaveWidth/2`. Since the result is already oscillating between `[-N, +N]`, we only need to add this `waveFactor` to our previous cutoff point... and now we have a fancy sine wave splitting the screen instead of a plain straight line.

You may have noticed that the edge of the sine wave is jaggy, not smooth. We will see how to improve that in later examples, so don't get too hung up on it for the time being.

If this was a bit difficult don't worry. The most important part is that you get exposed to these ideas and slowly become used to them. There are other parts of shaders that don't revolve so much around maths, but it's important to get some practice and become more familiar with these techniques if you really want to get good at them.

Let's summarize before the next section:
- Without the input position we can't do much cool stuff (yet).
- We have many useful built-in functions to help us out, like `mod(x, m)`.
- Maths and projections are a very important tool in the GPU realm.


### Table of Contents
Next up: [#5](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md).

0. [Main](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/00_main.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/03_position_input.md)
4. [**Built-in functions**](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/07_images.md)
8. [Screen vs sprite effects]()
9. [Performance considerations]()
10. [Graduation challenges]()
