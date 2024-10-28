# Snippets

A collection of snippets and common functions for Kage.

#### Table of contents
- [Logic](#logic)
	- [Selectors](#selectors)
	- [Area masks](#area-masks)
- [Sampling, coords and interpolation](#sampling-coords-and-interpolation)
	- [Relative target coordinates](#relative-target-coordinates)
	- [Remap `sourceCoords` between images](#remap-sourcecoords-between-images)
	- [Texture clamping and repeat](#texture-clamping-and-repeat)
	- [Bilinear interpolation](#bilinear-interpolation)
- [Misc.](#misc)
	- [Math functions](#math-functions)
	- [Colors](#colors)
	- [Aiding SDF visualization](#aiding-sdf-visualization)

## Logic

### Selectors

A set of functions for keeping/discarding results based on simple conditions[^1]:

[^1]: This can actually be done directly with conditionals most of the time, but some programmers prefer to avoid these in their shaders. The reason is quite long-winded and complex... Modern GPUs are much better optimized now, but in the past, conditionals were suboptimal in many cases. Since the same program is executed for a wavefront (group of pixels), whenever any of the pixels takes a different conditional branch, this branch must still be evaluated for all the pixels in the wavefront. If it's a short `if`, this is not a problem. If the branches are divergent and start doing a lot of completely distinct work, then this can have much more severe performance repercussions.

```Golang
// Returns 1 if a == b, 0 otherwise.
func whenEqual(a, b float) float {
	return 1.0 - abs(sign(a - b))
}

// Returns 1 if a >= b, 0 otherwise.
func whenGreaterOrEqualThan(a, b float) float {
	return step(b, a)
}

// Returns 1 if a < b, 0 otherwise.
func whenLessThan(a, b float) float {
	return 1 - step(b, a)
}

// Returns 1 if a > b, 0 otherwise.
func whenGreaterThan(a, b float) float {
	return 1 - step(a, b)
}

// Returns 1 if a <= b, 0 otherwise.
func whenLessOrEqualThan(a, b float) float {
	return step(a, b)
}
```

A/B'ing with selector results can also be made nicer with some helper functions:
```Golang
// Returns a if selector is 0, b if selector is 1.
func AB01(a, b vec2, selector float) vec2 {
	return a*(1.0 - selector) + b*selector
}

// Returns a if selector is 1, b if selector is 0.
func AB10(a, b vec2, selector float) vec2 {
	return a*selector + b*(1.0 - selector)
}
```
I read these as "pick A or B with 0 or 1", and "pick A or B with 1 or 0".

### Area masks

Similar concept to [selectors](#selectors), but with basic geometric shapes instead of simple conditions. If the area is satisfied, 1 is returned, otherwise, we get 0 back from the function. Margins are included for smoother transitions between `0..1`. The more generalized form of this are signed distance fields (SDFs).

```Golang
// Given a margin of zero, if value is within [lo, hi],
// the function returns 1, otherwise, returns 0.
// Margins greater than zero smooth the result at the
// edges.
func inBandMask(value, lo, hi float, margin float) float {
	rangeLen := hi - lo
	presence := clamp(value - lo, 0, rangeLen)
	in  := smoothstep(0, margin, presence)
	out := smoothstep(0, margin, rangeLen - presence)
	return in*out
}

// Return 1 if the current position is within 'hardRadius' of 'target',
// between 1 and 0 if within 'hardRadius + softRadius' of 'target',
// or zero otherwise.
func inDotMask(current vec2, target vec2, hardRadius, softRadius float) float {
	return 1.0 - smoothstep(hardRadius, hardRadius + softRadius, distance(current, target))
}
```

## Sampling, coords and interpolation

### Relative target coordinates

Get `targetCoords` in `[0, 1]` range. We all forget the origin:
```Golang
relativeX := (targetCoords.x - imageDstOrigin().x)/imageDstSize().x // [0, 1]
relativeY := (targetCoords.y - imageDstOrigin().y)/imageDstSize().y // [0, 1]
```

If you are going to use both, better get it directly as a `vec2`.
```Golang
relTargetCoords := (targetCoords.xy - imageDstOrigin())/imageDstSize()
```

### Remap `sourceCoords` between images

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

### Texture clamping and repeat

If you use `imageSrcNAt()`, coordinates outside bounds will return `vec4(0)` by default. If you use `imageSrcNUnsafeAt()`, there's no clamping at all and you might access the whole underlying internal atlas. Here are some snippets for other behaviors:

> [!NOTE]
> *Notice that default Ebitengine shaders for `DrawImage()` and `DrawTriangles()` already support the [`Address`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Address) type, which allows changing between unsafe, clamp and repeat. The snippets below are the manual version of this for custom shaders.*

Extend edge color:
```Golang
//kage:unit pixels
package main

var Cursor vec2 // in [0, 1] ranges, built-in with kage-desk/display

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	// as an example, we offset sourceCoords based on the cursor
	// so we can go out of bounds and see the clamping in action
	modifiedCoords := sourceCoords + (Cursor - vec2(0.5))*200.0

	minSrc0, maxSrc0 := GetSource0ClampCoords()
	clampedCoords := clamp(modifiedCoords, minSrc0, maxSrc0)
	return imageSrc0UnsafeAt(clampedCoords)
}

func GetSource0ClampCoords() (vec2, vec2) {
	const epsilon = 1.0/16384.0 // TODO: how small can we safely set this?
	origin := imageSrc0Origin()
	return origin, origin + imageSrc0Size() - vec2(epsilon)
}
```

Wrap:
```Golang
//kage:unit pixels
package main

var Cursor vec2 // in [0, 1] ranges, built-in with kage-desk/display

func Fragment(_ vec4, sourceCoords vec2, _ vec4) vec4 {
	// as an example, we offset sourceCoords based on the cursor
	// so we can go out of bounds and see the clamping in action
	modifiedCoords := sourceCoords + (Cursor - vec2(0.5))*400.0

	originSrc0, endSrc0 := GetSource0Region()
	wrappedCoords := Wrap(modifiedCoords, originSrc0, endSrc0)
	return imageSrc0UnsafeAt(wrappedCoords)
}

func GetSource0Region() (vec2, vec2) {
	origin := imageSrc0Origin()
	return origin, origin + imageSrc0Size()
}

func Wrap(coords vec2, origin, end vec2) vec2 {
	return origin + mod(coords - origin, end - origin)
}
```

### Bilinear interpolation

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

## Misc.

### Math functions

```Golang
func fmod(value, modulo float) float {
	return value - modulo*trunc(value/modulo)
}

func trunc(a float) float {
	if a >= 0 { return floor(a) }
	return ceil(a)
}
```

### Colors

Using RGB [0-255] directly can be useful while testing colors. Just make sure to optimize afterwards if needed!
```Golang
func rgb(r, g, b int) vec4 {
	return vec4(float(r)/255.0, float(g)/255.0, float(b)/255.0, 1.0)
}
```
*(TODO: many useful functions that we might want to copy from black-and-white and HSL-hue-shift examples)*

### Aiding SDF visualization

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
