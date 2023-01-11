# Using images

There's only one critical piece left to reveal before you can unlock the full power of shaders: images.

There are two main reasons why we might want to pass images to our shaders:
- To **apply an effect** to our game sprites or screen. This is the most common way to use images in 2D games. Pixelize an image, apply a blur or movement blur effect to it, deform it, control its color, apply chromatic aberrations, make simple lighting effects, screen transitions, etc.
- To **use the image as a texture** for the shader. This is very common in 3D, where textures are used for "painting" raw triangles and geometry and for lighting purposes. In 2D this is more unusual, but there are still some use-cases like creating animation effects on a sprite (e.g. being electrocuted), fancy glitches, reflections on the water, making "see behind the wall" effects and a few others that combine multiple images or textures to achieve a specific effect. Advanced lightning techniques with surfaces and normals can also be used in 2D, mainly in top-down view games, but this is rather uncommon, so we won't discuss it in this tutorial.

*(The words _image_ and _texture_ are often used interchangeably in the context of shaders. There are some nuances, but we won't get into it today)*

To get started, we will show how to pass an image to a shader in `main.go`: the draw options struct includes an `Images` array where we can set up to 4 images to be passed to the shader:
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// create draw options
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = display.ImageSpiderCatDog
	
	// draw shader
	bounds := display.ImageSpiderCatDog.Bounds()
	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), self.shader, opts)
}
```

Notice that we are importing `github.com/tinne26/kage-desk/display` again as it includes some images that we can use for our tests quite easily. You can try loading your own images if you want... but there's one **critical limitation** that you must know about:
- Both the shader's target rectangle and the source images must all have the same size.

We will discuss how to get around this in the next chapter, so keep calm and continue learning shaders.

For the shader, the first thing we will do is simply show the image. No fancy effects yet. Just make the shader compute the color of each pixel as the color of the corresponding pixel in the passed image:
```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	return imageColorAtPixel(position.xy)
}
```

This code is interesting because *it doesn't work*. But let me explain...

Kage does have some built-in functions for working with images:
- `imageSrc0At(texelCoords vec2) vec4`: returns the color of an input image at a specific point. As you may imagine, the variants `imageSrc1At()`, `imageSrc2At()` and `imageSrc3At()` also exist.
- `imageSrc0UnsafeAt(texelCoords vec2) vec4`: a variant of the previous that doesn't check bounds. The safe version will return `vec4(0)` if trying to read a position outside bounds. This one is faster, but it may return anything else that resides in the [texture atlas](https://en.wikipedia.org/wiki/Texture_atlas) that contains our image (which Ebitengine manages automatically by default) if we mess up the coordinates.
- `imageSrcTextureSize()`, `imageSrcRegionOnTexture()` and others that we won't discuss.

Despite mentioning these functions, the current advice is that you should not use them directly, because texels in Kage require understanding implementation details that we don't believe you should be concerned with[^1]. Instead, you should generally be using these two helper functions instead:

```Golang
// Helper function to access an image's color at the given pixel coordinates.
func imageColorAtPixel(pixelCoords vec2) vec4 {
	sizeInPixels := imageSrcTextureSize()
	offsetInTexels, _ := imageSrcRegionOnTexture()
	adjustedTexelCoords := pixelCoords/sizeInPixels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}

// Helper function to access an image's color at the given coordinates
// from the unit interval (e.g. top-left is (0, 0), center is (0.5, 0.5),
// bottom-right is (1.0, 1.0)).
func imageColorAtUnit(unitCoords vec2) vec4 {
	offsetInTexels, sizeInTexels := imageSrcRegionOnTexture()
	adjustedTexelCoords := unitCoords*sizeInTexels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}
```

You can copy `imageColorAtPixel()` into the previous shader in order to complete it. Notice also that if you want to use multiple images, you will need to write additional versions of these helper functions using `imageSrc1At` and the others instead of `imageSrc0At`.

To complete the program, make sure to use the spider-cat-dog bounds for the window and layout sizes in your `main.go` and try to run it! You should be seeing this silly creature, but with a black background:

![](https://github.com/tinne26/kage-desk/blob/main/display/spider_cat_dog.png?raw=true)

Ok, next part! Let's start with a very simple exercise to get you more familiar with working with images: modify the previous shader to display the same image, but with the rgb channels mixed up. For example, try to put the red channel into the green, the green into the blue, and the blue into the red!

<details>
<summary>Click to show the solution</summary>

```Golang
func Fragment(position vec4, _ vec2, _ vec4) vec4 {
	return imageColorAtPixel(position.xy).brga
}
```
*(Full program available at [examples/intro/color-swap](https://github.com/tinne26/kage-desk/blob/main/examples/intro/color-swap))*

Don't tell me you forgot about *swizzling*! So simple and yet so cool!
</details>


### Table of Contents
Next up: [#8](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/08_triangles.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/06_uniforms.md)
7. [**Using images**](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/07_images.md)
8. [`DrawTrianglesShader()`](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/08_triangles.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/09_loops.md)
10. [Tutorial end](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/10_end.md)

[^1]: If you want to understand texels in Ebitengine, we have written an [appendix for it](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/appx_texels.md). We consider it a burden for most people, but it can be useful for optimization or if you can't live without understanding the *why* of things.
