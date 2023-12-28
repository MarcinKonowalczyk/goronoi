package shutil

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

type Uniform struct {
	name  string
	utype uint32
	size  int32
}

func (sp ShaderProgram) GetActiveUniforms() []Uniform {
	num_uniforms := int32(0)
	gl.GetProgramiv(sp.program, gl.ACTIVE_UNIFORMS, &num_uniforms)

	if num_uniforms == 0 {
		return []Uniform{}
	}

	uniforms := make([]Uniform, num_uniforms)

	for i := int32(0); i < num_uniforms; i++ {
		var name_len int32
		var size int32
		var gl_type uint32
		name_null := make([]uint8, 256)
		gl.GetActiveUniform(sp.program, uint32(i), 256, &name_len, &size, &gl_type, &name_null[0])
		name := string(name_null[:name_len])
		uniforms[i] = Uniform{name, gl_type, size}
	}

	return uniforms
}

func (sp ShaderProgram) SetUniform1f(name string, x float32) {
	location := gl.GetUniformLocation(sp.program, gl.Str(name+"\x00"))
	if location == -1 {
		log.Fatalln("Invalid uniform name", name)
	}
	gl.Uniform1f(location, x)

	// read back the uniform value and check it
	var value float32
	gl.GetUniformfv(sp.program, location, &value)

	if value != x {
		log.Fatalln("Uniform value was not set correctly")
	}
}

func (sp ShaderProgram) SetUniform2f(name string, x, y float32) {
	location := gl.GetUniformLocation(sp.program, gl.Str(name+"\x00"))
	if location == -1 {
		log.Fatalln("Invalid uniform name", name)
	}
	gl.Uniform2f(location, x, y)

	// read back the uniform value and check it
	var value [2]float32
	gl.GetUniformfv(sp.program, location, &value[0])

	if value[0] != x || value[1] != y {
		log.Fatalln("Uniform value was not set correctly")
	}
}
