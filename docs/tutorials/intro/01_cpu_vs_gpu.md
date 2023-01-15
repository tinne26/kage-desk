# CPU vs GPU: different paradigms

Jumping right into Kage would be nice, but before that it may be wiser to try to understand how CPU and GPU programs differ.

When we talk about shaders, we have to talk about **parallelism**. While CPUs can tackle big problems sequentially, GPUs work solving many small problems at the same time. Instead of looping through each pixel of an image to compute their color, shader programs are executed on several pixels at once, with many instances of the same program running in parallel.

This can be hard to grasp at first, so let's begin with a simple example.

Let's try to create a vertical color gradient that goes from green to blue. We can make it 300x300 pixels, like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/intro_cpu_gradient.png?raw=true)

We will start with the CPU version. It's quite simple, as we only need to fill 300 rows with a progressively changing color for each row. Since the gradient is so simple, we can start with 100% green and 0% blue, and end at 0% green and 100% blue. It could look like this:
```Golang
func Gradient() *image.RGBA {
	// create the target image
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	// for each row...
	for y := 0; y < 300; y++ {
		// compute the gradient color for this row (amount of blue and green)
		greenLevel := (300.0 - float64(y) + 0.5)/300.0
		blueLevel  := ((float64(y) + 0.5)/300.0)

		// convert the float values to an rgba color, 8-bits per channel
		greenValue := uint8(255*greenLevel)
		blueValue  := uint8(255*blueLevel)
		clr := color.RGBA{0, greenValue, blueValue, 255}

		// apply the color to the whole row
		for x := 0; x < 300; x++ {
			img.SetRGBA(x, y, clr)
		}
	}	

	return img
}
```
*(Full program available at [examples/intro/gradient-cpu](https://github.com/tinne26/kage-desk/blob/main/examples/intro/gradient-cpu))*

To contextualize, one could use a gradient like this as the background for a videogame. We could compute it at the start of the game and store it in an image that would then be sent to the GPU. Otherwise, if we wanted to recompute this gradient on each frame... it could be better to use a shader instead. With the shader, we wouldn't need to recompute and transfer the image to the GPU on each frame: we would be passing *the shader program* to the GPU instead. In this simple case where the gradient doesn't change through time, the first approach may be simpler, but in more complex cases shaders would start to show more advantages and benefits.

With that out of the way, let's take a look at the GPU version (in pseudo-code):
```Golang
func Gradient(position vec2) vec4 {
	green := (300 - position.y)/300
	blue  := position.y/300
	return vec4(0, green, blue, 1.0)
}	
```
*(Full program available at [examples/intro/gradient](https://github.com/tinne26/kage-desk/blob/main/examples/intro/gradient))*

Notice that in this version there's no outer loop. We won't go into much detail yet, but you can see that all we are doing here is take the position of the pixel we are working on and compute its color. One pixel, one calculation.

In this case, the core code of the CPU and the GPU versions is very similar, but this isn't always the case. Thinking in the sequential version can be helpful sometimes, but you will have to learn many GPU-specific techniques in order to really get good at shaders. It's a different paradigm, so you need to start thinking differently.

Try to execute this first shader now. Open the terminal and run the following:
```
go run github.com/tinne26/kage-desk/examples/intro/gradient@latest
```
*(You need to have Golang 1.19 or above, and if you are on linux and have never used Ebitengine, you may have to install a few [additional dependencies](https://ebitengine.org/en/documents/install.html?os=linux#Installing_dependencies))*

That's it, you have successfully executed your first Kage shader!

Key takeaway:
> When creating shaders, we need to break a task into a pixel-level independent process.


### Table of Contents
Next up: [#2](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [**CPU vs GPU: different paradigms**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [`DrawTrianglesShader()`](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_triangles.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
