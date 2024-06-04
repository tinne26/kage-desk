# Loops are tricky

One thing we haven't discussed yet are loops. While in general shaders do support loops, they only do it in limited ways. In the case of Kage, we can only use classical `for` loops within ranges determined by *constants*:
```Golang
for i := -2; i < 2; i++ {
	// (do something)
}
```

This is a very harsh limitation: if the bounds must be given by constants, this means we can't use uniforms to dynamically control the size of the area we want to iterate over.

In this chapter we will be tackling some practical problems that involve loops and discuss how to relax or bypass some of these limitations.

The first challenge involving a loop will be a *pixelator effect* shader. Your first task will be to create something like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/pixelated_creature.png?raw=true)

The idea for this shader is not too complex:
1. Decide a pixelation cell size (e.g. 8x8 pixels).
2. Find the cell corresponding to the current position.
3. Average the colors of all pixels in that cell and return that value.

This last step is the part where you will use a `for` loop. For example, you can use this as a reference:
```Golang
for y := 0.0; y < CellHeight; y += 1.0 {
	for x := 0.0; x < CellWidth; x += 1.0 {
		// ...
	}
}
```

Where `CellHeight` and `CellWidth` are constants declared like this:
```Golang
const CellWidth  float = 12.0 // must be at least 1
const CellHeight float = 12.0 // must be at least 1
```

Try to write the pixelator effect by yourself now!

> [!TIP]
> *Fix the window size and the layout return to `display.ImageSpiderCatDog().Bounds()`.*

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

const CellWidth  float = 12.0 // must be at least 1
const CellHeight float = 12.0 // must be at least 1

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	// find the origin of the cell we are working with
	cellOX := floor(sourceCoords.x/CellWidth )*CellWidth
	cellOY := floor(sourceCoords.y/CellHeight)*CellHeight

	// iterate the pixelization cell
	colorAcc := vec4(0.0) // color accumulator
	for y := 0.0; y < CellHeight; y += 1.0 {
		for x := 0.0; x < CellWidth; x += 1.0 {
			pixCoords := vec2(cellOX + x, cellOY + y)
			colorAcc += imageSrc0At(pixCoords)
		}
	}

	// divide the color to average it
	return colorAcc/(CellWidth*CellHeight)
}
```
*(Full program available at [examples/intro/pixelize](https://github.com/tinne26/kage-desk/blob/main/examples/intro/pixelize))*
</details>

This technique can also be used to create outlines or implement blurs and other [image kernels](https://setosa.io/ev/image-kernels/), but keep in mind that looping can quickly get expensive and slow.

## Overcoming the limitations

This is all very nice, but having to use constants instead of uniforms is not cool. What if we wanted to animate the pixelization effect as in the following video, where the cell size changes through time?

https://user-images.githubusercontent.com/95440833/212187726-a42a80d6-42c6-4e74-95e4-f6bcc404bae2.mp4

There are three key tricks you should know about to get around these limitations:
- **Upper bounding**: if you make your looping constant act as an upper bound for the number of loop iterations, you can use the `break` keyword to break earlier based on the value of an actual uniform. It doesn't look as sleek as a regular loop, but... it works ¯\_(ツ)_/¯.
- **Not looping**: sometimes you can just fake it. For pixelation, for example, you could take the central pixel of the cell regardless of the cell size. This will be sacrificing accuracy, but for some animations it can work well enough!
- **Constant sampling**: sometimes you can't just *completely fake it*... but you can *half fake it*! What if you didn't check all the pixels within the cell, but only 6 values at properly distributed locations? This is a probabilistic method that can help you balance cost and accuracy while still allowing you to scale your cell size as you want.

The last challenge for this chapter will be to use the first approach to adapt our previous shader and make it animated as shown in the video above. Try to make cell sizes go from 1 to 32, then back to 1 and repeat. Show me what you have learned!

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

var CellSize float // uniform: max value is MaxCellSize
const MaxCellSize float = 32.0

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	// find the origin of the cell we are working with
	cellOrigin := floor(sourceCoords/CellSize)*CellSize

	// iterate the pixelization cell
	colorAcc := vec4(0.0) // color accumulator
	for y := 0.0; y < MaxCellSize; y += 1.0 {
		if y >= CellSize { break }
		for x := 0.0; x < MaxCellSize; x += 1.0 {
			if x >= CellSize { break }
			pixCoords := cellOrigin + vec2(x, y)
			colorAcc += imageSrc0At(pixCoords)
		}
	}

	// divide the color to average it
	return colorAcc/(CellSize*CellSize)
}
```
*(Full program available at [examples/intro/pixelize-anim](https://github.com/tinne26/kage-desk/blob/main/examples/intro/pixelize-anim))*

With this shader, if you open up your GPU software monitor you will already be able to observe that when cell size increases, the GPU load also increases, creating a sine wave of GPU load over time.

Notice also that for this to work correctly, we have to send a whole `CellSize` value, even when cell size is of type `float`. It's easy to make that mistake on the main file and get confused about why the shader does funny things.
</details>


### Table of Contents
0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [**Loops are tricky**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
