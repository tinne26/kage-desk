# Introduction to Kage

Kage is a programming language to implement shaders in Ebitengine.

If you are new to shaders, the short version is that they are programs that run on the GPU instead of the CPU. For our purposes, shaders are programs that allow us to create or modify images[^1], like recoloring them, adding noise or grain, distorting them and many others:

![](https://github.com/tinne26/kage-desk/blob/main/img/zombie_aliens_recolor.webp?raw=true)
*(Shader recolor based on lightness, with gamma correction applied)*

In games, you will see shaders being used for all kinds of things: rounded rectangles on UI frameworks, hit and death animations on sprites, power-up or alter state effects, full screen effects like blurs, warpings, chromatic aberrations and CRT effects, lighting, water generation, screen transitions... and many, many more.

<!-- TODO: video of a few cool shaders to make the contextualization stick -->

Shaders, in conclusion, are programs that allow the GPU to perform sophisticated computations for the individual pixels of an image in a highly parallel manner. When this process is repeated frame after frame —as you might have inferred from the previous examples—, we can even create cool animations and effects.

There are a few different languages in which shaders can be written: you may have heard of GLSL, HLSL and others. Ebitengine has it's own intermediate language, Kage, which allows us to write shaders with a Golang-like syntax without having to worry about other internal details. At runtime, Ebitengine will translate our Kage program to HLSL, MSL or whatever language is needed to make it work on the platform where the game's being played.

### Table of Contents
Next up: [#1](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md).

0. [**Introduction**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [`DrawTrianglesShader()`](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_triangles.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)

[^1]: Shaders can also be used for general computation, not just graphics, but that's outside the scope of this guide. We have thrown a few links [here](https://github.com/tinne26/kage-desk/blob/main/docs/general_links.md) if you want to learn more later.
