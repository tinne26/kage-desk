# Misc.

Miscellaneous snippets in no particular order.

## Relative target coordinates

Get `targetCoords` in `[0, 1]` range. We all forget the origin:
```Golang
relativeX := (targetCoords.x - imageDstOrigin().x)/imageDstSize().x // [0, 1]
relativeY := (targetCoords.y - imageDstOrigin().y)/imageDstSize().y // [0, 1]
```

If you are going to use both, better get it directly as a `vec2`.
```Golang
relTargetCoords := (targetCoords.xy - imageDstOrigin())/imageDstSize()
```

## Bilinear interpolation

In some cases you will prefer `imageSrcNUnsafeAt()` and clamping or others, but this is the general idea.
Pass `unit = vec2(1)` for default bilinear interpolation or use `unit = fwidth(sourceCoords)` for sharper results.
```Golang
func bilinearSampling(coords, unit vec2) vec4 {
	tl := imageSrc0At(coords - unit/2.0)
	tr := imageSrc0At(coords + vec2(+unit.x/2.0, -unit.y/2.0))
	bl := imageSrc0At(coords + vec2(-unit.x/2.0, +unit.y/2.0))
	br := imageSrc0At(coords + unit/2.0)
	delta  := min(fract(coords + unit/2.0), unit)/unit
	top    := mix(tl, tr, delta.x)
	bottom := mix(bl, br, delta.x)
	return mix(top, bottom, delta.y)
}
```

## Texture clamping / repeat

If you use `imageSrcNAt()`, coordinates outside bounds will return `vec4(0)` by default. If you use `imageSrcNUnsafeAt()`, there's no clamping at all and you might access the whole underlying internal atlas. Here are some snippets for other behaviors:

> [!NOTE]
> *Notice that default Ebitengine shaders for `DrawImage()` and `DrawTriangles()` already support the [`Address`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Address) type, which allows changing between unsafe, clamp and repeat. The snippets below are the manual version of this for custom shaders.*

Extend edge color:
```Golang
//kage:unit pixels
package main

var Cursor vec2 // in [0, 1] ranges, built-in with kage-desk/display

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	minSrc0, maxSrc0 := GetSource0ClampCoords()
	clampedCoords := clamp(sourceCoords + (Cursor - vec2(0.5))*100.0
	return imageSrc0UnsafeAt(clampedCoords)
}

func GetSource0ClampCoords() (vec2, vec2) {
	const epsilon = 1.0/65536.0 // TODO: how small can we safely set this?
	origin := imageSrc0Origin()
	return origin, origin + imageSrc0Size() - vec2(epsilon)
}
```

**TODO: UNTESTED** Wrap:
```Golang
//kage:unit pixels
package main

var Cursor vec2 // in [0, 1] ranges, built-in with kage-desk/display

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	originSrc0, endSrc0 := GetSource0Region()
	adjustedCoords := sourceCoords + (Cursor - vec2(0.5))*100.0
	return imageSrc0UnsafeAt(Wrap(adjustedCoords, originSrc0, endSrc0))
}

func GetSource0Region() (vec2, vec2) {
	origin := imageSrc0Origin()
	return origin, origin + imageSrc0Size()
}

func Wrap(coords vec2, origin, end vec2) vec2 {
	return origin + mod(coords - origin, end - origin)
}
```

## Remap `sourceCoords` between images

The `sourceCoords` passed to `Fragment(...)` are based on the source image 0. If you use multiple images in the same shader —especially if they are of different sizes—, you will need to remap these coordinates... but this is done in a rather surprising way:
```Golang
normCoords0   := (sourceCoords - imageSrc0Origin())/imageSrc0Size()
sourceCoords1 := (normCoords0*imageSrc1Size()) + imageSrc0Origin()
```
Yep, that's not a typo, we have to *add the origin of `Source0`*.

To be fair, the isolated normalization of coordinates can be quite useful on its own, so here you go:
```Golang
func normSrc0Coords(sourceCoords vec2) vec2 {
	return (sourceCoords - imageSrc0Origin())/imageSrc0Size()
}

func normDstCoords(targetCoords vec4) vec2 {
	return (targetCoords.xy - imageDstOrigin())/imageDstSize()
}
```

## Math functions

**TODO: UNTESTED**
```Golang
func fmod(value, modulo float) float {
	return value - modulo*trunc(value/modulo)
}

func trunc(a float) float {
	if a >= 0 { return floor(a) }
	return -floor(-a)
}
```

## Colors

Shouldn't stay there after you optimize, but occasionally helpful for testing colors:
```Golang
func rgb(r, g, b int) vec4 {
	return vec4(float(r)/255.0, float(g)/255.0, float(b)/255.0, 1.0)
}
```
*(TODO: many useful functions that we might want to copy from black-and-white and HSL-hue-shift examples)*

## Aiding SDF visualization

```Golang
func colorwaveA(sdf float) vec4 {
	mid := vec4(0.156, 0.682, 0.501, 1.0)
	dark, light := vec4(0.172, 0.447, 0.556, 1.0), vec4(0.368, 0.788, 0.384, 1.0)
	out := mix(mid, dark , 0.85 + 0.5*sin(clamp(pow(sdf*+0.25, 1.3), 0, 1_000_000)))
	in  := mix(mid, light, 0.85 + 0.5*sin(clamp(pow(sdf*-0.25, 1.4), 0, 1_000_000)))
	return out*clamp(+sdf, 0, 1) + in*clamp(-sdf, 0, 1)
}

func colorwaveB(sdf float) vec4 {
	out := mix(vec4(0.941, 0.227, 0.278, 1.0), vec4(1.000, 0.713, 0.152, 1.0), 0.1*sin(sdf*0.02))
	in  := mix(vec4(0.223, 0.486, 0.345, 1.0), vec4(0.352, 0.901, 0.980, 1.0), 0.85 + 0.3*sin(sdf*0.6))
	presenceOut := clamp(1.05 + 0.20*sin(sdf*0.44), 0, 1)*clamp(+sdf, 0, 1)
	presenceIn  := clamp(1.15 + 0.32*sin(sdf*0.83), 0, 1)*clamp(-sdf, 0, 1)
	return out*presenceOut + in*presenceIn
}
```


