package glu

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type VartexArray struct {
	Vao uint32
	Vbo uint32
}

// Create a Vertex Array Object (VAO) for the given vertices.
// Returns the VAO ID.
func NewVertexArray(size int32) (va VartexArray) {
	gl.GenVertexArrays(1, &va.Vao)
	gl.GenBuffers(1, &va.Vbo)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(va.Vao)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, va.Vbo)
	// gl.BufferData(gl.ARRAY_BUFFER, 0, nil, gl.STATIC_DRAW)

	// specify the format of our vertex input
	stride := int32(size * 4) // 4 bytes per float32
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, size, gl.FLOAT, false, stride, 0)
	defer gl.DisableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return va
}

func (va VartexArray) Bind() {
	gl.BindVertexArray(va.Vao)
}

func (va VartexArray) Unbind() {
	gl.BindVertexArray(0)
}

// Send new data to the VAO.
func (va VartexArray) BufferData(vertices []float32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, va.Vbo)
	defer gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Read the data back to make sure it was written correctly
	var data []float32 = make([]float32, len(vertices))
	gl.GetBufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(&data[0]))

	// Compare the data
	for i := 0; i < len(vertices); i++ {
		if vertices[i] != data[i] {
			fmt.Printf("vertices[%d] = %f, data[%d] = %f\n", i, vertices[i], i, data[i])
			panic("vertices != data")
		}
	}

}
