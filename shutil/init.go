package shutil

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func Init() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
	// gl.Enable(gl.BLEND)
	// gl.Enable(gl.SCISSOR_TEST)
	// gl.BlendEquation(gl.FUNC_ADD)
}
