## Invoking shaders manually

While writing the first shaders we have been using `kage-desk/display` to keep our `main.go` file really simple. This is great to get started, but if we want to keep expanding our powers, we need to take off the training wheels. We have learned how to make shaders, but we still haven't seen how to invoke them from Ebitengine by ourselves.

There are two options:
- [Image.DrawRectShader(...)](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawRectShader).
- [Image.DrawTrianglesShader(...)](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawTrianglesShader).

Let's start by reworking our `main.go` with `DrawRectShader`, the simplest of the two methods:
```Golang
package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	// ...
}
```

WIP
