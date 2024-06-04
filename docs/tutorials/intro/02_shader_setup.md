# Setting up your first shader

With the introduction out of the way, now it's time to get our hands dirty. If you want to learn shaders, you need to write shaders!

> [!IMPORTANT]
> *The entirety of kage-desk expects you to have Go1.19 or higher installed, and it also depends on Ebitengine's "pixel mode" for kage, which was added in version v2.6.0. All these have been around for a good while, so everything should be fine by default... but better safe than sorry.*

First, create a folder somewhere on your messy computer, run `go mod init first-shader` within it, and create a `main.go` with this content:
```Golang
package main

import "github.com/tinne26/kage-desk/display"

func main() {
	display.SetTitle("intro/first-shader")
	display.SetSize(512, 512)
	display.Shader("shader.kage")
}
```

Don't worry about the meaning of this `main.go` file for the moment, it's only helper code to get our shaders running more easily at the beginning.

Create also a `shader.kage` file with the following content:
```Golang
//kage:unit pixels
package main

func Fragment(_ vec4, _ vec2, _ vec4) vec4 {
	// ...
}
```
> [!TIP] *You can [configure your editor](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/config_editor.md) to highlight `.kage` files like `.go` programs)*

> [!CAUTION]
> *Notice that the `kage:unit pixels` is not just a comment, but a special directive in the same style as [Golang's compiler directives](https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives). Do not remove or modify it[^1].*

[^1]: The `kage:unit pixels` directive exists for compatibility reasons. Originally, Ebitengine used the "texels" mode for shaders, but after some time we saw that everything was kinda simpler with pixels, and a new mode was added. The default mode couldn't be changed without breaking compatibility, so the special directive was the best solution in the meantime. This new mode —the "pixels" mode— will become the default in Ebitengine v3, but that's not expected to arrive until at least 2026. There is a lot we could say about pixel vs texel modes, but all you probably need to know is that all advanced and long-time Ebitengine users have quickly switched to pixel mode.

In your typical CPU programs, the entry point is the `main()` function. In Kage, the entry point is the `Fragment()` function instead. The reason the entry function is called "fragment" is because there are multiple types of shader programs: vertex shaders, geometry shaders, compute shaders, tessellation shaders... but you don't need to remember those now, as Ebitengine only has fragment shaders. Fragment shaders, also called pixel shaders, are shaders that compute the color of a single fragment or pixel.

Anyway, back to work. We don't know what we want to do yet, but let's just execute the program anyway! `go run main.go`!

Oh... it didn't work? *Hmph*.
```
Failed to load shader:
3:1: function Fragment must have a return statement but not
```

Well, as we were just saying, a fragment shader should return the value of a pixel. Its color. We aren't returning anything yet, so let's fix the `Fragment` function by adding a `return vec4(1, 0, 0, 1)`. If you did it right, re-running `main.go` should now result in a red screen.

> [!NOTE]
> *If you are having any problems, the full working code for this first example can be found at [examples/intro/first-shader](https://github.com/tinne26/kage-desk/blob/main/examples/intro/first-shader).*

You probably figured it out on your own, but as you can see, we don't use `color.RGBA` within shaders, but vectorial types instead. Check the function signature again:
```Golang
func Fragment(_ vec4, _ vec2, _ vec4) vec4
```

We have a few input vectors of different sizes that we are ignoring for the moment, and one output `vec4`. That output vector is the color of the pixel! The vector has the four components that you may expect: red, green, blue and alpha.

Vectors in Kage are made of `float` values, which are 32-bit precision. In Golang we have both `float64` and `float32`; in kage we only have `float`. Using float values also means that unlike `color.RGBA`, `vec4` components are expected in the [0, 1] range (instead of [0, 255]). Keep this in mind because it's very easy to accidentally mix it up! If you use values outside that range in the returned color, they will be clamped to 0 - 1.

These vector types are actually quite cool, and you can do many weird operations with them:
```Golang
vec4(1) // creates a vector with 4 components, all set to 1
vec4(0, 0, 16, 32)[3] // access the last component of a vec4 with normal indexing
vec4(0.5, 0.5, 1.0, 1.0).rgb // gets a vec3 with the r, g, and b components from the vec4
vec3(32, 44.0, 0).xy // fields can't only be accessed as rgba, but also xyzw or stpq
vec2(3, 5).yxx // you can even mix up the order or repeat fields!
```
This field access magic thing is known as "swizzling". I'll test you again in a few chapters, so you better remember it!


### Table of Contents
Next up: [#3](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_position_input.md).

0. [Introduction](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/00_introduction.md)
1. [CPU vs GPU: different paradigms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/01_cpu_vs_gpu.md)
2. [**Setting up your first shader**](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/02_shader_setup.md)
3. [The `targetCoords` input parameter](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/03_target_coordinates.md)
4. [Built-in functions](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/04_built_in_functions.md)
5. [Manual shader invocation](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/05_invoke_shader.md)
6. [More input: uniforms](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/06_uniforms.md)
7. [Using images](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/07_images.md)
8. [Beyond one-to-one mapping](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/08_beyond.md)
9. [Loops are tricky](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/09_loops.md)
10. [What's next?](https://github.com/tinne26/kage-desk/blob/main/docs/tutorials/intro/10_what_next.md)
