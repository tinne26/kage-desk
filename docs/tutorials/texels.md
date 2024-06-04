## Texels in Kage

> [!WARNING]
> **Since the addition of pixels mode in Ebitengine v2.6.0, the content of this page is no longer relevant; it's offered only for historical context**.

From working with Ebitengine you will know that its graphical API is based around "images". Images are simply arrays of pixels with 4 channels (red, green, blue, alpha) per pixel. Images also have an integer width and height, in pixels.

Then, when it comes to drawing, Ebitengine takes care of sending this data to the GPU.

When we are working with shaders we can also specify up to 4 images per shader (all of the same size), to be used for our own nefarious purposes.

Once the image data is in the video memory of the GPU, we usually refer to them as **textures**.

And here is where **texels** first appear. When we want to access the data of these images from within a shader, we don't use pixel-based coordinates, but texel-based coordinates. As the Ebitengine shaders official document says:
> A pixel is a unit for one color dot. On the other hand, a texel is a unit covering the whole texture area with values between 0 and 1.

In other words: a texel is the area taken by a texture, normalized to cover the unit interval [0, 1]. If our texture was 511x511 and we wanted to access the color of the pixel at its center, instead of accessing `(255, 255)`, with texels we would access `(0.5, 0.5)` instead. This would still apply even if the aspect ratio of the texture was not square.

Image access built-in functions for Kage (like `imageSrc0At(texelCoords vec2) vec4`) expect texel coordinates.

Wait... if this is so easy, why was it relegated to an appendix?!

### Because it's not so simple

We have explained the basic theory, but... it only works if your **image size matches the atlas size**. Furthermore, due to the requirements of old OpenGL versions, the size of Ebitengine atlases are always powers of 2.
> Atlases are big images that contain other smaller images within them. When you create an image, Ebitengine automatically puts it into an atlas, which makes it more efficient to send the data to the GPU.

In practice, this means that the images passed to shaders have to take into account the **padding and offsets** required to locate them inside the atlas. The texel doesn't cover the area of your original image, but the area of the altas instead... and finding the texel coordinates of your actual image within the atlas container *is a bit of pain*. We will be calling these new coordinates **adjusted texels** for educative purposes.

This is why we offer two helper functions to work with textures, that we can finally explain in more detail. The first is this one:
```Golang
// Helper function to access an image's color at the given coordinates
// from the unit interval (e.g. top-left is (0, 0), center is (0.5, 0.5),
// bottom-right is (1.0, 1.0)).
func imageColorAtUnit(unitCoords vec2) vec4 {
	offsetInTexels, sizeInTexels := imageSrcRegionOnTexture()
	adjustedTexelCoords := unitCoords*sizeInTexels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}
```

You will notice that we used the "unit interval" terminology in order to avoid texels and having to explain all this mess to newcomers, but you could rename it `imageColorAtTexel` if you wanted.

The key idea in this code is that the built-in function `imageSrcRegionOnTexture()` gives us the offset and the size of the relevant part of the atlas texture that contains our image, both in texels. With this, we can project our "ideal texel coordinates" or `unitCoords`, that fall between 0 and 1, to this specific sub-region of adjusted texels. In the best case, when our image size is a power of 2 and matches the atlas size, `originInTexels` will be 0, `sizeInTexels` will be 1, and the `adjustedTexelCoords` will match our "ideal" texel or `unitCoords`.

The other helper function expects pixels directly:
```Golang
func imageColorAtPixel(pixelCoords vec2) vec4 {
	sizeInPixels := imageSrcTextureSize()
	offsetInTexels, _ := imageSrcRegionOnTexture()
	adjustedTexelCoords := pixelCoords/sizeInPixels + offsetInTexels
	return imageSrc0At(adjustedTexelCoords)
}
```

The conversion is very similar, but here we have to convert first from pixels to texels using `imageSrcTextureSize()`. Since the result of this built-in function is the containing atlas size (and not out actual shader image size), the resulting texels from the following calculations are already "adjusted", which is why we only need to apply the `offsetInTexels` but we can ignore the `_ = sizeInTexels` result of `imageSrcRegionOnTexture()`. The expanded expression would be:
```Golang
adjustedTexelCoords := (pixelCoords*sizeInTexels)/(sizeInPixels*sizeInTexels) + offsetInTexels
```
Where `sizeInTexels` cancels out, simplifying to `pixelCoords/sizeInPixels`.

### Avoid texels unless optimizing

Now you should understand why we prefer avoiding the concept of texels in the tutorials and recommend using the helper functions instead. That being said, it's true that if you want to get the most from your shaders, you may want to consider the following when it comes to optimization:
- Getting `imageSrcTextureSize()` and `imageSrcRegionOnTexture()` values only once and reusing them when you need multiple samples from the texture within a single shader (e.g. blurs and other [image kernels](https://setosa.io/ev/image-kernels/)).
- Making use of `imageSrc0UnsafeAt` instead of the slower safe version if you know that you are operating strictly within the relevant texels interval.
- Sometimes you can use the second input argument in `Fragment()`, which is the `texelCoords vec2` corresponding to the fragment. This works well for many `DrawRectShader()` invocations, but we don't even explain it in the regular tutorial because then we would also have to explain *when* it doesn't work and *why*, and that's... this document.

Future versions of Ebitengine [will probably rework how texels and access to images in shaders work](https://github.com/hajimehoshi/ebiten/issues/1431).
