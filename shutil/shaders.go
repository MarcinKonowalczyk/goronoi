package shutil

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type getGlParam func(uint32, uint32, *int32)
type getInfoLog func(uint32, int32, *int32, *uint8)

func checkGlError(
	glObject uint32,
	errorParam uint32,
	getParamFn getGlParam,
	getInfoLogFn getInfoLog) string {

	var success int32
	getParamFn(glObject, errorParam, &success)
	if success != 1 {
		var infoLog [512]byte
		getInfoLogFn(glObject, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		message := string(infoLog[:512])
		return message
	}
	return ""
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
	source_chars, free_func := gl.Strs(source + "\x00")
	defer free_func()
	gl.ShaderSource(program, 1, source_chars, nil)
	gl.CompileShader(program)
	message := checkGlError(program, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog)

	if message != "" {
		log.Fatalln("ERROR::SHADER::COMPILATION_FAILED. Source:\n", source, "\n", message)
	}

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
	message := checkGlError(program, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog)
	if message != "" {
		log.Fatalln("ERROR::PROGRAM::LINKING_FAILURE\n", message)
	}

	// shader objects are not needed after they are linked into a program object
	for _, shader := range shaders {
		gl.DeleteShader(shader.program)
	}

	return ShaderProgram{program}
}

func (sp ShaderProgram) Use() {
	gl.UseProgram(sp.program)
}

func (sp ShaderProgram) Unuse() {
	gl.UseProgram(0)
}

func (sp ShaderProgram) Delete() {
	gl.DeleteProgram(sp.program)
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

func (sp ShaderProgram) GetAttribLocation(name string) uint32 {
	sp.Use()
	location := gl.GetAttribLocation(sp.program, gl.Str(name+"\x00"))
	if location == -1 {
		log.Fatalln("Invalid attribute name", name)
	}
	return uint32(location)
}

func (sp ShaderProgram) GetUniformLocation(name string) uint32 {
	sp.Use()
	location := gl.GetUniformLocation(sp.program, gl.Str(name+"\x00"))
	if location == -1 {
		log.Fatalln("Invalid uniform name", name)
	}
	return uint32(location)
}

func (sp ShaderProgram) SetUniform1f(name string, x float32) {
	location := int32(sp.GetUniformLocation(name))
	gl.Uniform1f(location, x)

	// read back the uniform value and check it
	var value float32
	gl.GetUniformfv(sp.program, location, &value)

	if value != x {
		log.Fatalln("Uniform value was not set correctly")
	}
}

func (sp ShaderProgram) SetUniform2f(name string, vec [2]float32) {
	location := int32(sp.GetUniformLocation(name))
	gl.Uniform2f(location, vec[0], vec[1])

	// read back the uniform value and check it
	var value [2]float32
	gl.GetUniformfv(sp.program, location, &value[0])

	if value[0] != vec[0] || value[1] != vec[1] {
		log.Fatalf("Uniform value was not set correctly: %v != %v", value, vec)
	}
}

func (sp ShaderProgram) SetUniform3f(name string, vec [3]float32) {
	location := int32(sp.GetUniformLocation(name))
	gl.Uniform3f(location, vec[0], vec[1], vec[2])

	// read back the uniform value and check it
	var value [3]float32
	gl.GetUniformfv(sp.program, location, &value[0])

	if value[0] != vec[0] || value[1] != vec[1] || value[2] != vec[2] {
		log.Fatalln("Uniform value was not set correctly")
	}
}

func (sp ShaderProgram) SetUniform4f(name string, vec [4]float32) {
	location := int32(sp.GetUniformLocation(name))
	gl.Uniform4f(location, vec[0], vec[1], vec[2], vec[3])

	// read back the uniform value and check it
	var value [4]float32
	gl.GetUniformfv(sp.program, location, &value[0])

	if value[0] != vec[0] || value[1] != vec[1] || value[2] != vec[2] || value[3] != vec[3] {
		log.Fatalln("Uniform value was not set correctly")
	}
}

func (sp ShaderProgram) SetUniform1i(name string, x int32) {
	location := int32(sp.GetUniformLocation(name))
	gl.Uniform1i(location, x)

	// read back the uniform value and check it
	var value int32
	gl.GetUniformiv(sp.program, location, &value)

	if value != x {
		log.Fatalln("Uniform value was not set correctly")
	}
}
