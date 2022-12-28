## The `position` input parameter

So, we got started with our Kage shaders, running directly on Ebitengine! Cool!

But... not that cool, actually. If you tried to make something interesting with the previous shader, you probably noticed you could only change the color. So here's the big question: how do we actually make *each pixel different*? We need different colors if we want to make something awesome!

There are many ways, but the first step are the input parameters passed to the `Fragment` function. Let's reveal the first one:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4
```

As you may expect, the `position` argument tells us what pixel we are working with. If the position is `(0, 0)`, we are working with the pixel at the top-left corner of the screen. If it's `(320, 240)`, it's probably somewhere around the middle of the screen. Know also that position coordinates are global, not local. If we received a 400x400 screen image on our game's `Draw()` function and we drew a shader on the bottom right quadrant, on the shader we would see x and y coordinates go from 200 to 399, not 0 to 199.

Back to the topic: so we have this `vec4` position... but, if Ebitengine is a 2D game engine, why does the position have 4 components? The answer is that there's no good reason. In general 3D `Z` and `W` components are relevant (`W` is used for perspective projection), but in Ebitengine they are both always 0, so you can ignore them. You may even do this for comfort:
```Golang
func Fragment(pos4 vec4, _ vec2, _ vec4) vec4 {
	position := pos4.xy // now position is a vec2 with only x and y
	// ...
}
```

With this new tool we can finally start making something a bit more interesting... like the gradient from the first section. You can experiment with that if you want, but I'll also give you a new challenge: keeping the screen size at 512, make the left half be white and the right half be black.

Don't scroll down! Try it on your own!

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

We are almost ready to start making fun shaders now!


### Table of Contents
Next up: [#4](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md).

0. [Main](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/00_main.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md)
3. [**The `position` input parameter**](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms]()
7. [Using a texture]()
8. [Screen vs sprite effects]()
9. [Performance considerations]()
10. [Graduation challenges]()
