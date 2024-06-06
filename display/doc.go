// A package that makes it easy to execute shaders with very little setup.
// Here's a basic example:
//   package main
//   
//   import "github.com/tinne26/kage-desk/display"
//   func main() {
//       display.SetTitle("basic-shader")
//       display.SetSize(512, 512)
//       display.Shader() // will look for a .kage file in the local directory
//   }
//
// Here's a slightly more advanced example:
//   package main
//   
//   import _ "embed"
//   import "github.com/hajimehoshi/ebiten/v2"
//   import "github.com/tinne26/kage-desk/display"
//
//   //go:embed shader.kage
//   var shader []byte
// 
//   func main() {
//       display.SetTitle("nice-shader")
//       display.SetSize(512, 512, display.Resizable, display.HiRes)
//       display.SetBackColor(display.BCJade) // set a background color
//       display.LinkUniformKey("Mode", 1, ebiten.KeyDigit1, ebiten.KeyNumpad1)
//       display.LinkUniformKey("Mode", 2, ebiten.KeyDigit2, ebiten.KeyNumpad2)
//       display.SetUniformInfo("Mode", "%d")
//       display.Shader(shader) // shader program passed explicitly
//   }
// 
// In this example we have embedded the shader explicitly (great for WASM and single
// executable builds), made the screen resizable, made the game layout take device
// scaling into account, defined keys 1 and 2 to change a uniform named "Mode" and
// configured the program to show the value of that new uniform in screen with
// [SetUniformInfo]().
//
// You also get some additional uniforms, images, vertex colors and shortcuts for free;
// please refer to [Shader]() for more details on all that.
//
// Finally, the package also detects some flags like '--maxfps' (unlimit fps and
// display them on the title), '--fullscreen' and '--opengl' (Windows would use
// DirectX by default otherwise).
package display
