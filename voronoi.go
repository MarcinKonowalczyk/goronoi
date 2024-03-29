package main

import (
	"log"
	"runtime"
	"voronoi/glu"
	"voronoi/glu/font"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/font/gofont/gomono"

	_ "embed"
)

const windowWidth = 800
const windowHeight = 600

//go:embed quad.vert
var vertexShaderSource string

//go:embed quad.frag
var fragmentShaderSource string

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

	// Initialize gl in the glu package
	glu.Init()

	window.SetKeyCallback(keyCallback)

	// Set window size callback
	window.SetSizeCallback(windowSizeCallback)

	// Cap the framerate at 60fps
	glfw.SwapInterval(1)

	// Load the font

	font, err := font.LoadFontBytes(gomono.TTF, int32(32), windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	font.UpdateResolution(windowWidth, windowHeight)
	oldWindowSizeCallback := window.SetSizeCallback(nil)
	newSizeCallback := func(window *glfw.Window, width int, height int) {
		// fmt.Println("New window size:", width, height)
		oldWindowSizeCallback(window, width, height)
		font.UpdateResolution(width, height)
	}
	window.SetSizeCallback(newSizeCallback)

	programLoop(window, font)
}

func compileShaders() []glu.Shader {
	if vertexShaderSource == "" || fragmentShaderSource == "" {
		panic("vertexShaderSource or fragmentShaderSource is empty")
	}
	vertexShader := glu.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := glu.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	return []glu.Shader{vertexShader, fragmentShader}
}

func programLoop(window *glfw.Window, font *font.Font) {
	// the linked shader program determines how the data will be rendered
	shaders := compileShaders()
	shaderProgram := glu.LinkShaders(shaders)
	defer shaderProgram.Delete()

	quad_vertices := []float32{
		0.9, 0.9, 0.0,
		0.9, -0.9, 0.0,
		-0.9, 0.9, 0.0,
		-0.9, -0.9, 0.0,
	}

	quad := glu.CreateVertexArray(quad_vertices, 3)

	// We don't need to bind anything here because we only have one VAO
	quad.Bind()

	// Scale the resolution to the monitor's content scale
	// This is necessary for retina displays
	monitor := glfw.GetPrimaryMonitor()
	scale_x, scale_y := monitor.GetContentScale()

	shaderProgram.Use()

	// Augment the windowSizeCallback to update the resolution uniform
	oldWindowSizeCallback := window.SetSizeCallback(nil)
	newWindowSizeCallback := func(window *glfw.Window, width int, height int) {
		oldWindowSizeCallback(window, width, height)
		scale_x, scale_y := window.GetContentScale()
		f32_width := float32(float32(width) * scale_x)
		f32_height := float32(float32(height) * scale_y)
		shaderProgram.SetUniform2f("u_resolution", [2]float32{f32_width, f32_height})
	}
	window.SetSizeCallback(newWindowSizeCallback)
	newWindowSizeCallback(window, windowWidth, windowHeight)

	setMouseUniform(window, shaderProgram, scale_x, scale_y)

	frame := uint32(0)
	setTimeUniform(shaderProgram, frame)

	font.SetColor(1.0, 1.0, 1.0, 0.8)

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()
		shaderProgram.Use()
		quad.Bind()

		// Set the color to clear the screen with
		gl.ClearColor(0.137, 0.137, 0.137, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Get current mouse position
		setMouseUniform(window, shaderProgram, scale_x, scale_y)
		setTimeUniform(shaderProgram, frame)

		// NOTE: We're not calling quad.Bund and Unbind here because we only have one VAO
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		// Draw the text
		mouse_x_f64, mouse_y_f64 := window.GetCursorPos()
		font.Printf(-0.98, 0.98, 0.5, "Mouse: %07.1f, %07.1f. Frame: %07v", mouse_x_f64, mouse_y_f64, frame)

		// Swap in the rendered buffer
		window.SwapBuffers()

		frame++
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
	gl.Scissor(0, 0, int32(width)*int32(scale_x), int32(height)*int32(scale_y))
}

// Set the mouse coordinates uniform. We assume that the shader program is already in use.
func setMouseUniform(
	window *glfw.Window,
	shaderProgram glu.ShaderProgram,
	scale_x float32,
	scale_y float32,
) {
	mouse_x_f64, mouse_y_f64 := window.GetCursorPos()
	mouse_x := float32(mouse_x_f64 * float64(scale_x))
	mouse_y := float32(mouse_y_f64 * float64(scale_y))
	shaderProgram.SetUniform2f("u_mouse", [2]float32{mouse_x, mouse_y})
}

func setTimeUniform(shaderProgram glu.ShaderProgram, frame uint32) {
	var time float32
	if frame == 0 {
		time = 0.0
	} else {
		time = float32(glfw.GetTime())
	}
	shaderProgram.SetUniform1f("u_time", time)
	shaderProgram.SetUniform1i("u_frame", int32(frame))
}
