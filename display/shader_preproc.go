//go:build nope
package display

import "github.com/hajimehoshi/ebiten/v2"

type autoUniform interface {
	Update()
	Set(*ebiten.DrawTrianglesShaderOptions)
}

type autoUniformSet struct {
	tabPressed bool
	uniforms []autoUniform
}

func (self *autoUniformSet) DrawUI(screen *ebiten.Image) {

}

func (self *autoUniformSet) Update() {
	for i := 0; i < len(self.uniforms); i++ {
		self.uniforms[i].Update()
	}
}

func (self *autoUniformSet) Set(opts *ebiten.DrawTrianglesShaderOptions) {
	for i := 0; i < len(self.uniforms); i++ {
		self.uniforms[i].Set(opts)
	}
}

// Preprocess a program to extract an autoUniformSet. In the future we may
// want to modify the program itself, but for the moment we stick to
// comment-based macros only.
func preprocess(program []byte) *autoUniformSet {
	return nil
	//stage := "expect-package"
}
