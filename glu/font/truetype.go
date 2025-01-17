package font

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"voronoi/glu"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	fontChar map[rune]*character
	ttf      *truetype.Font

	// Font size in pixels. 12 is a good default.
	size int32

	// Font upscaling factor. This is to compensate for high DPI monitors where content scale is > 1.
	upscale float32

	// // Font compression factor
	// x_condense float32

	vertex_array glu.VartexArray
	program      glu.ShaderProgram
	windowWidth  int
	windowHeight int
	color        [4]float32
}

type character struct {
	textureID uint32 // ID handle of the glyph texture
	width     int    //glyph width
	height    int    //glyph height
	advance   int    //glyph advance
	bearingH  int    //glyph bearing horizontal
	bearingV  int    //glyph bearing vertical
}

// GenerateGlyphs builds a set of textures based on a ttf files gylphs
func (f *Font) GenerateGlyphs(low, high rune) error {
	//create a freetype context for drawing
	c := freetype.NewContext()
	c.SetDPI(float64(72 * f.upscale))
	c.SetFont(f.ttf)
	c.SetFontSize(float64(f.size))
	c.SetHinting(font.HintingFull)

	// Create new face to measure glyph dimensions
	ttfFace := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(f.size),
		DPI:     float64(72 * f.upscale),
		Hinting: font.HintingFull,
	})

	// Make each glyph
	for ch := low; ch <= high; ch++ {
		char := new(character)

		gBnd, gAdv, ok := ttfFace.GlyphBounds(ch)
		if !ok {
			return fmt.Errorf("ttf face glyphBounds error")
		}

		gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
		gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)

		// If glyph has no dimensions set to a max value
		if gw == 0 || gh == 0 {
			gBnd = f.ttf.Bounds(fixed.Int26_6(f.size))
			gw = int32((gBnd.Max.X - gBnd.Min.X) >> 6)
			gh = int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

			// Above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}

		//The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y.
		gAscent := int(-gBnd.Min.Y) >> 6
		gdescent := int(gBnd.Max.Y) >> 6

		//set w,h and adv, bearing V and bearing H in char
		char.width = int(gw)
		char.height = int(gh)
		char.advance = int(gAdv)
		char.bearingV = gdescent
		char.bearingH = (int(gBnd.Min.X) >> 6)

		//create image to draw glyph
		fg, bg := image.White, image.Black
		rect := image.Rect(0, 0, int(gw), int(gh))
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

		//set the glyph dot
		px := 0 - (int(gBnd.Min.X) >> 6)
		py := (gAscent)
		pt := freetype.Pt(px, py)

		// Draw the text from mask to image
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(fg)
		_, err := c.DrawString(string(ch), pt)
		if err != nil {
			return err
		}

		// Generate texture
		var texture uint32
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0,
			gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

		char.textureID = texture

		//add char to fontChar list
		f.fontChar[ch] = char
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return nil
}

// LoadTrueTypeFont builds OpenGL buffers and glyph textures based on a ttf file
func LoadTrueTypeFont(
	program glu.ShaderProgram,
	r io.Reader,
	size int32,
	upscale float32,
	low, high rune,
) (*Font, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font.
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	//make Font stuct type
	f := &Font{
		fontChar: make(map[rune]*character),
		ttf:      ttf,
		size:     size,
		upscale:  upscale,
		program:  program,
	}

	err = f.GenerateGlyphs(low, high)
	if err != nil {
		return nil, err
	}

	// Configure VAO/VBO for texture quads
	gl.GenVertexArrays(1, &f.vertex_array.Vao)
	gl.GenBuffers(1, &f.vertex_array.Vbo)
	gl.BindVertexArray(f.vertex_array.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, f.vertex_array.Vbo)

	// gl.BufferData(gl.ARRAY_BUFFER, 0*4*4, nil, gl.STATIC_DRAW)

	vertAttrib := program.GetAttribLocation("vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	defer gl.DisableVertexAttribArray(vertAttrib)

	texCoordAttrib := program.GetAttribLocation("vertTexCoord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	defer gl.DisableVertexAttribArray(texCoordAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return f, nil
}
