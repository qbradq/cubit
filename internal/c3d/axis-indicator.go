package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// AxisIndicator provides a simple axis indicator made of full-bright lines.
type AxisIndicator struct {
	o     Orientation
	vao   uint32 // Vertex Array Object ID
	vbo   uint32 // Vertex Buffer Object ID
	count int32  // Data element count
}

// newAxisIndicator creates a new axis indicator.
func newAxisIndicator(pos mgl32.Vec3, prg *program) *AxisIndicator {
	vertexes := []float32{
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
		o:     *NewOrientation(pos, 0, 0, 0),
		count: int32(len(vertexes)),
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
	var stride int32 = 3*4 + 3*4
	var offset int = 0
	gl.VertexAttribPointerWithOffset(uint32(aPos), 3, gl.FLOAT, false, stride,
		uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aPos))
	offset += 3 * 4
	gl.VertexAttribPointerWithOffset(uint32(aColor), 3, gl.FLOAT, false, stride,
		uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aColor))
	offset += 3 * 4
	gl.BindVertexArray(0)
	return ret
}

// draw draws the axis indicator.
func (a *AxisIndicator) draw() {
	gl.BindVertexArray(a.vao)
	gl.DrawArrays(gl.LINES, 0, a.count)
	gl.BindVertexArray(0)
}
