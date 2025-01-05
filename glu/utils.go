package glu

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func Init() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)
}

func ClearColor(r, g, b, a float32) {
	gl.ClearColor(0.137, 0.137, 0.137, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
