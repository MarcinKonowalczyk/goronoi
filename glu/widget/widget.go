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

	// The position of the widget top left corner in screen coordinates
	x int
	y int

	// The size of the widget in screen coordinates
	width  int
	height int

	// The color of the widget
	color [4]float32

	// The shader program used to draw the widget
	program      glu.ShaderProgram
	vertex_array glu.VartexArray
}

// For now we'll just work on drawing a simple rectangle which can be
// dragged around the screen and resized. No content yet.

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

void main()
{
	vec2 st = gl_FragCoord.xy/u_resolution.xy;
	vec2 mouse = u_mouse/u_resolution;
    mouse.y = 1.0 - mouse.y; // flip y-axis

	float mouse_dist = distance(st, mouse);
	if (mouse_dist < 0.1) {
		color = vec4(1.0, 0.0, 0.0, 0.5);
	} else {
		color = u_color;
	}

}
`

func newWidgetProgram(windowWidth int, windowHeight int) glu.ShaderProgram {
	shaders := []glu.Shader{
		glu.CompileShader(widgetVertexShader, glu.VERTEX_SHADER),
		glu.CompileShader(widgetFragmentShader, glu.FRAGMENT_SHADER),
	}

	return glu.LinkShaders(shaders)
}

func NewWidget(windowWidth int, windowHeight int, scaleX float32, scaleY float32) *Widget {
	program := newWidgetProgram(windowWidth, windowHeight)

	w := &Widget{
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		scaleX:       scaleX,
		scaleY:       scaleY,
		color:        [4]float32{1.0, 1.0, 1.0, 1.0},
		program:      program,
	}

	w.vertex_array = glu.NewVertexArray(2)

	quad_vertices := []float32{
		-0.8, 0.8,
		-0.8, 0.2,
		-0.2, 0.8,
		-0.2, 0.2,
	}

	w.vertex_array.BufferData(quad_vertices)

	w.SetColor(1.0, 1.0, 1.0, 1.0)
	w.SetWindowResolution(windowWidth, windowHeight)

	return w
}

func (w *Widget) SetColor(red float32, green float32, blue float32, alpha float32) {
	w.program.Use()
	defer w.program.Unuse()
	color := [4]float32{red, green, blue, alpha}
	w.program.SetUniform4f("u_color", color)
	w.color = color
}

// func (w *Widget) UpdateResolution(windowWidth int, windowHeight int) {
// 	w.windowWidth = windowWidth
// 	w.windowHeight = windowHeight
// }

func (w *Widget) SetPosition(x int, y int) {
	w.x = x
	w.y = y
}

func (w *Widget) SetWindowResolution(width int, height int) {
	w.program.Use()
	defer w.program.Unuse()

	w.windowWidth = width
	w.windowHeight = height
	w.program.SetUniform2f("u_resolution", [2]float32{float32(w.windowWidth), float32(w.windowHeight)})
}

func (w *Widget) SetMouseUniform(mouse_x_f64 float64, mouse_y_f64 float64) {
	w.program.Use()
	defer w.program.Unuse()

	mouse_x := float32(mouse_x_f64 * float64(w.scaleX))
	mouse_y := float32(mouse_y_f64*float64(w.scaleY)) - float32(w.windowHeight)

	w.program.SetUniform2f("u_mouse", [2]float32{mouse_x, mouse_y})
}

func (w *Widget) Draw() {
	w.program.Use()
	defer w.program.Unuse()

	w.vertex_array.Bind()
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

}
