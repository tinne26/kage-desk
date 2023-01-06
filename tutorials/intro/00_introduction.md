# Introduction to Kage

Kage is a programming language to implement shaders in Ebitengine.

If you are new to shaders, the short version is that they are programs that run on the GPU instead of the CPU. For our purposes, shaders are programs that allow us to create or modify images[^1], like recoloring them, adding noise or grain, distorting them and many others:

![](https://github.com/tinne26/kage-desk/blob/main/img/zombie_aliens_recolor.webp?raw=true)
*Simple recoloring example based on lightness.*

Shaders are programs that allow us to perform sophisticated computations for the individual pixels of an image in a highly parallel manner.

There are a few different languages in which shaders can be written: you may have heard of GLSL, HLSL and others. Ebitengine has it's own intermediate language, Kage, which allows us to write shaders in a Golang-like syntax and forget about the rest. At runtime, Ebitengine will translate that Kage program to HLSL, MSL or whatever language is needed to make it work on the platform where the game's being run.

### Table of Contents
Next up: [#1](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md).

0. [**Introduction**](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/01_cpu_vs_gpu.md)
2. [Setting up your first shader](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/02_shader_setup.md)
3. [The `position` input parameter](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/03_position_input.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/tutorials/intro/07_images.md)
8. [Screen vs sprite effects]()
9. [Performance considerations]()
10. [Graduation challenges]()

[^1]: Shaders can also be used for general computation, not just graphics, but that's outside the scope of this guide. We have thrown a few links [here](https://github.com/tinne26/kage-desk/blob/main/tutorials/general_links.md) if you want to learn more later.
