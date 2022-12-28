package main

// In general you can use display.Shader() without arguments to
// search for a .kage shader in the working directory, or pass
// the filename directly, e.g.:
//    display.Shader("shader.kage")
//
// In these examples, though, we want the programs to be executable
// remotely from your terminal, so we are embedding the shader program
// data directly into the Golang program and exposing a variable.

import _ "embed"

//go:embed shader.kage
var shaderProgram []byte
