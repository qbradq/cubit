package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// AxisIndicator provides a simple axis indicator made of full-bright lines.
type AxisIndicator struct {
	o         Orientation
	vao       uint32 // Vertex Array Object ID
	vbo       uint32 // Vertex Buffer Object ID
	uWorld    int32  // World uniform ID
	uDrawMode int32  // Draw mode uniform ID
	count     int32  // Data element count
}

// newAxisIndicator creates a new axis indicator.
func newAxisIndicator(pos mgl32.Vec3, prg *Program) *AxisIndicator {
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
		o:         *NewOrientation(pos, 0, 0, 0),
		uWorld:    prg.GetUniformLocation("world"),
		uDrawMode: prg.GetUniformLocation("drawMode"),
		count:     int32(len(vertexes)),
	}
	aPos := prg.GetAttributeLocation("pos")
	aTex := prg.GetAttributeLocation("texPos")
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
	gl.VertexAttribPointerWithOffset(uint32(aTex), 3, gl.FLOAT, false, stride,
		uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aTex))
	offset += 3 * 4
	gl.BindVertexArray(0)
	return ret
}

// Draw draws the axis indicator.
func (a *AxisIndicator) Draw() {
	mTranslate := mgl32.Translate3D(a.o.Position[0], a.o.Position[1],
		a.o.Position[2])
	mTransform := mTranslate.Mul4(a.o.RotationMatrix())
	gl.Uniform1i(a.uDrawMode, 3)
	gl.UniformMatrix4fv(int32(a.uWorld), 1, false, &mTransform[0])
	gl.BindVertexArray(a.vao)
	gl.DrawArrays(gl.LINES, 0, a.count)
	gl.BindVertexArray(0)
}
