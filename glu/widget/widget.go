package widget

import (
	"voronoi/glu"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Widget struct {
	// The size of the window in screen coordinates
	windowWidth  int
	windowHeight int
	scaleX       float32
	scaleY       float32

	// Mouse position in screen coordinates
	mouseX    float32
	mouseY    float32
	mouseDown bool
	mouseOver bool

	mouseXPrev    float32
	mouseYPrev    float32
	mouseDownPrev bool
	mouseOverPrev bool

	// The position of the widget
	// x      float32
	// y      float32
	// width  float32
	// height float32

	// The color of the widget
	color [4]float32

	// The shader program used to draw the widget
	program      glu.ShaderProgram
	vertex_array glu.VartexArray
}

var widgetVertexShader = `
#version 330 core

layout (location = 0) in vec2 position;

void main()
{
	gl_Position = vec4(position.xy, 0.0, 1.0);
}
`

var widgetFragmentShader = `
#version 330 core

out vec4 color;

uniform vec2 u_resolution;
uniform vec4 u_color;
uniform vec2 u_mouse;
uniform bool u_mouse_down;
uniform bool u_mouse_over;

void main()
{
	vec2 st = gl_FragCoord.xy/u_resolution.xy;
	vec2 mouse = u_mouse/u_resolution;
    mouse.y = 1.0 - mouse.y; // flip y-axis

	float mouse_dist = distance(st, mouse);
	if (mouse_dist < 0.1) {
		if (u_mouse_down) {
			// Yellow
			color = vec4(1.0, 1.0, 0.0, 0.5);
		} else {
			// Red
			color = vec4(1.0, 0.0, 0.0, 0.5);
		}
	} else {
		if (u_mouse_over) {
				color = vec4(u_color.rgb, 1.0);
			} else {
				color = vec4(u_color.rgb, 0.5);
		}
	}

}
`

func newWidgetProgram() glu.ShaderProgram {
	shaders := []glu.Shader{
		glu.CompileShader(widgetVertexShader, glu.VERTEX_SHADER),
		glu.CompileShader(widgetFragmentShader, glu.FRAGMENT_SHADER),
	}

	return glu.LinkShaders(shaders)
}

func NewWidget() *Widget {
	program := newWidgetProgram()

	w := &Widget{program: program}

	w.vertex_array = glu.NewVertexArray(2)

	// Default position of the widget
	quad_vertices := []float32{
		-0.8, 0.8,
		-0.8, 0.2,
		-0.2, 0.8,
		-0.2, 0.2,
	}

	w.vertex_array.BufferData(quad_vertices)
	w.SetColor(1.0, 1.0, 1.0, 1.0)

	return w
}

func (w *Widget) SetColor(red float32, green float32, blue float32, alpha float32) {
	w.program.Use()
	defer w.program.Unuse()
	color := [4]float32{red, green, blue, alpha}
	w.program.SetUniform4f("u_color", color)
	w.color = color
}

// Set all the stuff to do with the window size.
func (w *Widget) SetWindow(width int, height int, scale_x float32, scale_y float32) {
	w.windowWidth = width
	w.windowHeight = height
	w.scaleX = scale_x
	w.scaleY = scale_y

	w.program.Use()
	defer w.program.Unuse()
	w.program.SetUniform2f("u_resolution", [2]float32{float32(w.windowWidth), float32(w.windowHeight)})
}

// Set the mouse position and state.
func (w *Widget) SetMouse(mouse_x_f64 float64, mouse_y_f64 float64, mouse_down int) {
	// Update the previous mouse state
	w.mouseXPrev = w.mouseX
	w.mouseYPrev = w.mouseY
	w.mouseDownPrev = w.mouseDown

	// Update the current mouse state
	w.mouseX = float32(mouse_x_f64 * float64(w.scaleX))
	w.mouseY = float32(mouse_y_f64*float64(w.scaleY)) - float32(w.windowHeight)
	w.mouseDown = mouse_down != 0

	// Update the mouse over state
	// mouse_over := w.mouseX < 100
	mouse_over := false
	w.mouseOverPrev = w.mouseOver

	w.mouseOver = mouse_over

	// Set the shader uniforms
	mouse_x := float32(mouse_x_f64 * float64(w.scaleX))
	mouse_y := float32(mouse_y_f64*float64(w.scaleY)) - float32(w.windowHeight)
	mouse_over_int := 0
	if mouse_over {
		mouse_over_int = 1
	}

	w.program.Use()
	defer w.program.Unuse()
	w.program.SetUniform2f("u_mouse", [2]float32{mouse_x, mouse_y})
	w.program.SetUniform1i("u_mouse_down", int32(mouse_down))
	w.program.SetUniform1i("u_mouse_over", int32(mouse_over_int))
}

func (w *Widget) Draw() {
	w.program.Use()
	defer w.program.Unuse()

	w.vertex_array.Bind()
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}
