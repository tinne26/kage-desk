# What's next?

First of all, congratulations on getting this far! There are still so many things you will be unsure about, and yet, I hope that now you feel more ready to start experimenting and exploring on your own.

Depending on what's your motivation to learn shaders, many big and important topics might have been omitted. I'll add a brief list with quick explanations here, not as proper tutorials, but to make you aware of some ideas you might not have considered yet, or as possible suggestions of where to look next.

The third `color` parameter:
> We haven't discussed the third parameter of `Fragment(targetCoords vec4, sourceCoords vec2, color vec4) vec4`. This is actually mentioned in passing in the [triangles](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/triangles.md) article, but the basic idea is that you can set a color for each vertex you pass to `DrawTrianglesShader()` and access that value as an input parameter.
>
> This is specially useful when creating color related shaders that will darken or recolor the image, but it can also be used for hackier purposes. The idea of passing additional data to vertices is actually a general concept; for example, in GLSL (OpenGL's shading language), we call this 'vertex attributes'. There are some open issues in Ebitengine [discussing vertex attributes](https://github.com/hajimehoshi/ebiten/issues/2640).

Additional exercises for practice:
> If you want to practice more, we have the whole [learn-by-example](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/learn_by_example.md) section in kage-dest, which is awfully underdeveloped, but... it contains a lot of interesting ideas that you can use to practice! In fact, if you do really well on any of the problems, you could even consider contributing some code back..?

Interpolation methods:
> One of the biggest problems that we didn't discuss in the tutorial and that you will find in practice is that, most often, your shaders need to be able to adapt to different sizes and resolutions. When this happens, we often need to apply interpolation algorithms... preferably something better than nearest neighbour.
>
> I would like to write more about the topic and probably create some nice snippets for everyone, but for the moment the best I can do is point you towards [`mipix/filters`](https://github.com/tinne26/mipix/tree/main/filters), which contains a fair amount of interpolation and filtering methods: bilinear, bicubic, hermite, pixel-art sampling, etc. You can actually test some of the behaviors online at [mipix-examples/gametest](https://tinne26.github.io/mipix-examples/gametest/). Most of these shaders use `SourceRelativeTextureUnitX` and `SourceRelativeTextureUnitY` uniforms, but in most cases you can derive those out directly with `units := fwidth(sourceCoords)` or other tricks.

Use of multiple images:
> Using multiple source images is very common in 3D due to the use of heightmaps, lightmaps and so on. In 2D it's more uncommon, but you can definitely mix images and textures and do hacky stuff for more advanced use-cases. Using multiple images is actually quite straight-forward, so I don't know if a tutorial would even be that helpful... but one interesting point is that since Ebitengine v2.8.0, it's possible to use source images of different sizes in a single shader. Using multiple images of different sizes gets a bit trickier, since the `sourceCoords` only refer to image `0`, not all of them... but there are many ways around this, and it depends mostly on your particular use-case.

Debugging tips:
> We do a lot of binary searches with `if value >= X && value <= Y { return vec4(1, 0, 0, 1) }` to figure out stuff. Some early returns too. I can't pretend it gets much more glamorous than that, sorry.

Signed distance fields:
> There are a lot of cool things you can do with shaders, but using maths to draw all kinds of geometry is pretty high up the list. You can just google around and find an introduction to the topic in your favorite format and platform. Afterwards, you should check out Íñigo Quilez's legendary work. Here's a [video showcasing what you can do with SDFs](https://youtu.be/8--5LwHRhjk), and here's a link to a [collection of SDFs functions on his site](https://iquilezles.org/articles/distfunctions/). You might accidentally get trapped on his site unable to escape for a few days, but as far as rabbit holes go, this is a good one to get lost in.


## Feedback

[Discussions](https://github.com/tinne26/kage-desk/discussions) are open for this repository, so if you want to give any feedback, ask additional questions, suggest improvements or anything else, feel free to open a discussion and let me know. The discord server is also very active if you want to join and share something with the community.

Otherwise, it might be time to return to the [main page](https://github.com/tinne26/kage-desk)!


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
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [**What's next?**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
