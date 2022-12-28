//go:build nope
package display

// TODO: add ESC to handling.

type AutoUniform uint8
const (
	UTime AutoUniform = iota // current time in seconds. Uniform name: "Time" (float)
	UTick // current tick. Uniform name: "Tick" (float)
	UCursor // cursor position. Uniform name: "Cursor" (vec2)
	UClick // last click position or (0, 0). Uniform name: "LastClickPos" (vec2)
	UBackColor // background color. Uniform name: "BackColor" (vec4)
	USliderUnit // custom slider from 0 to 1
	USliderCent // custom slider from 0 to 100
	USliderSigned // custom slider from -1 to 1
	USlider2Unit // TODO: sliders for ints would also be cool. maybe have actual slider uniform. or a couple.
	// TODO: colormap for interpolation. with N positions?
	USlider2Cent
	USlider2Signed
	UPlayerPos // simulates a player position with gamepads or keyboard wasd/arrows
)

type changingUniformID uint8
const (
	cuRand changingUniformID = iota
	cuAngleDeg
	cuAngleRad
)

type AutoChangingUniform struct {
	key changingUniformID
	hz float32
}

func checkHz(hz float64) {
	if hz < 0 {
		panic("AutoChangingUniform doesn't accept negative hz values")
	}
}

// Hz indicate how many times we re-roll the rand per second.
// For example, 0.1 would mean re-roll once every 10 seconds.
func URand(hz float64) AutoChangingUniform {
	checkHz(hz)
	return AutoChangingUniform{ cuAngleRad, float32(hz) }
}

// Hz indicate how many loops we do per second (from 0 to 360 degrees).
func UAngleDeg(hz float64) AutoChangingUniform {
	checkHz(hz)
	return AutoChangingUniform{ cuAngleDeg, float32(hz) }
}

// Hz indicate how many loops we do per second (from 0 to 360 degrees).
func UAngleRad(hz float64) AutoChangingUniform {
	checkHz(hz)
	return AutoChangingUniform{ cuAngleRad, float32(hz) }
}

// Predefined textures that can be used with your shader.
// You can pass up to 4 to the shader options, and they will
// be assigned from image 0 to image 3, but if all the images
// need to be resizable or have the same fixed size.
//
// Fixed size images have a number at the end of their name.
// Resizable images will be resized to the canvas size.
type AutoTexture uint8
const (
	TexTriangleRed AutoTexture = iota
	TexTriangleCyan
	TexTriangleColor
	TexNoiseMono
	TexNoiseColor
	// TODO: something that has normals for it too. Like, TeapotMask, TeapotNormals

	texFixedAnchor // unexported for internal use (resizable above, fixed size below)
	TexNoiseMono32x32
	TexEbiten32x32
	TexSprite32x32
	TexAniSprite32x32
)

