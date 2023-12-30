package shutil

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type VAO struct {
	id   uint32
	vbo  uint32
	size int32
}

// Create a Vertex Array Object (VAO) for the given vertices.
// Returns the VAO ID.
func CreateVAO(vertices []float32, size int32) VAO {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(vao)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// specify the format of our vertex input
	stride := int32(size * 4) // 4 bytes per float32
	gl.VertexAttribPointer(0, size, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO{
		id:   vao,
		vbo:  vbo,
		size: size,
	}
}

func (vao VAO) Bind() {
	gl.BindVertexArray(vao.id)
}

func (vao VAO) Unbind() {
	gl.BindVertexArray(0)
}

// Send new data to the VAO.
func (vao VAO) BufferData(vertices []float32) {
	gl.BindVertexArray(vao.id)
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.vbo)
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
