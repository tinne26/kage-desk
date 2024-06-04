# Beyond one to one mapping

We have now learned to use images, but we are still only scratching the surface of what we can do with them. The next exercise I wanted you to try consists in mirroring the spider-cat-dog to obtain something like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/mirrored_creature.webp?raw=true)

...There's only a small problem: the one-to-one mapping techniques we have been using up to this point don't seem to be sufficient for the task.

While it's true that we could draw the first part normally and the second using the shader, I want you to do everything in the same shader instead. After all, it won't be a proper training arc if I don't make you suffer a little bit.

To get started, I will give you some reference code to help you out. You can copy the code from the previous chapter into a new project, and modify it with the following bits:
```Golang
func main() {
	// ...
	
	// set a specific window size: like spider-cat-dog, but twice as tall
	bounds := display.ImageSpiderCatDog().Bounds()
	ebiten.SetWindowSize(bounds.Dx(), bounds.Dy()*2)

	// ...
}

// request a canvas with a specific size matching the logical window size
func (self *Game) Layout(_, _ int) (int, int) {
	bounds := display.ImageSpiderCatDog().Bounds()
	return bounds.Dx(), bounds.Dy()*2
}

func (self *Game) Draw(screen *ebiten.Image) {
	// map the vertices to the target image
	bounds := screen.Bounds()
	self.vertices[0].DstX = float32(bounds.Min.X) // top-left
	self.vertices[0].DstY = float32(bounds.Min.Y) // top-left
	self.vertices[1].DstX = float32(bounds.Max.X) // top-right
	self.vertices[1].DstY = float32(bounds.Min.Y) // top-right
	self.vertices[2].DstX = float32(bounds.Min.X) // bottom-left
	self.vertices[2].DstY = float32(bounds.Max.Y) // bottom-left
	self.vertices[3].DstX = float32(bounds.Max.X) // bottom-right
	self.vertices[3].DstY = float32(bounds.Max.Y) // bottom-right

	// set the source image sampling coordinates
	srcBounds := display.ImageSpiderCatDog().Bounds()
	self.vertices[0].SrcX = float32(srcBounds.Min.X) // top-left
	self.vertices[0].SrcY = float32(srcBounds.Min.Y) // top-left
	self.vertices[1].SrcX = float32(srcBounds.Max.X) // top-right
	self.vertices[1].SrcY = float32(srcBounds.Min.Y) // top-right
	self.vertices[2].SrcX = float32(srcBounds.Min.X) // bottom-left
	self.vertices[2].SrcY = float32(srcBounds.Max.Y) // bottom-left
	self.vertices[3].SrcX = float32(srcBounds.Max.X) // bottom-right
	self.vertices[3].SrcY = float32(srcBounds.Max.Y) // bottom-right

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Images[0] = display.ImageSpiderCatDog()
	shaderOpts.Uniforms = make(map[string]interface{})
	shaderOpts.Uniforms["MirrorAlphaMult"] = float32(0.2)
	shaderOpts.Uniforms["VertDisplacement"] = 28

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```

The quick breakdown is that we are using specific window, layout and canvas sizes to make the exercise more manageable. We want to draw the spider-cat-dog and its reflection, so we need twice the vertical space. For the uniforms, here's what we are doing:
- `MirrorAlphaMult` is an opacity multiplier for the reflection. It is fairly easy to apply, but you don't need to worry about it until your are cleaning up the reflection effect.
- `VertDisplacement` is an advanced and optional uniform that you should ignore until you get everything else working. The idea is that the mirrored image can have too much padding around the edges and look too disconnected from the reflection. We can use this configurable factor to bring the two closer to the center of the shader's target rectangle and make it look better. Notice that this uniform is an `int`! It could be a `float` too, but I wanted to throw an `int` in a shader at some point just so you didn't forget about them.

Now the first step will be running this as it is and see what you get:
```Golang
package main

var MirrorAlphaMult float // uniform: reflection opacity multiplier
var VertDisplacement int  // uniform: displacement towards the center

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	return imageSrc0At(sourceCoords)
}
```

You should be getting a stretched spider-dog-cat. That's good. When you are writing shaders, since you can't add `Printf`s and debug code so easily, it's best to start very slowly and try to make progress step by step. You are almost ready to take it from here, but let me give you a few more tools:
- You can use [`imageDstSize()`](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(images)) to get the target image size, in pixels, within a shader.
- You can use [<code>imageSrc<b><i>N</i></b>Size()</code>](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(images)) to get the size of the source image N, in pixels, within a shader.
- You don't really need it in this case since origins are (0, 0) everywhere, but you can use `imageDstOrigin()` and <code>imageSrc<b><i>N</b></i>Origin()</code> to get the origin coordinates of both target and source images.

Ok, the time has come: with the current setup, try to write the mirror shader by yourself. This is probably the hardest shader you will be asked to write in the tutorial, so take your time... and don't get frustrated if you fail, but at least try to come out of the attempt with *more and more concrete questions* than you had going in. This whole chapter is about tackling a non-trivial problem for the first time, so try to put some real effort into it.

<details>
<summary>Click to show tips if stuck for more than 15 minutes</summary>

Some steps that would be relevant:
- Make sure you know how to draw the top half of the screen in one color, and the lower half in another.
- Make sure you can draw the spider-cat-dog on the top half, isolated, without stretching.
- Make sure you can draw the spider-cat-dog on the lower half, isolated, without stretching.
</details>

<details>
<summary>Click to show the solution</summary>

I have two solutions for you. The first is a "simple" solution without using `VertDisplacement`, and the second is a more general solution.

Let's see the "simple one" first:
```Golang
//kage:unit pixels
package main

var MirrorAlphaMult float // uniform: reflection opacity multiplier
var VertDisplacement int  // uniform: displacement towards the center

func Fragment(targetCoords vec4, sourceCoords vec2, _ vec4) vec4 {
	if targetCoords.y < imageDstSize().y/2 {
		return imageSrc0At(vec2(sourceCoords.x, sourceCoords.y*2))
	} else {
		adjustedY := (sourceCoords.y - imageSrc0Size().y/2)*2
		invertedY := imageSrc0Size().y - adjustedY
		samplingCoords := vec2(sourceCoords.x, invertedY)
		return imageSrc0At(samplingCoords)*MirrorAlphaMult
	}
}
```
The code shouldn't be too hard to understand. If we are on the top half of the screen, we sample the source image... but since it would be stretched, we multiply the `y` by 2 to make it fit properly into the top area. The x value is always already what we want, so it never has to be modified. For the lower half, the normalization of the `y` value is slightly trickier, but it's not much different. Since we are already past the midpoint position through the source image, we need to offset that and then apply the same `y*2` idea of the first branch. To invert the image, we simply use `height - y`. Finally, we sample the value at the relevant position and multiply it by `MirrorAlphaMult`, which yes, is a simple product on the sampled pixel color.

This code does rely on the fact that `imageSrc0Size()` is exactly half `imageDstSize()`, and we could actually use only one of them for everything. This is not very clean, not very nice and whatever, but if you understand the constraints of your context... you don't have to be a perfectionist, just get the job done. You could parametrize all this with further uniforms, or automatically center everything on an arbitrarily-sized target, or whatever. That's a pain and not fun, so do it as homework if you really have nothing better to do.

Now let's jump onto the more general solution using `VertDisplacement`:
```Golang
//kage:unit pixels
package main

var MirrorAlphaMult float // uniform: reflection opacity multiplier
var VertDisplacement int  // uniform: displacement towards the center

func Fragment(targetCoords vec4, sourceCoords vec2, _ vec4) vec4 {
	// compute top contribution
	uprightColor := imageSrc0At(vec2(sourceCoords.x, sourceCoords.y*2 - float(VertDisplacement)))
	
	// compute bottom contribution
	adjustedY := (sourceCoords.y - imageSrc0Size().y/2)*2
	invertedY := imageSrc0Size().y - adjustedY
	samplingCoords := vec2(sourceCoords.x, invertedY - float(VertDisplacement))
	mirrorColor := imageSrc0At(samplingCoords)*MirrorAlphaMult

	// return the sum of contributions
	return uprightColor + mirrorColor
}
```
*(Full program available at [examples/intro/mirror](https://github.com/tinne26/kage-desk/blob/main/examples/intro/mirror))*

This second shader is not that different, but it has a few subtle ideas worth explaining:
- We are computing different parts of the image all in a single pass and adding them at the end. This is a common pattern used in many shaders, but it often requires gating or filtering the partial results before summing them.
- The reason why we didn't "filter" the results here is that `imageSrc0At()` will return `vec4(0, 0, 0, 0)` if we are requesting positions out of bounds, and *it just happens that for this particular example*, the different calculations for different parts of the image do not collide (given reasonable `VertDisplacement` values, at least).
- For proper filtering, you could use [selectors](https://github.com/tinne26/kage-desk/blob/main/docs/snippets/selectors.md) to keep or discard specific results. Again, this is not necessary in this specific situation, but you could totally add something like `uprightColor *= whenLessThan(sourceCoords.y, imageSrc0Size().y/2)` and `mirrorColor *= whenGreaterThan(sourceCoords.y, imageSrc0Size().y/2)` to be a bit safer.
</details>

Good work! We are getting close to the end now!


### Table of Contents
Next up: [#9](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [**Beyond one-to-one mapping**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
