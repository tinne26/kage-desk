# Common pitfalls

You will get no warning for integer divisions in float contexts:
```Golang
if color.r >= 2/3 {
	return vec4(1, 0, 0, 1)
}
```

You will forget to account for the origin of target and sources:
```Golang
// Assumes that Cursor.x is in [0, 1] range.
if Cursor.x > (targetCoords.x - imageDstOrigin().x)/imageDstSize().x {
	return vec4(1, 1, 1, 1)
} else {
	return vec4(0, 0, 0, 1)
}
```
