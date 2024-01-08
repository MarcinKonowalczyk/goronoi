package font

import (
	"bytes"
	"fmt"
	"voronoi/glu"

	"github.com/go-gl/gl/v3.3-core/gl"

	_ "embed"
)

//go:embed font.frag
var fontFragmentShaderSource string

//go:embed font.vert
var fontVertexShaderSource string

func newFontProgram(windowWidth int, windowHeight int) glu.ShaderProgram {
	shaders := []glu.Shader{
		glu.CompileShader(fontVertexShaderSource, glu.VERTEX_SHADER),
		glu.CompileShader(fontFragmentShaderSource, glu.FRAGMENT_SHADER),
	}

	return glu.LinkShaders(shaders)
}

// Loads the specified font bytes at the given scale.
func NewFont(buf []byte, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	program := newFontProgram(windowWidth, windowHeight)

	fd := bytes.NewReader(buf)
	f, err := LoadTrueTypeFont(program, fd, scale, 32, 127)
	if err != nil {
		return nil, err
	}

	f.UpdateResolution(windowWidth, windowHeight)
	f.SetColor(1.0, 1.0, 1.0, 1.0) // Set the default color to white

	return f, nil
}

// // LoadFont loads the specified font at the given scale.
// func LoadFont(file string, scale int32, windowWidth int, windowHeight int) (*Font, error) {
// 	fd, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer fd.Close()

// 	program := newFontProgram(windowWidth, windowHeight)

// 	f, err := LoadTrueTypeFont(program, fd, scale, 32, 127)
// 	if err != nil {
// 		return nil, err
// 	}

// 	f.UpdateResolution(windowWidth, windowHeight)
// 	f.SetColor(1.0, 1.0, 1.0, 1.0) // Set the default color to white

// 	return f, nil
// }

// SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(red float32, green float32, blue float32, alpha float32) {
	f.program.Use()
	defer f.program.Unuse()
	color := [4]float32{red, green, blue, alpha}
	f.program.SetUniform4f("textColor", color)
	f.color = color
}

// UpdateResolution used to recalibrate fonts for new window size
func (f *Font) UpdateResolution(windowWidth int, windowHeight int) {
	f.program.Use()
	defer f.program.Unuse()
	f.program.SetUniform2f("u_resolution", [2]float32{float32(windowWidth), float32(windowHeight)})
	f.windowWidth = windowWidth
	f.windowHeight = windowHeight
}

// Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font) Printf(x_norm, y_norm float32, scale float32, fs string, argv ...interface{}) error {

	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return nil
	}

	// *_norm is the normalized * position of the text in the range [-1, 1]
	x := (x_norm + 1) / 2 * float32(f.windowWidth)
	y := (y_norm + 1) / 2 * float32(f.windowHeight)

	// Activate corresponding render state
	f.program.Use()
	defer f.program.Unuse()

	gl.ActiveTexture(gl.TEXTURE0)
	f.vertex_array.Bind()
	defer f.vertex_array.Unbind()
	defer gl.BindTexture(gl.TEXTURE_2D, 0)

	// Set blending options
	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// defer gl.Disable(gl.BLEND)

	// Iterate through all characters in string
	for i := range indices {

		// get rune
		runeIndex := indices[i]

		// find rune in fontChar list
		ch, ok := f.fontChar[runeIndex]

		// load missing runes in batches of 32
		if !ok {
			low := runeIndex - (runeIndex % 32)
			f.GenerateGlyphs(low, low+31)
			ch, ok = f.fontChar[runeIndex]
		}

		// skip runes that are not in font character range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		// calculate position and size for current rune
		xpos := x + float32(ch.bearingH)*scale
		ypos := y - float32(ch.height-ch.bearingV)*scale
		w := float32(ch.width) * scale
		h := float32(ch.height) * scale

		// order: bottom left, top left, bottom right, top right
		vertices := []float32{
			xpos, ypos, 0.0, 0.0,
			xpos + w, ypos, 1.0, 0.0,
			xpos, ypos + h, 0.0, 1.0,
			xpos + w, ypos + h, 1.0, 1.0,
		}

		gl.BindTexture(gl.TEXTURE_2D, ch.textureID)
		f.vertex_array.BufferData(vertices)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		x += float32((ch.advance >> 6)) * scale
	}

	return nil
}

// Width returns the width of a piece of text in pixels
func (f *Font) Width(scale float32, fs string, argv ...interface{}) float32 {

	var width float32

	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return 0
	}

	// Iterate through all characters in string
	for i := range indices {

		// get rune
		runeIndex := indices[i]

		// find rune in fontChar list
		ch, ok := f.fontChar[runeIndex]

		// load missing runes in batches of 32
		if !ok {
			low := runeIndex & rune(32-1)
			f.GenerateGlyphs(low, low+31)
			ch, ok = f.fontChar[runeIndex]
		}

		// skip runes that are not in font character range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

	}

	return width
}
