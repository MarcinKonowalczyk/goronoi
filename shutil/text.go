package shutil

import (
	_ "embed"
	"fmt"
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

// For now we will only support ASCII characters
const MIN_CHAR = ' ' // 32
const MAX_CHAR = '~' // 126

// Text renderer holds all the state necessary for rendering text. The size is fixed at creation time.
// New text renderers can be created for different sizes.
type TextRenderer struct {
	program ShaderProgram
	vao     VAO
	size    float32
	atlas   [MAX_CHAR - MIN_CHAR]Glyph
	Color   [4]float32
}

// All the information needed to render a glyph
type Glyph struct {
	character rune
	width     int
	height    int
	bearingX  int
	bearingY  int
	advance   int
	texture   uint32
}

//go:embed font.vert
var fontVertexShaderSource string

//go:embed font.frag
var fontFragmentShaderSource string

func createAtlas(font *truetype.Font, size float64) [MAX_CHAR - MIN_CHAR]Glyph {

	// Make sure the openGL context is active
	gl.ActiveTexture(gl.TEXTURE0)

	// Create font face
	face := truetype.NewFace(font, &truetype.Options{
		Size: size,
		DPI:  72,
	})
	// For each character in the font, create a glyph
	atlas := [MAX_CHAR - MIN_CHAR]Glyph{}
	var dr image.Rectangle
	var mask image.Image
	// var maskp image.Point
	var advance fixed.Int26_6
	var ok bool

	for r := MIN_CHAR; r < MAX_CHAR; r++ {
		dot := fixed.P(0, 0)
		dr, mask, _, advance, ok = face.Glyph(dot, r)
		if !ok {
			panic(fmt.Sprintf("failed to get glyph for character %c", r))
		}

		// fmt.Println("Got glyph for character", r, "with advance", advance)

		// Create a texture for the glyph
		var texture uint32
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// Set the texture wrapping/filtering options (on the currently bound texture object)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE) // Set texture wrapping to GL_REPEAT (usually basic wrapping method)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST) // Set texture filtering (options: GL_LINEAR, GL_NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

		// Upload the glyph to the texture
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(dr.Dx()), int32(dr.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(mask.(*image.Alpha).Pix))

		// Convert bearing to pixels by dividing by 64
		advance = advance >> 6

		// Add the glyph to the atlas
		atlas[r-MIN_CHAR] = Glyph{
			character: r,
			width:     dr.Dx(),
			height:    dr.Dy(),
			bearingX:  0,
			bearingY:  0,
			advance:   int(advance),
			texture:   texture,
		}
	}

	return atlas
}

func NewTextRenderer(font *truetype.Font, size float64) TextRenderer {

	atlas := createAtlas(font, size)

	// Create the shader program for rendering text
	font_vertex_shader := CompileShader(fontVertexShaderSource, VERTEX_SHADER)
	font_fragment_shader := CompileShader(fontFragmentShaderSource, FRAGMENT_SHADER)
	program := LinkShaders([]Shader{font_vertex_shader, font_fragment_shader})

	// Create the VAO for rendering text
	// This is a quad that will be rendered for each character
	// The VBO will be updated for each character

	vertices := []float32{
		0.9, 0.9, 0.0,
		0.9, -0.9, 0.0,
		-0.9, 0.9, 0.0,
		-0.9, -0.9, 0.0,
	}
	vao := CreateVAO(vertices, 3)

	return TextRenderer{
		program: program,
		vao:     vao,
		size:    float32(size),
		atlas:   atlas,
		Color:   [4]float32{1.0, 1.0, 1.0, 1.0},
	}
}

func (t *TextRenderer) RenderText(text string, pos [2]float32, screen_size [2]float32) {
	// pos is in -1 to 1 coordinates
	// screen_size is in pixels

	apos := [2]float32{
		(pos[0] + 1) / 2 * screen_size[0],
		(pos[1] + 1) / 2 * screen_size[1],
	}

	screen_width := screen_size[0]
	screen_height := screen_size[1]

	// Activate corresponding render state
	t.program.Use()
	t.program.SetUniform4f("u_color", t.Color)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	t.vao.Bind()

	// Apply the scale to the text such that 12 point font is 12 pixels tall
	scale := float32(2.0)
	fmt.Println("scale:", scale)

	aorigin := apos[0]
	for _, c := range text {
		// Get the glyph info for the character
		glyph := t.atlas[c-MIN_CHAR]
		// fmt.Printf("Rendering character '%c' with advance %d\n", c, glyph.advance)
		// fmt.Printf("Bearing: %d, %d\n", glyph.bearingX, glyph.bearingY)

		// Print the glyph info
		// fmt.Printf("glyph: %+v\n", glyph)

		// Calculate the absolute position of the glyph in pixels
		xpos := aorigin
		ypos := apos[1]

		// Calculate the origin of the next glyph
		aorigin += float32(glyph.advance) * scale

		// Calculate the absolute size of the glyph in pixels
		width := float32(glyph.width) * scale
		height := float32(glyph.height) * scale

		// fmt.Printf("xpos: %f, ypos: %f, width: %f, height: %f\n", xpos, ypos, width, height)

		// Calculate the relative position of the glyph in -1 to 1 coordinates
		// a_vertices := []float32{
		// 	xpos, ypos, 0.0,
		// 	xpos, (ypos + height), 0.0,
		// 	(xpos + width), ypos, 0.0,
		// 	(xpos + width), (ypos + height), 0.0,
		// }

		vertices := []float32{
			(xpos - screen_width/2) / (screen_width / 2), (ypos - screen_height/2) / (screen_height / 2), 0.0,
			(xpos - screen_width/2) / (screen_width / 2), (ypos + height - screen_height/2) / (screen_height / 2), 0.0,
			(xpos + width - screen_width/2) / (screen_width / 2), (ypos - screen_height/2) / (screen_height / 2), 0.0,
			(xpos + width - screen_width/2) / (screen_width / 2), (ypos + height - screen_height/2) / (screen_height / 2),
		}

		t.vao.BufferData(vertices)

		// fmt.Println(vertices)

		// Render the glyph
		gl.BindTexture(gl.TEXTURE_2D, glyph.texture)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	}
}
