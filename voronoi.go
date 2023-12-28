package main

import (
	"log"
	"runtime"
	"voronoi/shutil"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

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

	window.SetKeyCallback(keyCallback)

	// Cap the framerate at 60fps
	glfw.SwapInterval(1)

	programLoop(window)
}

func compileShaders() []shutil.Shader {
	if vertexShaderSource == "" || fragmentShaderSource == "" {
		panic("vertexShaderSource or fragmentShaderSource is empty")
	}
	vertexShader := shutil.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := shutil.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	return []shutil.Shader{vertexShader, fragmentShader}
}

func programLoop(window *glfw.Window) {
	// the linked shader program determines how the data will be rendered
	shaders := compileShaders()
	shaderProgram := shutil.LinkShaders(shaders)

	quad_vertices := []float32{
		0.9, 0.9, 0.0,
		0.9, -0.9, 0.0,
		-0.9, 0.9, 0.0,
		-0.9, -0.9, 0.0,
	}

	quad := shutil.CreateVAO(quad_vertices)

	// Scale the resolution to the monitor's content scale
	// This is necessary for retina displays
	monitor := glfw.GetPrimaryMonitor()
	scale_x, scale_y := monitor.GetContentScale()

	// active_uniforms := shaderProgram.GetActiveUniforms()
	// fmt.Println("Active Uniforms:")
	// for _, uniform := range active_uniforms {
	// 	fmt.Println(uniform)
	// }

	shaderProgram.Use()
	shaderProgram.SetUniform2f("u_resolution", float32(windowWidth*scale_x), float32(windowHeight*scale_y))
	shaderProgram.SetUniform2f("u_mouse", float32(0.0), float32(0.0))
	shaderProgram.SetUniform1f("u_time", 0.0)

	var mouse_x, mouse_y float32

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// Get current mouse position
		mouse_x_f64, mouse_y_f64 := window.GetCursorPos()
		mouse_x = float32(mouse_x_f64)
		mouse_y = float32(mouse_y_f64)

		// Get current time
		time := glfw.GetTime()

		// Set uniforms
		shaderProgram.SetUniform2f("u_mouse", float32(mouse_x*scale_x), float32(mouse_y*scale_y))
		shaderProgram.SetUniform1f("u_time", float32(time))

		// perform rendering
		gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// draw loop
		quad.Bind()
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		quad.Unbind()
		// end of draw loop

		// swap in the rendered buffer
		window.SwapBuffers()
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
