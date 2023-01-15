# `DrawTrianglesShader()`

We have now learned to use images, but we are still only scratching the surface of what we can do with them. The next exercise I wanted you to try consists in mirroring the spider-cat-dog to obtain something like this:

![](https://github.com/tinne26/kage-desk/blob/main/img/mirrored_creature.webp?raw=true)

...There's only a small problem.

You may remember from the previous chapters that the `DrawRectShaderOptions` image sizes and the shader's target rectangle must be the same size. While it's true that we could draw the first part normally and the second using the shader, I want you to do everything in the same shader instead. After all, it won't be a proper training arc if I don't make you suffer a little bit.

Since `DrawRectShader()` won't work, we need to switch to `DrawTrianglesShader()` instead.

Triangle-drawing calls in Ebitengine are confusing for *a lot of people*, so we have made a separate [tutorial](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/misc/triangles.md) for them. If you want to properly learn how triangles work, jump there and read. Otherwise, we have also prepared a helper function in `kage-desk/display` so you can adapt your `main.go` file and keep going even without understanding the thing that we titled this chapter with:
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawTrianglesShaderOptions{}

	// set images for the shader
	opts.Images[0] = display.ImageSpiderCatDog()

	// set uniforms for the shader
	bounds := display.ImageSpiderCatDog().Bounds()
	rect := image.Rect(0, 0, bounds.Dx(), bounds.Dy()*2)
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["TargetRect"] = display.RectToUniform(rect)
	opts.Uniforms["MirrorAlphaMult"] = float32(0.2)
	opts.Uniforms["VertDisplacement"] = 28
	
	// draw shader
	display.DrawShader(screen, rect, self.shader, opts)
}
```

The situation is very similar to what we've been seeing in the previous chapters, but we are using [`DrawTrianglesShaderOptions`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#DrawTrianglesShaderOptions) for the draw options instead of `DrawRectShaderOptions` (which also have `Images` and `Uniforms` fields, but not `GeoM`), `display.DrawShader()` instead of `ebiten.DrawRectShader()`, and we are using `bounds.Dy()*2`. This new size is only necessary for this specific shader, but you will also need to adapt the layout function and `SetWindowSize()` to use `bounds.Dy()*2`.

For the uniforms, here's what we are doing:
- `TargetRect` is important so you can compute both the main position of the creature and its reflection in relation to the shader's target rectangle. The actual content of `TargetRect` is `vec4(minX, minY, maxX, maxY)`.
- `MirrorAlphaMult` is an opacity multiplier for the reflection. It is fairly easy to apply, but you don't need to worry about it until your are cleaning up the reflection effect.
- `VertDisplacement` is an advanced and optional uniform that you should ignore until you get everything else working. The idea is that the mirrored image can have too much padding around the eges and look too disconnected from the reflection. We can use this configurable factor to bring the two closer to the center of the shader's target rectangle and make it look better. Notice that this uniform is an `int`! It could be a `float` too, but I wanted to throw an `int` in a shader at some point so you didn't forget about them.

Remember also the `imageColorAtPixel(vec2)` helper function from the previous chapter, you will definitely need it! Or its sibling `imageColorAtUnit(vec2)`.

Ok, the time has come: with the current setup, try to write the mirror shader by yourself. This is probably the hardest shader you will be asked to write in the tutorial, so take your time... and don't get frustrated if you fail, but at least try to come out of the attempt with *more and more concrete questions* than you had going in.

<details>
<summary>Click to show the solution</summary>

```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	// prepare helper information
	relativePos := vec2(position.x - TargetRect[0], position.y - TargetRect[1])
	rectHeight := TargetRect[3] - TargetRect[1]
	yCenter := rectHeight/2

	// apply displacement
	relativePos.y += float(VertDisplacement)*sign(relativePos.y - yCenter)

	// top part (unmodified creature)
	mainColor := imageColorAtPixel(relativePos)

	// bottom part (inverted and alpha-adjusted creature)
	mirrorPosition := vec2(relativePos.x, rectHeight - relativePos.y)
	mirrorColor := imageColorAtPixel(mirrorPosition)*MirrorAlphaMult

	// compose the result
	return mainColor + mirrorColor
}
```
*(Full program available at [examples/intro/mirror](https://github.com/tinne26/kage-desk/blob/main/examples/intro/mirror))*
</details>


### Table of Contents
Next up: [#9](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [**`DrawTrianglesShader()`**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_triangles.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
