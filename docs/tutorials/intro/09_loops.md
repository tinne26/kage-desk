# Loops are tricky

One thing we haven't discussed yet are loops. In general, shaders do support loops, but in limited ways. In the case of Kage, we can only loop within ranges determined by constants:
```Golang
for i := -2; i < 2; i++ {
	// (do something)
}
```

This is a very harsh limitation: if the bounds must be given by constants, we can't use uniforms to dynamically control the size of the area we want to iterate over.

In this chapter we will be tackling some practical problems that involve loops and discuss how to relax or bypass some of these limitations.

The first challenge involving a loop will be a *pixelator effect* shader. Your first exercise will be to create something like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/pixelated_creature.png?raw=true)

The idea for this shader is not too complex:
1. Set a pixelation cell size (e.g. 8x8 pixels).
2. Find the cell corresponding to the current position.
3. Average the colors of all pixels in that cell and return that value.

This last step is the part where you will use a for loop. For example, you can use this as a reference:
```Golang
for y := 0.0; y < CellHeight; y += 1.0 {
	for x := 0.0; x < CellWidth; x += 1.0 {
		// ...
	}
}
```

Where `CellHeight` and `CellWidth` are constants declared in the following way:
```Golang
const CellWidth  float = 12.0 // must be at least 1
const CellHeight float = 12.0 // must be at least 1
```

Try to write the pixelator effect by yourself!

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

This is all very nice, but having to use constants instead of uniforms is not great. What if we want to animate the pixelization effect like in the following video, changing the cell size?

https://user-images.githubusercontent.com/95440833/212187726-a42a80d6-42c6-4e74-95e4-f6bcc404bae2.mp4

There are three key tricks you should know about to get around these limitations:
- **Upper bounding**: if you make your looping constant act as an upper bound for the number of loop iterations, you can use the `break` keyword to break earlier based on the value of an actual uniform. It doesn't look as sleek as a regular loop, but it works well enough.
- **Not looping**: sometimes you can just fake it. For pixelation, for example, you could take the central pixel of the cell regardless of the cell size. This will be sacrificing accuracy, but for some animations it can work well enough!
- **Constant sampling**: sometimes you can't just completely fake it... but you still can *half fake it*! What if you didn't check all the pixels within the cell, but only 6 values at properly distributed locations? This is a probabilistic method that can help you balance cost and accuracy, while still allowing you to scale your cell size as you want.

The last challenge for this chapter will be to use the first approach to adapt our previous shader and make it animated as shown in the video above, making cell sizes go from 1 to 32, then back to 1 and repeat. Show me what you have learned!

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

With this shader, if you open up your GPU software monitor you will already be able to observe that when the cell sizes increase, the GPU load also increases, creating a sine wave of GPU load over time.

This shader can still be optimized by manually inlining the helper function, moving the reused values outside the loop and computing the texture coordinates as fixed deltas before entering the loop. With this we can avoid the divisions on the inner part of the loop and get a performance improvement somewhere between 15-20%. The optimized code can be found at [examples/intro/pixelize-anim-opt](https://github.com/tinne26/kage-desk/blob/main/examples/intro/pixelize-anim-opt), but it requires you to have read the [tutorial explaining texels](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/texels.md) to really understand everything that's going on. This is offered as an optimization exercise, but it's not part of the main tutorial (optimization is not one of the goals of the introduction).
</details>


### Table of Contents
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
