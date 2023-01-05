## CPU vs GPU: different paradigms

Jumping right into Kage would be nice, but before that it may be wiser to try to understand how CPU and GPU programs differ.

When we talk about shaders, we have to talk about **parallelism**. If CPUs can solve big problems sequentially, GPUs prefer solving many small problems at the same time. Instead of looping through each pixel of an image to decide their color, shader programs are executed on many pixels at once, executing many instances of the same program in parallel.

This can be hard to grasp at first, so let's start with a simple example. Let's try to create a vertical color gradient that goes from green to blue. We can make it 300x300 pixels, like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/intro_cpu_gradient.png?raw=true)

We will start with the CPU version. It's quite simple, as we only need to fill 300 rows with a progressively changing color for each row. Since the gradient is so simple, we can start with 100% green and 0% blue, and end at 0% green and 100% blue. It could look like this:
```Golang
func Gradient() *image.RGBA {
	// create the target image
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	// for each row...
	for y := 0; y < 300; y++ {
		// compute the gradient color for this row (amount of blue and green)
		greenLevel := (299.0 - float64(y))/299.0
		blueLevel  := (float64(y)/299.0)

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
*(Full program available at [examples/intro/cpu/gradient](https://github.com/tinne26/kage-desk/blob/main/examples/intro/gradient-cpu))*

That should be easy enough to understand. Let's take a look at the GPU version now, in pseudo-code:
```Golang
func Gradient(position vec2) vec4 {
	green := (299 - position.y)/299
	blue  := position.y/299
	return vec4(0, green, blue, 1.0)
}	
```
*(Full program available at [examples/intro/gpu/gradient](https://github.com/tinne26/kage-desk/blob/main/examples/intro/gradient))*

Notice that in this version there's no outer loop. We won't go into much detail, but you can see that we only take the position of the pixel we are working on and compute its color. One pixel, one calculation.

In this case, the core code of the CPU and the GPU versions is almost the same, but this isn't always the case. Thinking in the sequential version can be helpful sometimes, but you will have to learn many GPU-specific techniques in order to really get good at shaders. It's a different paradigm, so you need to start thinking different.

Try to execute this first shader now. Open the terminal and run the following:
```
go run github.com/tinne26/kage-desk/examples/intro/gradient@latest
```
*(You need to have Golang 1.19 or above, and if you are on linux and have never used Ebitengine, you may have to install a few [additional dependencies](https://ebitengine.org/en/documents/install.html?os=linux#Installing_dependencies))*

That's it, you have successfully executed your first Kage shader!

Key takeaway:
> When creating shaders, we need to break a task into a pixel-level independent process.


### Table of Contents
Next up: [#2](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md).

0. [Main](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/00_main.md)
1. [**CPU vs GPU: different paradigms**](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/06_images.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/07_images.md)
8. [Screen vs sprite effects]()
9. [Performance considerations]()
10. [Graduation challenges]()
