# Kage for game devs

Quick reference for developers with previous knowledge of how shaders work.

Kage is the language used on Ebitengine to write shaders. It has a Golang-like syntax and it's internally translated to GLSL, HSL or MSL at runtime as required. Only fragment shaders are supported at the moment.
- Supported types: `bool`, `int`, `float` (float32), `vec2`, `vec3`, `vec4` (float vectors), `mat2`, `mat3`, `mat4`.
- Vectors support swizzling with `rgba`, `xyzw` and `stpq`.
- You can write helper functions, but there are no slices, maps, strings, structs, import, switch, etc.
- Textures limited to 4 RGBA images of the same size per shader invocation.

#### Table of contents
- [Built-in functions](#built-in-functions)
- [Basic examples](#basic-shaderkage-example)
- [Load and invoke]()
- [Uniforms]()
- [Textures]()

## Built-in functions

Almost all of them can be applied both to numerical types like `int` and `float`, but also vectors. E.g.: `abs(vec2(-1, 0)) == vec2(1, 0)`.

Key single-argument functions:
```Golang
len(vec) // for vec2, vec3, vec4. same as in Golang
length(x) // mathematical length / magnitude of a vector
abs(x)
sign(x) // returns -1, 0 or 1
sin(x), cos(x), tan(x) // plus asin, acos, atan, etc.
sqrt(x)
floor(x), ceil(x), fract(x)
```

Key two-argument functions:
```Golang
mod(x, m) // %
min(a, b), max(a, b)
pow(x, exp)
step(s, x) // 0 if `x < s`, 1 otherwise
dot(x, y), cross(x, y vec3) // dot and cross products
distance(pointA, pointB) // == length(pointA - pointB)
```

Key three-argument functions:
```Golang
clamp(x, min, max)
mix(a, b, t) // lerp, linear interpolation
```

Full [official reference](https://ebitengine.org/en/documents/shader.html#Built-in_functions_(mathematics)).

## Basic `shader.kage` example


```Golang
package main

var Time float // [[ time, 60Hz ]]
var Size vec2 // [[ 400x400 * DeviceScale ]]
const Size vec2 // [[ 200x400 ]]

// TODO: try how can everything appear. uniforms are always uppercase, let's consider utf8 letters
//       though. let's see what to link. implement it. have slices of logic blocks for each param type.

func Fragment(pos4 vec4, _ vec2, _ vec4) vec4 {
	position := pos4.xy // z and w are always 0


}

```

Make sure to [configure your editor](https://github.com/tinne26/kage-desk/blob/main/tutorials/config_editor.md) so you get syntax highlight.

For quick testing, you can put the code into a `shader.kage` file and run it with the following `main.go`:
```Golang
package main
import _ "embed"
import "github.com/tinne26/kage-desk/tools/display"

//go:embed shader.kage
var source []byte
func main() {
	display.Shader(source, "kage/sine", 512, 512)
}
// TODO: autodetect some uniforms by parsing the program?
```
