package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// AxisIndicator provides a simple axis indicator made of full-bright lines.
type AxisIndicator struct {
	o   Orientation
	vao uint32 // Vertex Array Object ID
	vbo uint32 // Vertex Buffer Object ID
}

// newAxisIndicator creates a new axis indicator.
func newAxisIndicator(pos mgl32.Vec3, prg *program) *AxisIndicator {
	vertexes := []byte{
		// X axis
		0, 0, 0, 1, 0, 0,
		1, 0, 0, 1, 0, 0,
		// Y axis
		0, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 1, 0,
		// Z axis
		0, 0, 0, 0, 0, 1,
		0, 0, 1, 0, 0, 1,
	}
	ret := &AxisIndicator{
		o: *NewOrientation(pos, 0, 0, 0),
	}
	prg.use()
	aPos := prg.attr("aVertexPosition")
	aColor := prg.attr("aVertexColor")
	gl.GenVertexArrays(1, &ret.vao)
	gl.GenBuffers(1, &ret.vbo)
	gl.BindVertexArray(ret.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, ret.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexes)*4, gl.Ptr(vertexes),
		gl.STATIC_DRAW)
	var stride int32 = 3*1 + 3*1
	var offset int = 0
	gl.VertexAttribPointerWithOffset(uint32(aPos), 3, gl.UNSIGNED_BYTE, false,
		stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aPos))
	offset += 3 * 1
	gl.VertexAttribPointerWithOffset(uint32(aColor), 3, gl.UNSIGNED_BYTE, false,
		stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aColor))
	offset += 3 * 1
	gl.BindVertexArray(0)
	return ret
}

// draw draws the axis indicator.
func (a *AxisIndicator) draw() {
	gl.BindVertexArray(a.vao)
	gl.DrawArrays(gl.LINES, 0, 6)
	gl.BindVertexArray(0)
}
