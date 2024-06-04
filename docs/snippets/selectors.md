# Selectors

In shaders, it's often useful to compose the results of multiple operations based on some conditions. This can actually be done directly in conditionals most of the time, but some programmers prefer to avoid these in their shaders[^1].

[^1]: The reason for programmers wanting to avoid conditionals on shaders is quite long-winded and complex; in the past, they were kinda suboptimal in many cases. Nowadays, most of them are not a problem, but since the same program is executed for a wavefront (groups of pixels), whenever any of the pixels takes a different conditional branch, this branch must still be evaluated for all the pixels in the wavefront. A short `if` is *very different* from a set of completely distinct, long branches, but everyone is allowed to have their own preferences.

I'm not going to advocate against using conditionals in shaders or anything, but here's a set of functions that can be useful to keep/discard results based on simple conditions:
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
