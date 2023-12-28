package shutil

import "github.com/go-gl/gl/v3.3-core/gl"

type VAO struct {
	id uint32
}

// Create a Vertex Array Object (VAO) for the given vertices.
// Returns the VAO ID.
func CreateVAO(vertices []float32) VAO {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(vao)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// specify the format of our vertex input
	// (shader) input 0
	// vertex has size 3
	// vertex items are of type FLOAT
	// do not normalize (already done)
	// stride of 3 * sizeof(float) (separation of vertices)
	// offset of where the position data starts (0 for the beginning)
	// gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO{
		id: vao,
	}
}

func (vao VAO) Bind() {
	gl.BindVertexArray(vao.id)
}

func (vao VAO) Unbind() {
	gl.BindVertexArray(0)
}
