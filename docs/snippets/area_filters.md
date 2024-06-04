# Area filters

Similar concept to the [selectors](https://github.com/tinne26/kage-desk/blob/main/docs/snippets/selectors.md) snippets, but with basic geometric shapes instead of simple conditions. If the area is satisfied, 1 is returned, otherwise, we get 0 back from the function. Margins are included for smoother transitions between `0..1`. The more generalized form of this are signed distance fields (SDFs).

```Golang
// Given a margin of zero, if value is within [lo, hi],
// the function returns 1, otherwise, returns 0.
// Margins greater than zero smooth the result at the
// edges.
func band(value, lo, hi float, margin float) float {
	rangeLen := hi - lo
	presence := clamp(value - lo, 0, rangeLen)
	in  := smoothstep(0, margin, presence)
	out := smoothstep(0, margin, rangeLen - presence)
	return in*out
}

// Return 1 if the current position is within 'hardRadius' of 'target',
// between 1 and 0 if within 'hardRadius + softRadius' of 'target',
// or zero otherwise.
func dot(current vec2, target vec2, hardRadius, softRadius float) float {
	return 1.0 - smoothstep(hardRadius, hardRadius + softRadius, distance(current, target))
}
```
