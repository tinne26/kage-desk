# The `targetCoords` input parameter

So, we got started with our Kage shaders, running directly on Ebitengine! Cool!

But... not that cool, actually. If you tried to make something interesting with the previous shader, you probably noticed you could only change the color. So here's the big question: how do we actually make *each pixel different*? We need different colors if we want to make something awesome!

There are many ways, but the first step are the input parameters passed to the `Fragment` function. Let's reveal the first one:
```Golang
func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4
```

As you may expect, the `targetCoords` argument tells us what pixel we are modifying. If the coordinates are `(0.5, 0.5)`, we are working with the pixel at the top-left corner of the target image. If the target is the main canvas and the `targetCoords` are `(320.5, 240.5)`, then we are probably somewhere closer to the middle of the screen. Notice that we receive the positions at the *center* of the pixel, not the corner! So, we have `(0.5, 0.5)`, not `(0, 0)`! This is kind of a subtle detail that you will almost never need to think about consciously, but if you are a dev with subpixel perfect centering OCD... it's *critical* to know.

Another important fact to be aware of is that position coordinates are *global*, not local. If we received a 400x400 screen image on our game's `Draw()` function and we drew a shader on the bottom right quadrant, on the shader we would see x and y coordinates go from 200.5 to 399.5, not 0.5 to 199.5.

Back to the topic: so we have this `vec4` position... but, if Ebitengine is a 2D game engine, why does the position have 4 components? Mostly convention. In 3D, `Z` and `W` components are relevant (`W` is used for perspective projection), but in Ebitengine they are both always 0, so you can ignore them. You may even do this for comfort:
```Golang
func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	position := targetCoords.xy // now we have a vec2 with only x and y
	// ...
}
```

With this new tool we can finally start making something a bit more interesting... like the gradient from the first section. You can experiment with that if you want, but I'd rather give you a new challenge: keeping the screen size at 512, make the left half be white and the right half be black. Try it on your own!

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	if targetCoords.x < 256 {
		return vec4(1) // white
	} else {
		return vec4(vec3(0), 1) // black
		// ^ most people just return vec4(0, 0, 0, 1),
		//   but now you know this also works
	}
}
```
*(Full program available at [examples/intro/half-half](https://github.com/tinne26/kage-desk/blob/main/examples/intro/half-half))*
</details>

We are almost ready to start making fun shaders now!


### Table of Contents
Next up: [#4](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [**The `targetCoords` input parameter**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
