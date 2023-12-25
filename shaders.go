package main

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type getGlParam func(uint32, uint32, *int32)
type getInfoLog func(uint32, int32, *int32, *uint8)

func checkGlError(glObject uint32, errorParam uint32, getParamFn getGlParam,
	getInfoLogFn getInfoLog, failMsg string) {

	var success int32
	getParamFn(glObject, errorParam, &success)
	if success != 1 {
		var infoLog [512]byte
		getInfoLogFn(glObject, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		log.Fatalln(failMsg, "\n", string(infoLog[:512]))
	}
}

func checkShaderCompileErrors(shader uint32) {
	checkGlError(shader, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"ERROR::SHADER::COMPILE_FAILURE")
}

func checkProgramLinkErrors(program uint32) {
	checkGlError(program, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"ERROR::PROGRAM::LINKING_FAILURE")
}

type Shader struct {
	program uint32
}

type ShaderType uint32

const (
	VERTEX_SHADER   ShaderType = gl.VERTEX_SHADER
	FRAGMENT_SHADER ShaderType = gl.FRAGMENT_SHADER
)

// Compile the provided shader source and return the shader object.
func CompileShader(source string, shader_type ShaderType) Shader {
	program := gl.CreateShader(uint32(shader_type))
	source_chars, free_func := gl.Strs(source)
	defer free_func()
	gl.ShaderSource(program, 1, source_chars, nil)
	gl.CompileShader(program)
	checkShaderCompileErrors(program)
	return Shader{program}
}

type ShaderProgram struct {
	program uint32
}

// Link the provided shaders in the order they were given and return the linked program.
// The shader objects are not needed after they are linked into a program object, and they
// should be deleted.
func LinkShaders(shaders []Shader) ShaderProgram {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader.program)
	}
	gl.LinkProgram(program)
	checkProgramLinkErrors(program)

	// shader objects are not needed after they are linked into a program object
	for _, shader := range shaders {
		gl.DeleteShader(shader.program)
	}

	return ShaderProgram{program}
}

func (sp ShaderProgram) Use() {
	gl.UseProgram(sp.program)
}
