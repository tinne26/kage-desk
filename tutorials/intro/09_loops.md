# Loops are tricky

**WIP**

One thing we haven't discussed yet are loops. At the start of the tutorial we showed that since shaders are run in parallel for many pixels, they are like the core part of the loop in an equivalent CPU program.

Adding loops within shaders, then, is like having a loop inside a loop inside a loop inside... well, something like that.

This is not the most common thing, but it's not super uncommon either. If you have programmed for a while, you know how nesting absurd levels of loops are a thing, even if you try to modularize it into separate functions or whatever.

One good example of a shader that can be done using a loop is the *pixelization effect*. So this will be your next exercise:

![](https://github.com/tinne26/kage-desk/blob/main/img/pixelated_creature.png?raw=true)

*TODO: base demo code, talk about sampling only 1 pixel, sampling multiple pixels, uniforms vs constants, etc. If we trained them right, they should be able to do most of the setup with only a few indications.*
