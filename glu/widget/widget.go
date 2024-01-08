package widget

import (
	"voronoi/glu"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Widget struct {
	windowWidth  int
	windowHeight int

	// The position of the widget top left corner in screen coordinates
	x int
	y int

	// The size of the widget in screen coordinates
	width  int
	height int

	color [4]float32

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

uniform vec4 u_color;

void main()
{
	// color = u_color;

	// Round the corners of the rectangle
	float radius = 0.1;
	vec2 pos = gl_FragCoord.xy;
	vec2 size = vec2(0.6, 0.6);
	vec2 center = vec2(0.0, 0.0);
	vec2 dist = abs(pos - center) - size + vec2(radius);
	float alpha = 1.0 - clamp(max(dist.x, dist.y), 0.0, radius);

	color = u_color;



	// color = vec4(1.0, 1.0, 1.0, 1.0);
}
`

func newWidgetProgram(windowWidth int, windowHeight int) glu.ShaderProgram {
	shaders := []glu.Shader{
		glu.CompileShader(widgetVertexShader, glu.VERTEX_SHADER),
		glu.CompileShader(widgetFragmentShader, glu.FRAGMENT_SHADER),
	}

	return glu.LinkShaders(shaders)
}

func NewWidget(windowWidth int, windowHeight int) *Widget {
	program := newWidgetProgram(windowWidth, windowHeight)

	w := &Widget{
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		color:        [4]float32{1.0, 1.0, 1.0, 1.0},
		program:      program,
	}

	w.vertex_array = glu.NewVertexArray(2)

	quad_vertices := []float32{
		0.3, 0.3,
		0.3, -0.3,
		-0.3, 0.3,
		-0.3, -0.3,
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

// func (w *Widget) UpdateResolution(windowWidth int, windowHeight int) {
// 	w.windowWidth = windowWidth
// 	w.windowHeight = windowHeight
// }

func (w *Widget) SetPosition(x int, y int) {
	w.x = x
	w.y = y
}

func (w *Widget) SetSize(width int, height int) {
	w.width = width
	w.height = height
}

func (w *Widget) Draw() {
	w.program.Use()
	defer w.program.Unuse()

	w.vertex_array.Bind()
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

}
