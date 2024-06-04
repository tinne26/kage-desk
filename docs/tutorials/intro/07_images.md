# Using images

It's time to reveal one of the most important pieces left of Kage! If we want to unlock the true power of shaders, we will need... *images*.

There are two main reasons why we might want to pass images to our shaders:
- To **apply an effect** to our game sprites or screen. This is the most common way to use images in 2D games. Pixelize an image, apply a blur or movement blur effect to it, deform it, control its color, apply chromatic aberrations, make simple lighting effects, screen transitions, etc.
- To **use the image as a texture** for the shader. This is very common in 3D, where textures are used for "painting" raw triangles, geometry and lighting purposes. In 2D this is more unusual, but there are still some use-cases like creating animation effects on a sprite (e.g. being electrocuted), fancy glitches, reflections on the water, making "see behind the wall" effects and a few others that combine multiple images or textures to achieve a specific effect. Advanced lighting techniques with surfaces and normals can also be used in 2D, mainly in top-down view games, but this is rather uncommon, so we won't discuss it in this tutorial.

> [!NOTE]
> *The words _image_ and _texture_ are often used interchangeably in the context of shaders. There are some nuances, but you can basically consider them fully equivalent.*

To get started, we will show how to pass an image to a shader in `main.go`: the draw options struct includes an `Images` array where we can set up to 4 images to be passed to the shader:
```Golang
func (self *Game) Draw(screen *ebiten.Image) {
	// ... (some stuff)

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Images[0] = display.ImageSpiderCatDog()

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```
> [!TIP]
> *We are importing `github.com/tinne26/kage-desk/display` again as it includes some images that we can use for our tests quite easily. You can try loading your own images if you want.*

Hmmmm... ok, but... there's something missing here.

Now that we have an input image, how do we map it to the target? Like, the target and the source[^1] can have different sizes, so how do we tell the shader what part of the source texture maps to the target?

[^1]: Quick reminder in case you are a bit lost with the "target" (also known as "destination") and "source" terminology, as this applies to Ebitengine as a whole, not just shaders: the texture or image that we are modifying is the "target". Sometimes we are simply filling it with a color, gradient or noise, and we don't have any other element into the equation. More commonly though, we modify targets by blending graphical data from another image, the "source(s)". The basic operation that does this in Ebitengine is `target.DrawImage(source, options)`. Knowing what's a target and what's a source is not only important conceptually, but [also for performance](https://ebitengine.org/en/documents/performancetips.html#Avoid_changing_render_sources'_pixels), as Ebitengine can operate much more efficiently if sources and targets are clearly differentiated and they consistently reside in separate internal atlases. In the case of shaders, the only big novelty is that you can have *multiple* source images for a single draw. Don't worry, you will internalize all this in due time.

Well, we only need to set the `SrcX` and `SrcY` fields of the [vertices](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Vertex). Vertices have two pairs of coordinates:
- `DstX` and `DstY`, which indicate the *target* coordinates for the vertex.
- `SrcX` and `SrcY`, which indicate the *texture sampling* coordinates for the source images that we pass to the shader.

Again, this might be tricky to visualize with only an explanation, but the [triangles](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/triangles.md) article goes a bit more in depth if you need additional explanations. Here's an image from that article that might help you visualize everything:

![](https://github.com/tinne26/kage-desk/raw/main/img/triangle_B.png?raw=true)

So, *in order to map our source texture to the target region*, we get the source image bounds and set the `Vertex.SrcX` and `Vertex.SrcY` fields like this:
```Golang
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

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	screen.DrawTrianglesShader(self.vertices[:], indices, self.shader, &shaderOpts)
}
```

Oof. Ok, that was tedious, but we had to deal with it at some point. Let's go back to the shaders now!

For the shader, the first thing we will try do is show the image. No effects yet. Just make the shader compute the color of each pixel as the color of the corresponding pixel in the passed image:
```Golang
//kage:unit pixels
package main 

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	return imageSrc0At(targetCoords.xy)
}
```

This code showcases a new function: <code>imageSrc<b><i>N</i></b>At()</code>. This function allows us to get the color of a source image at a given position. We can have up to four images, and you can use `imageSrc0At()`, `imageSrc1At()`, `imageSrc2At()` and `imageSrc3At()` to sample colors from each one. There are also [a few more functions](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(images)) to get a source image size, origin and so on, but we can deal with that another day.

You should try to run all this now... and notice that the image doesn't properly fill the screen.

Of course: `targetCoords` refer to the destination image; what we need are the *source texture sampling coordinates* corresponding to the current `targetCoords`! Luckily enough, that's actually the second argument to `Fragment(...)` that we hadn't unveiled yet:
```Golang
//kage:unit pixels
package main 

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	return imageSrc0At(sourceCoords)
}
```
*(Full program available at [examples/intro/spider-cat](https://github.com/tinne26/kage-desk/blob/main/examples/intro/spider-cat))*

Now you could even add `ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)` to the main program and test that no matter the size of the screen, the spider-cat-dog will resize alongside it. You might want to fill the screen white first so it's easier to tell how much space the automatic Ebitengine borders take.

> [!TIP]
> *You might have noticed that the spider-cat-dog doesn't look particularly smooth. This is because by default, the sampling uses nearest neighbour interpolation, instead of bilinear or something else. This is outside the scope of this tutorial, but just know that bilinear interpolation would be quite easy to add to our shader to make the results much smoother. In the meantime, you could fix the window size and layout to 384x384 instead, which is the exact size of the spider-cat-dog.*

If you managed to put it all together, you should see something similar to this:

![](https://github.com/tinne26/kage-desk/blob/main/display/spider_cat_dog.png?raw=true)

Well done! Let's end this chapter with a very simple exercise to get you more familiar with working with images: modify the previous shader to display the same image, but with the rgb channels mixed up. For example, try to put the red channel into the green, the green into the blue, and the blue into the red!

<details>
<summary>Click to show the solution</summary>

```Golang
//kage:unit pixels
package main

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	return imageSrc0At(sourceCoords).brga
}
```
*(Full program available at [examples/intro/color-swap](https://github.com/tinne26/kage-desk/blob/main/examples/intro/color-swap))*

Don't tell me you forgot about *swizzling*! I told you in chapter 2 that I would ask again! So simple and yet so cool!
</details>


### Table of Contents
Next up: [#8](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [**Using images**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
