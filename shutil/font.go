package shutil

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"

	_ "embed"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

// Known directions.
const (
	LeftToRight Direction = iota // E.g.: Latin
	RightToLeft                  // E.g.: Arabic
	TopToBottom                  // E.g.: Chinese
)

//go:embed font.frag
var fragmentFontShader string

//go:embed font.vert
var vertexFontShader string

func newFontProgram(windowWidth int, windowHeight int) ShaderProgram {
	shaders := []Shader{
		CompileShader(vertexFontShader, VERTEX_SHADER),
		CompileShader(fragmentFontShader, FRAGMENT_SHADER),
	}

	return LinkShaders(shaders)
}

// // Activate corresponding render state

// // set screen resolution
// resUniform := gl.GetUniformLocation(program, gl.Str("resolution\x00"))
// gl.Uniform2f(resUniform, float32(windowWidth), float32(windowHeight))

// LoadFontBytes loads the specified font bytes at the given scale.
func LoadFontBytes(buf []byte, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	program := newFontProgram(windowWidth, windowHeight)

	fd := bytes.NewReader(buf)
	return LoadTrueTypeFont(program, fd, scale, 32, 127, LeftToRight)
}

// LoadFont loads the specified font at the given scale.
func LoadFont(file string, scale int32, windowWidth int, windowHeight int) (*Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	program := newFontProgram(windowWidth, windowHeight)

	return LoadTrueTypeFont(program, fd, scale, 32, 127, LeftToRight)
}

// SetColor allows you to set the text color to be used when you draw the text
func (f *Font) SetColor(red float32, green float32, blue float32, alpha float32) {
	f.program.Use()
	defer f.program.Unuse()
	f.program.SetUniform4f("textColor", [4]float32{red, green, blue, alpha})
}

// UpdateResolution used to recalibrate fonts for new window size
func (f *Font) UpdateResolution(windowWidth int, windowHeight int) {
	f.program.Use()
	defer f.program.Unuse()
	f.program.SetUniform2f("resolution", [2]float32{float32(windowWidth), float32(windowHeight)})
}

// Printf draws a string to the screen, takes a list of arguments like printf
func (f *Font) Printf(x, y float32, scale float32, fs string, argv ...interface{}) error {

	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return nil
	}

	// Activate corresponding render state
	f.program.Use()
	defer f.program.Unuse()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(f.vao)
	defer gl.BindVertexArray(0)
	defer gl.BindTexture(gl.TEXTURE_2D, 0)
	defer f.program.Unuse()

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

		// skip runes that are not in font chacter range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		// calculate position and size for current rune
		xpos := x + float32(ch.bearingH)*scale
		ypos := y - float32(ch.height-ch.bearingV)*scale
		w := float32(ch.width) * scale
		h := float32(ch.height) * scale
		vertices := []float32{
			xpos + w, ypos, 1.0, 0.0,
			xpos, ypos, 0.0, 0.0,
			xpos, ypos + h, 0.0, 1.0,

			xpos, ypos + h, 0.0, 1.0,
			xpos + w, ypos + h, 1.0, 1.0,
			xpos + w, ypos, 1.0, 0.0,
		}

		// Render glyph texture over quad
		gl.BindTexture(gl.TEXTURE_2D, ch.textureID)
		// Update content of VBO memory
		gl.BindBuffer(gl.ARRAY_BUFFER, f.vbo)

		// BufferSubData(target Enum, offset int, data []byte)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices)) // Be sure to use glBufferSubData and not glBufferData
		// Render quad
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
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

		// skip runes that are not in font chacter range
		if !ok {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}

		// Now advance cursors for next glyph (note that advance is number of 1/64 pixels)
		width += float32((ch.advance >> 6)) * scale // Bitshift by 6 to get value in pixels (2^6 = 64 (divide amount of 1/64th pixels by 64 to get amount of pixels))

	}

	return width
}
