package widget

import (
	"fmt"
	"voronoi/glu"

	"github.com/go-gl/gl/v3.3-core/gl"

	_ "embed"
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

	// The position of the widget. x and y are in screen coordinates. width and height are in pixels.
	x      float32
	y      float32
	width  float32
	height float32

	// The color of the widget
	color [4]float32

	// The shader program used to draw the widget
	program      glu.ShaderProgram
	vertex_array glu.VartexArray
}

//go:embed widget.vert
var widgetVertexShader string

//go:embed widget.frag
var widgetFragmentShader string

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

	// w.vertex_array.BufferData(quad_vertices)
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
	// w.mouseX = float32(mouse_x_f64 * float64(w.scaleX))
	// w.mouseY = float32(mouse_y_f64*float64(w.scaleY)) - float32(w.windowHeight)
	// w.mouseY = float32(mouse_y_f64) - float32(w.windowHeight)
	fmt.Println(mouse_x_f64)
	w.mouseX = float32(mouse_x_f64)
	w.mouseY = float32(mouse_y_f64)
	w.mouseDown = mouse_down != 0

	// Maybe move the window
	// if w.mouseDown && w.mouseDownPrev && w.mouseOver && w.mouseOverPrev {
	// 	delta_x := w.mouseX - w.mouseXPrev
	// 	delta_y := w.mouseY - w.mouseYPrev
	// 	w.SetPosition(w.x+delta_x, w.y+delta_y, w.width, w.height)
	// }

	// Update the mouse over state
	x_pixels := (w.x + 1) / 2 * float32(w.windowWidth)
	y_pixels := (w.y + 1) / 2 * float32(w.windowHeight)
	// fmt.Println(x_pixels, y_pixels)
	// mouse_over := (w.mouseX > w.x) && (w.mouseX < (w.x + w.width)) && (w.mouseY > w.y) && (w.mouseY < (w.y + w.height))
	mouse_over := (w.mouseX > x_pixels) && (w.mouseX < (x_pixels + w.width)) && (w.mouseY > y_pixels) && (w.mouseY < (y_pixels + w.height))
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

func (w *Widget) SetPosition(
	x float32,
	y float32,
	width float32,
	height float32,
) {

	// y = float32(w.windowWidth) - y
	w.x = x
	w.y = y
	w.width = width
	w.height = height

	// Calculate the position in screen units and send them over to the buffer array
	xn := (w.x + 1) / 2.0
	yn := (w.y + 1) / 2.0
	wn := w.width / float32(w.windowWidth)
	hn := w.height / float32(w.windowHeight)

	t := func(x float32) float32 { return (x - 0.5) * 2 }

	// Default position of the widget
	quad_vertices := []float32{
		t(xn), -t(yn),
		t(xn), -t(yn + hn),
		t(xn + wn), -t(yn),
		t(xn + wn), -t(yn + hn),
	}
	w.vertex_array.BufferData(quad_vertices)
}

func (w *Widget) Draw() {
	w.program.Use()
	defer w.program.Unuse()

	w.vertex_array.Bind()
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}
