# Learn by example

This is a list of Kage shaders that you may use as references or in order to learn specific techniques. It's recommended that you have some basic knowledge about shaders and / or Kage before diving into this, but we have no bouncer at the door.

The shaders are roughly sorted from simple to complex, and they include tags that indicate what topics they touch on. Here's the full list of tags, which you may use to `ctrl + f` if you are looking for anything in particular:
> geometry, 3D, lighting, triangles, deformation, projection, masking, animation, UI, FX, screen, sprite, image, color, artistic.

> [!NOTE]
> This page is still very much in development. There's already a fair amount of content, but nowhere as much as I'd like. The accessibility and immediacy also need to be severely improved. Please have some leniency.

We have also created a page for [community contributed examples](https://github.com/tinne26/kage-desk/blob/main/docs/community_examples.md). Check those out too and feel free to contribute your own!

Here we go!
- [Filled circle](https://github.com/tinne26/kage-desk/blob/main/examples/learn/filled-circle): geometry.
- [Unfilled circle](https://github.com/tinne26/kage-desk/blob/main/examples/learn/unfilled-circle): geometry.
- [Gamma correction](https://github.com/tinne26/kage-desk/blob/main/examples/learn/gamma-correction): color.
- [Filled triangle](https://github.com/tinne26/kage-desk/blob/main/examples/learn/filled-triangle): geometry.
- [Unfilled triangle](https://github.com/tinne26/kage-desk/blob/main/examples/learn/unfilled-triangle): geometry.
- [Filled rounded rectangle](https://github.com/tinne26/kage-desk/blob/main/examples/learn/filled-rounded-rectangle): geometry, UI.
- [Unfilled rounded rectangle](https://github.com/tinne26/kage-desk/blob/main/examples/learn/unfilled-rounded-rectangle): geometry, UI.
- `TODO` [Filled polygon](): geometry.
- `TODO` [Unfilled polygon](): geometry.
- [Black and white](https://github.com/tinne26/kage-desk/blob/main/examples/learn/black-and-white): image, color.
- [Simple tint](https://github.com/tinne26/kage-desk/blob/main/examples/learn/simple-tint): image, color.
- [HSL hue rotation](https://github.com/tinne26/kage-desk/blob/main/examples/learn/hsl-hue-rotation): image, color.
- [Oklab chroma shift](https://github.com/tinne26/kage-desk/blob/main/examples/learn/oklab-chroma-shift): image, color.
- `TODO` [Pixelize](): FX, image.
- `TODO` [Pixelize (sampled)](): FX, image.
- `TODO` [Blur](): FX, image.
- `TODO` [Grain](): screen, FX, animation.
- `TODO` [Vignette](): screen, FX, ¿animation?.
- `TODO` [Noise](): FX.
- `TODO` [Simple sphere](): geometry, 3D, lighting.
- `TODO` [Complex wave](): geometry, animation.
- `TODO` [Bézier conic curve](): geometry.
- `TODO` [Bézier cubic curve](): geometry.
- `TODO` [Complex polygon](): geometry.
- `TODO` [Simple glow](): sprite, FX.
- `TODO` [Simple outline](): sprite, FX.
- `TODO` [Behind obstruction](): masking, sprite.
- `TODO` [Behind wall](): masking, sprite.
- `TODO` [Dynamic motion blur](): FX, image, animation.
- `TODO` [Breathing](): sprite, animation, FX, deformation.
- `TODO` [Character boost](): sprite, FX, animation.
- `TODO` [Character poison](): sprite, FX, animation.
- `TODO` [Windmill transition](): image, animation.
- `TODO` [Destroy (disintegrate + spread)](): sprite, animation, FX.
- `TODO` [Destroy (cut + vanish)](): sprite, animation, FX.
- `TODO` [Destroy (spiral collapse)](): sprite, animation, FX.
- `TODO` [Destroy (radial crystal)](): sprite, animation, FX.
- `TODO` [Chromatic aberration](): image, FX, color.
- `TODO` [Hit (pixelize + wave)](): sprite, animation, FX.
- `TODO` [Water reflection](): image, FX, projection, deformation, animation.
- `TODO` [Inside water](): image, FX, deformation, animation.
- `TODO` [Heat shimmer](): image, FX, deformation.
- `TODO` [Cube](): geometry, 3D, lighting.
- `TODO` [Animated cube](): geometry, 3D, lighting, animation.
- `TODO` [Ebiten model](): geometry, 3D, lighting, animation.
- `TODO` [Ring loader](): UI, projection, animation.
- `TODO` [Cylindrical HUD](): UI, projection.
- `TODO` [Hemispherical HUD](): UI, projection.
- `TODO` [SMAA (antialiasing)](): screen, FX. https://github.com/iryoku/smaa
- `TODO` [Subpixel antialiasing](): screen, FX, color.
- `TODO` [Simple moon phases](): geometry, masking, animation.
- `TODO` [Geodesic dome](): geometry, 3D, triangles, lighting.
- `TODO` [Rotating geodesic dome](): geometry, 3D, triangles, animation, lighting.

At this point, I feel it's important to highlight that many fancy effects can also be achieved *without shaders*. Keep this in mind and avoid shader tunnel vision!
- **Hit effect**: compress the sprite with regular `GeoM` transformations.
- **Morphing**: saturate both images to white with `ColorM` transformations, and then scale the first down and the second up at the same center position. You can create "pokemon evolution" effects or other kinds of morphs between two sprites.
- **Skews and transitions**: you can use `GeoM` to skew images, combine color blending operators, masks and other tricks to create screen transitions and all kinds of animated effects.
- **Alter states**: it's easy to create flashes of different colors on characters to represent alter states, like poison, paralysis and others. Often you can add custom small animations to make these effects even much better.
- **Manual effects with game art**: outlines, glows, glitches and others can be done directly with art and animations. If you need a cool effect only in one specific situation, sometimes making the art for it may be easier than coding the shader.
