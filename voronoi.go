package main

import (
	"fmt"
	"image"
	"log"
	"runtime"
	"unsafe"
	"voronoi/shutil"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"

	_ "embed"
)

// const fontfile = "luxisr.ttf"

const windowWidth = 800
const windowHeight = 600

//go:embed quad.vert
var vertexShaderSource string

//go:embed quad.frag
var fragmentShaderSource string

//go:embed test.vert
var testVertexShaderSource string

//go:embed test.frag
var testFragmentShaderSource string

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Hello!", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)

	// Initialize gl in the shutil package
	shutil.Init()

	window.SetKeyCallback(keyCallback)

	// Set window size callback
	window.SetSizeCallback(windowSizeCallback)

	// Cap the framerate at 60fps
	glfw.SwapInterval(1)

	// Load the font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	// // Create the text renderer
	// // make sure the NewTextRenderer is called in the same thread as the OpenGL context
	// text_renderer := shutil.NewTextRenderer(font, 36)

	// // // Render 'hello world' to the top left of the screen
	// text_renderer.RenderText(
	// 	"Quick brown fox jumps over the lazy dog",
	// 	[2]float32{-0.98, 0},
	// 	[2]float32{float32(windowWidth), float32(windowHeight)},
	// )

	testShaders := []shutil.Shader{
		shutil.CompileShader(testVertexShaderSource, gl.VERTEX_SHADER),
		shutil.CompileShader(testFragmentShaderSource, gl.FRAGMENT_SHADER),
	}

	testShaderProgram := shutil.LinkShaders(testShaders)
	defer testShaderProgram.Delete()

	testShaderProgram.Use()

	// Make a face
	face := truetype.NewFace(font, &truetype.Options{
		Size: 1,
		DPI:  10,
	})
	// Make a glyph
	dot := fixed.P(0, 0)
	dr, mask, _, _, ok := face.Glyph(dot, 'A')
	if !ok {
		panic("failed to get glyph")
	}
	// Create a texture for the glyph
	var texture uint32
	gl.GenTextures(1, &texture)
	fmt.Println("texture", texture)

	// Set the texture as active
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// Set the texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE) // set texture wrapping to GL_REPEAT (default wrapping method)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST) // set texture filtering (options: GL_LINEAR, GL_NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// Print glyph mask
	fmt.Println("mask", mask)

	// Upload the glyph to the texture
	pixels_unsafe_pointer := unsafe.Pointer(&mask.(*image.Alpha).Pix[0])
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(dr.Dx()), int32(dr.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, pixels_unsafe_pointer)

	// order: bottom left, top left, bottom right, top right
	vertices := []float32{
		-0.9, -0.9, 0.0, 0.0,
		-0.9, 0.9, 0.0, 1.0,
		0.9, -0.9, 1.0, 0.0,
		0.9, 0.9, 1.0, 1.0,
	}

	// Create a VAO
	vao := shutil.CreateVAO(vertices, 4)

	// Bind the VAO
	vao.Bind()

	gl.ClearColor(0.137, 0.137, 0.137, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	window.SwapBuffers() // Swap the rendered buffer with the window
	dummyLoop(window)
	// programLoop(window, text_renderer)
}

func compileShaders() []shutil.Shader {
	if vertexShaderSource == "" || fragmentShaderSource == "" {
		panic("vertexShaderSource or fragmentShaderSource is empty")
	}
	vertexShader := shutil.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := shutil.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	return []shutil.Shader{vertexShader, fragmentShader}
}

func programLoop(window *glfw.Window, text_renderer shutil.TextRenderer) {
	// the linked shader program determines how the data will be rendered
	shaders := compileShaders()
	shaderProgram := shutil.LinkShaders(shaders)
	defer shaderProgram.Delete()

	quad_vertices := []float32{
		0.9, 0.9, 0.0,
		0.9, -0.9, 0.0,
		-0.9, 0.9, 0.0,
		-0.9, -0.9, 0.0,
	}

	quad := shutil.CreateVAO(quad_vertices, 3)

	// We don't need to bind anything here because we only have one VAO
	quad.Bind()

	// Scale the resolution to the monitor's content scale
	// This is necessary for retina displays
	monitor := glfw.GetPrimaryMonitor()
	scale_x, scale_y := monitor.GetContentScale()

	shaderProgram.Use()

	// Augment the windowSizeCallback to update the resolution uniform
	newWindowSizeCallback := func(window *glfw.Window, width int, height int) {
		windowSizeCallback(window, width, height)
		scale_x, scale_y := window.GetContentScale()
		f32_width := float32(float32(width) * scale_x)
		f32_height := float32(float32(height) * scale_y)
		shaderProgram.SetUniform2f("u_resolution", [2]float32{f32_width, f32_height})
	}
	window.SetSizeCallback(newWindowSizeCallback)
	newWindowSizeCallback(window, windowWidth, windowHeight)

	setMouseUniform(window, shaderProgram, scale_x, scale_y)
	setTimeUniform(shaderProgram)

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// Get current mouse position
		setMouseUniform(window, shaderProgram, scale_x, scale_y)
		setTimeUniform(shaderProgram)

		// Set the color to clear the screen with
		gl.ClearColor(0.137, 0.137, 0.137, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// NOTE: We're not calling quad.Bund and Unbind here because we only have one VAO
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		// Swap in the rendered buffer
		window.SwapBuffers()
	}
}

// Dummy loop that just polls events and does nothing else. Useful for testing.
func dummyLoop(window *glfw.Window) {
	for !window.ShouldClose() {
		glfw.PollEvents()
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func windowSizeCallback(window *glfw.Window, width int, height int) {
	scale_x, scale_y := window.GetContentScale()
	gl.Viewport(0, 0, int32(width)*int32(scale_x), int32(height)*int32(scale_y))
}

// Set the mouse coordinates uniform. We assume that the shader program is already in use.
func setMouseUniform(
	window *glfw.Window,
	shaderProgram shutil.ShaderProgram,
	scale_x float32,
	scale_y float32,
) {
	mouse_x_f64, mouse_y_f64 := window.GetCursorPos()
	mouse_x := float32(mouse_x_f64 * float64(scale_x))
	mouse_y := float32(mouse_y_f64 * float64(scale_y))
	shaderProgram.SetUniform2f("u_mouse", [2]float32{mouse_x, mouse_y})
}

func setTimeUniform(shaderProgram shutil.ShaderProgram) {
	time := glfw.GetTime()
	shaderProgram.SetUniform1f("u_time", float32(time))
}
