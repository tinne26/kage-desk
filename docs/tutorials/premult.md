# Premultiplied alpha in Ebitengine

So you are happily writing code, creating some nice visuals to share with your imaginary friends...

![](https://github.com/tinne26/kage-desk/blob/main/img/premult_opaque.png?raw=true)

"Yeah, it's looking stylish!", you say to yourself.

"But perhaps we can make a couple slight changes here and there... *yo, what's going on with color?!*"

![](https://github.com/tinne26/kage-desk/blob/main/img/premult_transparent.png?raw=true)

And suddenly color gets weird.

No, you are not going insane. No, Ebitengine isn't broken. It's our good friend the **premultiplied alpha**.

*(and this one is sadly very real)*

## Huh..?

Most of the time we are working with RGBA values. One color, four channels: red, green, blue and finally alpha, which is kinda special.

When it comes to image processing (rendering, composition, and so on), there are multiple ways to interpret the alpha value of a given color. Computer graphics benefit most from the *premultiplied alpha* format, which considers that the alpha value of a color *has already been applied to the other color channels* in a pre-multiplication process.

In other words, if we take the non-premultiplied color `clr := vec4(1.0, 1.0, 1.0, 0.5)` and we want to premultiply it, we would do `clr = vec4(clr.RGB*clr.A, clr.A)`, which would give us `vec4(0.5, 0.5, 0.5, 0.5)`.

This might not mean much to you yet, but don't worry, just keep going. Wihle the concept may be tricky to grasp, its shape is actually quite easy to detect:
- If you use a `color.RGBA` or a `vec4` where the alpha value is lower than any of the other channels (e.g., `color.RGBA{255, 255, 0, 128}`, `vec4(0, 0, 128, 64)`), you are creating an *invalid premultiplied alpha color* and bad things will happen.

If you take a look at the documentation in [image/color](https://pkg.go.dev/image/color), you will see that it's full of references to "premultiplied alpha" and "non-premultiplied alpha". When Ebitengine asks you for `color.Color` values in the API, you have to pay attention to the actual types you are using. If you are using the most common color type, `color.RGBA`, you have to obey the spec and *use premultiplied values when you create the color*. Other types, like `color.NRGBA`, can be easier to use at the beginning. If you are using shaders, though, Ebitengine always requires you to use the premultiplied alpha format.

For the moment, just keep in mind that you don't want to be using `color.RGBA` or `vec4` color values where the alpha channel value is lower than any of the other three color channels. You can come back later to review this section and it will eventually click.

## How to solve the issues

There are three main ways to navigate the troubles that premultiplied alpha brings to your life:
1. Understanding the enemy and adjusting your maths to avoid creating invalid premultiplied alpha colors in the first place.
2. Same as 1, but guarding yourself actively with code. For example, instead of creating colors like `color.RGBA{r, g, b, a}`, use a function like `NewRGBA(r, g, b, a uint8) color.RGBA` that panics if `r`, `g` or `b` are greater than `a`.
3. Use `color.NRGBA` instead of `color.RGBA`. The format of `color.NRGBA` is non-premultiplied alpha, so any value you can pass to it during creation will be valid.

In general, the third solution is the easiest one to apply; you don't need to understand anything and can kinda forget that premultiplied alpha even exists. That being said, I think understanding the basic maths behind the conversions is beneficial: 
- You can work with shaders —where premultiplied alpha is unavoidable— without trouble.
- You become aware of the overhead that `color.NRGBA` can have (*admittedly small and either irrelevant or unavoidable in most cases anyway*).

## Case study

Remember the first image in this page? Well, let's say we want to make the circular shapes *slowly fade in* in our Ebitengine game. We could try something like this:

```Golang
package main

import ( "log" ; "image/color" )
import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/vector"

type Game struct { fadeIn uint8 }
func (self *Game) Layout(w, h int) (int, int) { return w, h }
func (self *Game) Update() error {
	if self.fadeIn < 255 { self.fadeIn += 1 }
	return nil
}

func (self *Game) Draw(canvas *ebiten.Image) {
	canvas.Fill(color.RGBA{246, 197, 175, 255})
	bounds := canvas.Bounds()
	cx := float32(bounds.Min.X + bounds.Dx()/2)
	cy := float32(bounds.Min.Y + bounds.Dy()/2)
	unit := float32(bounds.Dy())/100.0
	aa := false // no antialiasing
	
	circleColor := color.RGBA{218, 65, 103, self.fadeIn}
	vector.DrawFilledCircle(canvas, cx, cy, unit*18.0, circleColor, aa)
	vector.StrokeCircle(canvas, cx, cy, unit*25.0, unit*1.5, circleColor, aa)
	vector.StrokeCircle(canvas, cx, cy, unit*33.3, unit*0.8, circleColor, aa)
}

func main() {
	ebiten.SetWindowSize(360, 360)
	err := ebiten.RunGame(&Game{})
	if err != nil { log.Fatal(err) }
}
```

Can you spot the mistake? Yeah, of course. This is what we would expect:



This is what we actually get:



At the beginning of the animation process, since `Game.fadeIn` starts increasing from 0, the color is *not a valid premultiplied alpha color*. The color only becomes valid when alpha reaches a value of 218. The solution in this program is, of course, replacing `color.RGBA` with `color.NRGBA`.


## Why is pre-multiplied alpha needed at all?

Color composition can be done more efficiently when alpha has been premultiplied. Basically, the programs don't have to perform additional products and conversions, and they can ignore the alpha value in most scenarios. You don't really need to know any of this, but it's always reassuring to know there's a justified reason for your pain.

You will likely make mistakes with premultiplied alpha from time to time, but as long as you can quickly remember what's going on, it will always be trivial to fix. Just keep this in mind: if something weird is going on with color, as the first step, always ask yourself if premultiplied alpha could be your culprit.


