# Loops are tricky

One thing we haven't discussed yet are loops. In Kage, loops can be used... but only within ranges determined by constants:
```Golang
for i := -2; i < 2; i++ {
	// (do something)
}
```

This is a very harsh limitation. For example, it means that we can't use a uniform to dynamically determine the size of the area we want to iterate over.

In this chapter we will tackle some practical problems and discuss how to relax or bypass some of these limitations.

The first example of a shader that can be done using a loop will be a *pixelator effect*. Your first exercise will be to create something like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/pixelated_creature.png?raw=true)

The idea for this shader is not too complex: set a pixelation cell size (e.g. 8x8 pixels), find the cell corresponding to the current position, average the colors of all pixels in that cell and return that value. This last step is the part where you will use a for loop.

Make a shader that uses this loop as a reference:
```Golang
for y := 0.0; y < CellHeight; y += 1.0 {
	for x := 0.0; x < CellWidth; x += 1.0 {
		// ...
	}
}
```

And where `CellHeight` and `CellWidth` are constants declared like this:
```Golang
const CellWidth  float = 12.0 // must be at least 1
const CellHeight float = 12.0 // must be at least 1
```

<details>
<summary>Click to show the solution</summary>

```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	// find the position of the cell we are working on
	baseX := floor(position.x/CellWidth)*CellWidth
	baseY := floor(position.y/CellHeight)*CellHeight

	// iterate the pixelization cell
	colorAcc := vec4(0.0) // color accumulator
	for y := 0.0; y < CellHeight; y += 1.0 {
		for x := 0.0; x < CellWidth; x += 1.0 {
			pixCoords := vec2(baseX + x, baseY + y)
			colorAcc += imageColorAtPixel(pixCoords)
		}
	}

	// divide the color to average it
	return colorAcc/(CellWidth*CellHeight)
}
```
*(Full program available at [examples/intro/pixelize](https://github.com/tinne26/kage-desk/blob/main/examples/intro/pixelize))*
</details>

This technique can also be used to create blurs, motion blurs, implement [image kernels](https://setosa.io/ev/image-kernels/) and much more, but keep in mind that looping can quickly get expensive and slow.

## Overcoming the limitations

This is all very nice, but having to use constants instead of uniforms can get a bit painful. What if we wanted to animate the pixelization effect, for example, making cells bigger and bigger each time?

https://user-images.githubusercontent.com/95440833/211917307-81d31314-8746-4efe-a5eb-a15db877485e.mp4

There are three tricks I can share with you:
- **Upper bounding**: your looping constant will act as an upper bound, and you will use the `break` keyword to break earlier if possible. Since the uniform value will be the same for all shaders, they should all finish at the same time (as opposed to have to wait for the slowest one).
- **Not looping**: sometimes you can just fake it. For pixelation, for example, you could take the central pixel of the cell regardless of the cell size. This sacrifices accuracy, but in some cases it's may be ok! For some animations it may work well enough.
- **Constant sampling**: sometimes you can't just fake it... but you still can kinda fake it. For example, what if you didn't check all the pixels within the cell, but only 6 values at properly distributed locations? This is a probabilistic method that can help you balance cost and accuracy, while still allowing you to scale your cell size as you want.

Try to use the first approach to adapt our previous shader and make it animated as shown in the video above, making cell sizes go from 1 to 32, then back to 1 and so on.

<details>
<summary>Click to show the solution</summary>

```Golang
var CellSize float // uniform: max value is MaxCellSize
const MaxCellSize float = 32.0

func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	// find the position of the cell we are working on
	baseX := floor(position.x/CellSize)*CellSize
	baseY := floor(position.y/CellSize)*CellSize

	// iterate the pixelization cell
	colorAcc := vec4(0.0) // color accumulator
	for y := 0.0; y < MaxCellSize; y += 1.0 {
		if y >= CellSize { break }
		for x := 0.0; x < MaxCellSize; x += 1.0 {
			if x >= CellSize { break }
			pixCoords := vec2(baseX + x, baseY + y)
			colorAcc += imageColorAtPixel(pixCoords)
		}
	}

	// divide the color to average it
	return colorAcc/(CellSize*CellSize)
}
```
*(Full program available at [examples/intro/pixelize-anim](https://github.com/tinne26/kage-desk/blob/main/examples/intro/pixelize-anim))*

With this shader, if you open up your GPU software monitor you will already be able to observe that when the cell sizes increase, the GPU load also increases, creating a sine wave of GPU load over time. This shader can still be optimized by manually inlining the helper function, moving the reused values outside the loop and computing the texture coordinates as fixed deltas before entering the loop. With this we can avoid the divisions in the inner part of the loop and get a performance improvement somewhere between 15-20%. The optimized code can be found at [examples/intro/pixelize-anim-opt](https://github.com/tinne26/kage-desk/blob/main/examples/intro/pixelize-anim-opt), but it requires you to have read the [tutorial explaining texels](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/texels.md) to really understand everything that's going on. This is offered as an optimization exercise, but it's not part of the main tutorial (optimization is not one of the goals of the introduction).
</details>


### Table of Contents
<!-- Next up: [#9](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md). -->

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [`DrawTrianglesShader()`](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_triangles.md)
9. [**Loops are tricky**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
