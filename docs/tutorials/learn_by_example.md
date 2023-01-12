# Learn by example

**TODO: none of the examples are written yet, this is a "what would be awesome to have" list**

This is a list of Kage shaders that you may use as references or in order to learn specific techniques. It's recommended that you have some basic knowledge about shaders and / or Kage before diving into this, but we have no bouncer at the door.

The shaders are roughly sorted from simple to complex, and they include tags that indicate what topics they touch on. Here's the full list of tags, which you may use to `ctrl + f` if you are looking for anything in particular:
> geometry, 3D, lighting, triangles, deformation, projection, masking, animation, UI, FX, screen, sprite, image, color, artistic.

Here we go!
- [Filled circle](): geometry.
- [Unfilled circle](): geometry.
- [Gamma correction](): color.
- [Filled smooth triangle](): geometry.
- [Unfilled smooth triangle](): geometry.
- [Filled rounded rectangle](): geometry, UI.
- [Unfilled rounded rectangle](): geometry, UI.
- [Filled polygon](): geometry.
- [Unfilled polygon](): geometry.
- [Black and white](): image, color.
- [Simple tint](): image, color.
- [Simple sphere](): geometry, 3D, lighting.
- [Complex wave](): geometry, animation.
- [Blur](): FX, image.
- [Pixelize](): FX, image.
- [Pixelize (sampled)](): FX, image.
- [Bézier conic curve](): geometry.
- [Bézier cubic curve](): geometry.
- [Complex polygon](): geometry.
- [Grain](): screen, FX, animation.
- [Vignette](): screen, FX, ¿animation?.
- [Simple glow](): sprite, FX.
- [Simple outline](): sprite, FX.
- [Behind obstruction](): masking, sprite.
- [Behind wall](): masking, sprite.
- [Dynamic motion blur](): FX, image, animation.
- [Breathing](): sprite, animation, FX, deformation.
- [Character boost](): sprite, FX, animation.
- [Character poison](): sprite, FX, animation.
- [Destroy (disintegrate + spread)](): sprite, animation, FX.
- [Destroy (cut + vanish)](): sprite, animation, FX.
- [Destroy (spiral collapse)](): sprite, animation, FX.
- [Destroy (radial crystal)](): sprite, animation, FX.
- [Chromatic aberration](): image, FX, color.
- [Hit (pixelize + wave)](): sprite, animation, FX.
- [Water reflection](): image, FX, projection, deformation, animation.
- [Inside water](): image, FX, deformation, animation.
- [Heat shimmer](): image, FX, deformation.
- [Ring loader](): UI, projection, animation.
- [Cylindrical HUD](): UI, projection.
- [Hemispherical HUD](): UI, projection.
- [SMAA (antialiasing)](): screen, FX. https://github.com/iryoku/smaa
- [Subpixel antialiasing](): screen, FX, color.
- [Simple moon phases](): geometry, masking, animation.
- [Geodesic dome](): geometry, 3D, triangles, lighting.
- [Rotating geodesic dome](): geometry, 3D, triangles, animation, lighting.

At this point, I feel it's important to highlight that many fancy effects can also be achieved *without shaders*. Keep this in mind and avoid shader tunnel vision!
- Hit effect: compress the sprite with regular `GeoM` transformations.
- Morphing: saturate both images to white with `ColorM` transformations, and then scale the first down and the second up at the same center position. You can create "pokemon evolution" effects or other kinds of morphs between two sprites.
- Skews and transitions: you can use `GeoM` to skew images, combine color blending operators, masks and other tricks to create screen transitions and all kind of animated effects.
- Alter states: it's easy to create flashes of different colors on characters to represent alter states, like poison, paralysis and others. Often you can add custom small animations to make these effects even much better.
- Many other effects done manually with game art: outlines, glows, glitches and others can be done directly with art. If you need a cool effect only in one specific situation, making the art for it may be easier than coding the shader.
