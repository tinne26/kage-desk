# Ebitengine's `Game` interface

Ebitengine games look very simple to implement. Just create a struct that implements the following methods:
- `Draw()`: draw the current game state.
- `Update()`: update the game state.
- `Layout()`: specify the game canvas size.
And then pass an instance of that struct to `ebiten.RunGame(myGame)`.

Easy, right?

```Golang
type Game struct {}
func (self *Game) Layout(width, height int) (int, int) {
	return 512, 512
}
func (self *Game) Update() error {
	return nil
}
func (self *Game) Draw(screen *ebiten.Image) {
	// ...
}
```

Not so fast!

The truth is that while *the general idea* sounds very accessible, when you start implementing a complex game... you will start asking yourself some questions:
- In which order are the functions called? Can `Update()` and `Draw()` be called at the same time?
- Can I update logic inside `Draw()`? Should I check any input there?
- Should I use elapsed times inside `Update()`?
- Can I control the frame rate? What's the difference between FPS and TPS?
- What's layout *really* doing? Is it *Ebitengine telling me* the screen size, or is it *me telling Ebitengine* the screen size?
- What's the "game size"? If my game has a resolution of 128x72, how do I scale that up depending on the display's resolution? What will happen on fullscreen?
- How do I support HiDPI screens and display scaling?

## The guidelines

We will explain everything in more detail over the next sections, but first let's give some general guidelines that will keep you on the right path... most of the time:
- `Update()`, `Draw()` and `Layout()` can be called in any order, including multiple calls to one before the other, but they are never called concurrently. You should *not* be worrying about any of this.
- Update your logic and check your inputs in `Update()`. Never use elapsed time. The TPS are the ticks per second or updates per second, which control the rate at which the logic of your game is evaluated. The default TPS of 60 should be good for most games, so don't touch that either unless you have a really good reason to.
- `Draw()` should only draw your current game state and not evaluate or update any logic. Don't worry about the refresh rate.
- There are two main ways to use `Layout()`:
	- For simple games that *only* use pixel art, you can use a fixed size layout. Just return a fixed size from `Layout()` and let Ebitengine worry about it.
	- For any games that use high resolution assets, scalable text or anything like that, use `LayoutF()`, `DeviceScaleFactor()` and adjust draw sizes based on your current "canvas" size. There's a snippet for this on the [high quality assets and vectorial art](#high-quality-assets-and-vectorial-art) section.

## Update and Draw

Consider reading [tinne26/tps-vs-fps](https://github.com/tinne26/tps-vs-fps) if you need to understand them in detail.

## Layout

The layout is one of the most important parts to get right if you don't want your game to be poorly scaled on different displays and you want your graphics to look sharp. Sadly, layouts are a tricky topic and you really need to understand what you are doing.

First, let's see the method signature:
```Golang
Layout(logicWinWidth, logicWinHeight int) (canvasWidth, canvasHeight int)
```

Ebitengine calls `Layout()` with the current window size in logical or *device-independent pixels*, which are pixels divided by the display scaling[^1]. Then, based on this "window size", you can tell Ebitengine what do you want your "canvas size" to be. The canvas size is the specific size in pixels of the "canvas" you want to draw your game graphics to. If the canvas size is different from the window size with the display scale applied, then Ebitengine will automatically apply some scaling to your "canvas" so it fits the current window.

You are probably still super confused. Logical pixels? Display scaling? Window size? Canvas size? Are there "real pixels" too then? Weeeeell... on some platforms you can have "real pixels" and on others not (e.g. macOS). Sometimes you simply don't know exactly how your pixels will be projected to the screen. We are still going to try our best.

### One step at a time

Let's start with a low-resolution game. Imagine you make a game with a very small 128x72 resolution (128 pixels wide, 72 pixels tall) and want it to look as good as possible on Ebitengine. When we say 128x72, we are referring to our "canvas size". We will make pixel art assets that fit within a game screen of that size. Canvas sizes are not some abstract dimensions: your sprites have a concrete amount of pixels, and you position these sprites in concrete positions of the canvas. Artists can set each pixel with love, care and certainty. Simple stuff, great stuff.

So we make our beautiful pixel art, our fantastic pixelated world... and when everything seems perfect, the time to project it to a screen of arbitrary size[^2] arrives.

The simplest option in this case would be the following:
```Golang
func (_ *Game) Layout(_, _ int) (int, int) { return 128, 72 }
```
We tell Ebitengine that we *don't care* what the current screen size is, that we want a canvas size of 128x72 and we want to forget about the rest. Ebitengine will then do its best to scale this canvas to the final window size.

This simple approach is actually perfectly ok for most pixel art games that *only contain pixel art*. Ebitengine won't change the aspect ratio of your game and will keep black bars on the top or the sides of the screen if necessary, but beyond that it will only try to scale the canvas to fill as much of the window as possible.

You may notice that this can still produce distortions. For example, if the actual screen size is 1920x1080, the canvas will be zoomed-in x15 in both dimensions. Since this is a whole number, what before was one pixel now will become a block of 15x15 pixels, but the game will visually remain the same. But if the screen size is 1366x768 instead, the scaling will be x10.67! Scaling the canvas as much as possible can cause distortions.

On the topic of distortions and scaling, you should know that Ebitengine supports both simple interpolation and nearest neighbour scaling, which can be configured through `ebiten.SetScreenFilterEnabled(bool)`. This filter is enabled by default, but in some very specific cases the game may look better with the filter disabled. Just keep it in mind as a possibility.

Sadly, this will only change the "type of distortions" that you get. Depending on your game's original resolution and the screen it's being projected to, these distortions can be more or less noticeable.

One idea I personally like is adding an option in the game to do only integer scaling. Instead of scaling by x10.67, one can do the scaling manually and truncate the scaling factor so it's x10 instead, and then center the result on the screen. Depending on the game resolution, this can end up wasting a lot of space on the screen, but it's an idea to keep in mind if you want to go the extra mile to allow a truly pixel-perfect experience. Unless the game is being played on macOS, of course. You can't do pixel-perfect on macOS, period.

### High quality assets and vectorial art

In the previous section we discussed the many problems that we have when projecting our game canvas to the final screen. As we have seen, when we are working with small resolutions and pure pixel art assets, we can return a fixed layout size and kinda forget about the rest.

Unfortunately for you, pure pixel art games are *very uncommon*. In almost all cases, you will also want high resolution text, camera effects, slow-scrolling-but-smooth parallaxes or many other visual effects that aren't really pixel-art in nature. Or you may be using high-resolution assets from the start.

In any of these cases we will need to make use of the full screen resolution:
```Golang
func (_ *Game) Layout(logicWinWidth, logicWinHeight int) (int, int) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth  := int(math.Ceil(float64(logicWinWidth )*scale))
	canvasHeight := int(math.Ceil(float64(logicWinHeight)*scale))
	return canvasWidth, canvasHeight
}
```

You will notice that the canvas size will be bigger than the given window size if the display scaling factor is greater than one. This is expected because the given window size is *logical*, which means that it has been pre-divided by the display scaling factor.

This is still not enough to make use of the full screen resolution, though, as logical window sizes can be fractional. Using the `int` values from `Layout()` can still result in a lossy and blurry projection at certain sizes if window resizing is allowed... which is why Ebitengine v2.5.0 added `LayoutF()` supporting `float64` values. If your game struct implements it, Ebitengine will call `LayoutF()` instead of `Layout()`. The correct conversion is the following:
```Golang
func (_ *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth  := math.Ceil(logicWinWidth*scale)
	canvasHeight := math.Ceil(logicWinHeight*scale)
	return canvasWidth, canvasHeight
}
```

Notice that `Layout()` still needs to exist with a dummy implementation (e.g. `panic("unused")`) so your struct complies with the `ebiten.Game` interface. Can't wait for `LayoutF()` to become the default in v3!

The need for `math.Ceil` is an implementation detail that we shouldn't have to be aware of, but that's what we have for the moment.

With all this, the "canvas" that we receive in `Draw()` will have the maximum size possible, meaning we can draw at the highest resolution supported by the screen... unless the OS gets in the way somehow (*ehem ehem macOS*).

Now, drawing your assets with the proper position and scale is kind of your problem. It's all basic maths, but it can get a bit tricky. If you are working on an almost-fully-pixel-art game, instead of trying to scale everything as you draw it, you may also consider keeping a fixed size offscreen canvas to draw the pixel art to, projecting this offscreen to the main canvas with the proper scaling, and finally do a second pass where you draw any remaining high resolution or scalable assets, like text. Of course, this only works if the high resolution or vectorial graphical assets are always on top of the pixel-art canvas. Otherwise you will start having too many layers and it can get very messy too, but now you know.

### Summary

Layout gives us the **window size** and we return our desired **canvas size** based on it. For pure pixel-art games, we may return a fixed size and let Ebitengine do the scaling internally. Otherwise, we will need to manage this ourselves and **make use of the full screen resolution**.

For high resolution games, we will need to make use of **`LayoutF()`**, **`DeviceScaleFactor()`** and keep this scaling factor into account for many game elements.

[^1]: The display scaling is the "zoom level" applied to the screen contents. Screens can have different resolutions, and the pixel sizes can also vary between screens. For example, a screen with a very high resolution may have pixels that are half the size compared to another screen. In practice, this means we should not be drawing fixed size graphics (e.g. 32x32 pixels) on a screen. We need to worry about the display scaling and apply it in order to properly dimension the content we want to draw. In Ebitengine, the display scaling can be obtained through the `DeviceScaleFactor()` function.
[^2]: To be fair, screen sizes are not arbitrary. The most common aspect ratio is 16:9. When making a pure pixel-art game, you should choose a multiple of that (e.g. 128x72, 256x144, 512x288, 768x432...). The most common screen resolution is 1920x1080, so choosing a size that divides that evenly is almost always a good idea. You may also go for a 4:3 ratio, also quite common in games, but then you will have black borders on most modern screens (you could do screen stretching yourself, but that's *really ugly*).
