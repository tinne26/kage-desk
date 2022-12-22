# Introduction to Kage

Kage is a programming language to implement shaders in Ebitengine.

If you are new to shaders, the short version is that they are programs that run on the GPU instead of the CPU. For our purposes, shaders are programs that allow us to create or modify images[^1], like recoloring them, adding noise or grain, distorting them and many others. Shaders are programs that allow us to perform sophisticated computations for the individual pixels of an image in a highly parallel manner.

[^1]: Shaders can also be used for general computation, not just graphics, but that's outside the scope of this guide. We have thrown a few links [here](https://github.com/tinne26/kage-desk/blob/main/tutorials/general_links.md) if you want to learn more later.

There are a few different languages in which shaders can be written: you may have heard of GLSL, HLSL and others. Ebitengine has it's own intermediate language, Kage, which allows us to write shaders in a Golang-like syntax and forget about the rest. At runtime, Ebitengine will translate that Kage program to HLSL, MSL or whatever language is needed to make it work on the platform where the game's being run.

#### Table of Contents  
1. [CPU vs GPU](#cpu-vs-gpu-a-change-of-mindset)  
2. [The first Kage program](#the-first-kage-programs)
3. [The `position` input parameter](#the-position-input-parameter)
4. [Built-in functions](#built-in-functions)
5. [The hidden infrastructure](#the-hidden-infrastructure)
6. [More input: uniforms](#more-input-uniforms)
7. [Using a texture](#using-a-texture)
8. [Screen vs sprites](#screen-effects-vs-sprite-effects)


## CPU vs GPU: a change of mindset

While we would love to jump right away into Kage, understanding in which ways CPU programs are different from GPU programs will allow us to move forward much faster later. When are shaders better than their CPU equivalents... and why?

The key concept behind shaders is **parallelization**. GPUs don't solve one big problem sequentially, but many small problems at the same time. Instead of looping through each pixel of an image to decide its color, a single shader program is run on many pixels in parallel, at the same time. Shader programs, therefore, only need to determine the color of the specific pixel they are being executed for. Yes! A single shader program is executed many times, one for each pixel, without any information being shared between runs or any specific order!

This can be hard to grasp at first, so let's start with an example. Let's create a vertical gradient of 300x300, from green to magenta, like the following:

![](https://github.com/tinne26/kage-desk/blob/main/img/intro_cpu_gradient.png?raw=true)

The CPU program would go like this:
```Golang
func Gradient() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	for y := 0; y < 300; y++ {
		// compute the gradient color for this row
		floatGreen   := (300.0 - float64(y + 1))/300.0
		floatMagenta := (float64(y + 1)/300.0)

		// convert the color from float to rgba, 8-bits per channel
		green   := uint8(255*floatGreen)
		magenta := uint8(255*floatMagenta)
		clr := color.RGBA{magenta, green, magenta, 255}

		// apply the color to the whole row
		for x := 0; x < 300; x++ {
			img.SetRGBA(x, y, clr)
		}
	}

	return img
}
```
*(Full program available at [examples/intro/cpu/gradient](https://github.com/tinne26/kage-desk/blob/main/examples/intro/cpu/gradient))*

The GPU program, instead, would go like this (pseudo-code only for the moment):
```Golang
func Gradient(position vec2) vec4 {
	green   := (299 - position.y)/299
	magenta := position.y/299
	return vec4(magenta, green, magenta, 1.0)
}	
```
*(Full program available at [examples/intro/gpu/gradient](https://github.com/tinne26/kage-desk/blob/main/examples/intro/gpu/gradient))*

Notice that there's no outer loop. Without going into details, we can see that we only take the position of the pixel to compute and figure out its color. An operation for a single pixel.

In this case, having the CPU version of the program makes it easy to come up with the GPU counterpart, but this won't always be the case. Thinking in the sequential version can be helpful sometimes, but you will have to learn many GPU-specific techniques in order to really get good at shaders. It's a different paradigm, so you have to start thinking different.

Try it out! Open the terminal and copy the following:
```
go run github.com/tinne26/kage-desk/examples/intro/gpu/gradient@latest
```
*(You need to have Golang 1.19 or above, and if you are on linux and have never used Ebitengine, you may have to install a few [additional dependencies](https://ebitengine.org/en/documents/install.html?os=linux#Installing_dependencies))*

Congratulations, you have successfully executed your first Kage shader!

Key takeaway:
> When creating shaders, we need to break a task into a pixel-level independent process.


## The first Kage program

We have got the general idea now, so it's time to get our hands dirty.

First, create a folder, run `go mod init first-shader` within it, and create a `main.go` with this content:
```Golang
package main
import _ "embed"
import "github.com/tinne26/kage-desk/tools/display"

//go:embed shader.kage
var shaderSource []byte
func main() {
	display.Shader(shaderSource, "intro/first-shader", 512, 512)
}
```

Ignore this for the moment, it's only helper code to get our shaders running more easily at the start. Create also a `shader.kage` file with the following content:
```Golang
package main

func Fragment(_ vec4, _ vec2, _ vec4) vec4 {
	// ...
}
```
*(Tip: you can [configure your editor](https://github.com/tinne26/kage-desk/blob/main/tutorials/config_editor.md) to highlight `.kage` files like `.go` programs)*

In main programs, the "main" function is the entry point. In Ebitengine shaders, the "Fragment" function is the entry point instead. The reason the entry function is called "Fragment" is because there are multiple types of shader programs: vertex shaders, geometry shaders, compute shaders, tessellation shaders... and fragment shaders. Fragment shaders, also called pixel shaders, are shaders that compute the color of a single pixel or fragment. Ebitengine only has fragment shaders at the moment.

Anyway, back to work. We don't know what we want to do yet, but let's just execute the program anyway! `go run main.go`!

Oh... it didn't work?
```
Failed to load shader:
3:1: function Fragment must have a return statement but not
```

Well, as we were just saying, a fragment shader should return the value of a pixel. It's color. We aren't returning anything yet, so let's fix the `Fragment` function by adding a `return vec4(1, 0, 0, 1)`. If you did it right, re-running `main.go` should now result in a red screen.

You probably figured it out on your own, but as you can see, we don't use `color.RGBA` within shaders, but rather vector types. Check the function signature again:
```Golang
func Fragment(_ vec4, _ vec2, _ vec4) vec4
```

We have a few input vectors of different sizes that we are ignoring for the moment, and one output `vec4`. That output vector is the color of the pixel! With alpha included after the red, green and blue values.

Vectors in Kage are made of `float` values, which have 32-bit precision. That's why the color values don't range from 0 to 255 but rather 0 to 1. If you use values outside that range in the returned color, they will be clamped to 0 - 1.

These vectors are actually quite cool, and you can do many weird operations with them:
```Golang
vec4(1) // creates a vector with 4 components, all set to 1
vec4(0.5, 0.5, 1.0, 1.0).rgb // gets a vec3 with the r, g, and b components from the vec4
vec3(32, 44.0, 0).xy // fields can't only be accessed as rgba, but also xyzw or stpq
vec2(3, 5).yxx // you can even use a different order or repeated fields!
```
The field access thingie is called "swizzling", if you needed the fancy name.


## The `position` input parameter

So, we got started with our Kage shaders, running directly on Ebitengine! Cool!

But not that cool, actually. If you tried to make something interesting with the previous shader, you probably noticed you could only change the color. So here's the big question: how do we actually make *each pixel different* to make these shaders flashier?

There are many ways, but the first step are the input parameters passed to the `Fragment` function. Let's uncover the first one:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4
```

As you may expect, the `position` argument tells us what pixel we are working with. If the position is `(0, 0)`, we are working with the pixel at the top-left corner of the screen. If it's `(320, 240)`, it's probably somewhere around the middle of the screen.

But you may be wondering: if Ebitengine is a 2D game engine... why does the position have 4 components? The answer is that in general 3D, the `Z` component is obviously also important, and the `W` too for the camera projection... but in Ebitengine they are both always 0, so you should just ignore that. You may even do this for comfort:
```Golang
func Fragment(pos4 vec4, _ vec2, _ vec4) vec4 {
	position := pos4.xy // now position is a vec2 with only x and y
	// ...
}
```

With this new tool we can finally start making something a bit more interesting... like the gradient from the first section. You can try to reproduce it on your own if you want. Otherwise, let's jump to the next challenge: keeping the screen size at 512, make the left half be white and the right half be white.

*(Yes, please, try it by yourself if you are actually trying to learn)*

Did you do it? Here's a possible answer:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	if position.x < 256 {
		return vec4(1) // white
	} else {
		return vec4(vec3(0), 1) // black
		// ^ normal people just return vec4(0, 0, 0, 1),
		//   but now you know this also works
	}
}
```

Before the next section, this is probably a good place to mention that the position coordinates are global, not local. If we had a 400x400 screen on our game's `Draw()` function and we drew a shader on the bottom right quadrant, we would see x and y coordinates go from 200 to 399, not 0 to 199.

## Built-in functions

The half-white half-black shader was too easy, so let's complicate it a bit more now. For example, if I asked you to make a shader that creates a pattern by making a pixel white if it's on an even position, and black otherwise... would you be able to make it?

Well, with what you have been explained up to this point, you won't. Because you don't know how to check if a number is even or odd in a Kage shader. It's time to present the *built-in* functions. Although there's a [full list](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(mathematics)) at the official ebitengine.org page for shaders, I don't recommend going there yet. Just know it exists. Instead, I'll share the most common ones here:
- Single argument functions: `abs`, `sign` (returns -1, 0 or 1), `sin`, `cos`, `sqrt`, `floor`, `ceil`, `fract` (returns the fractional part of a number), `length` (length of a vector), `cap` and `len` (both with the same meaning as in Golang).
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

If your screen has any form of DPI scaling, you may see the result somewhat aliased. What would you do to make the pattern bigger? Like, alternating 2x2 pixels white and 2x2 pixels black? This one is a bit trickier, but that's the kind of math problems you will often find with shaders. Think about it for a while!

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

You can also divide by 4, for example, to see the result more clearly. If you got it, well done! If you didn't, don't worry. There's a common idea here that's worth explaining as it comes up all the time when writing shaders: conceptually, we wanted to do the same as in the previous example... but at a different scale. We wanted to project the original "canvas" to one that was half the size, and then apply the same `mod` function as we had been doing. This idea of scaling / projecting / deforming a space or surface is quite common in shaders. In general, when you have to work at the context of a single pixel, mathematical transformations are a very powerful tool.

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


## The hidden infrastructure

While writing the first shaders we have been using `tools/display` to keep our `main.go` file really simple. This is great and all, but if we want to keep expanding our powers, we need to take a look at what's behind that. We have started to see how to make shaders, but we still haven't seen how to invoke them from Ebitengine.

There are two options: [Image.DrawRectShader(...)](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawRectShader) and [Image.DrawTrianglesShader(...)](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawTrianglesShader).

Let's start by reworking our `main.go` with `DrawRectShader`, the simplest of the two methods:
```Golang
package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	// ...
}
```


## More input: uniforms
...


## Using a texture
...


## Screen effects vs sprite effects
...

