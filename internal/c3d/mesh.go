package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// Mesh provides a simple, drawable mesh.
type Mesh struct {
	vao    uint32 // Vertex Array Object ID
	vbo    uint32 // Vertex Buffer Object ID
	count  int32  // Vertex count
	uWorld int32  // World uniform ID
}

// newMesh creates a new generic mesh.
func newMesh(vertexes []float32, prg *Program) *Mesh {
	ret := &Mesh{
		count:  int32(len(vertexes)),
		uWorld: prg.GetUniformLocation("world"),
	}
	aPos := prg.GetAttributeLocation("pos")
	aTex := prg.GetAttributeLocation("texPos")
	gl.GenVertexArrays(1, &ret.vao)
	gl.GenBuffers(1, &ret.vbo)
	gl.BindVertexArray(ret.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, ret.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexes)*4, gl.Ptr(vertexes), gl.STATIC_DRAW)
	var stride int32 = 3*4 + 3*4
	var offset int = 0
	gl.VertexAttribPointerWithOffset(uint32(aPos), 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aPos))
	offset += 3 * 4
	gl.VertexAttribPointerWithOffset(uint32(aTex), 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(aTex))
	offset += 3 * 4
	gl.BindVertexArray(0)
	return ret
}

// draw draws the mesh.
func (m *Mesh) draw(o *Orientation) {
	mTranslate := mgl32.Translate3D(o.Position[0], o.Position[1], o.Position[2])
	mTransform := mTranslate.Mul4(o.RotationMatrix())
	gl.UniformMatrix4fv(int32(m.uWorld), 1, false, &mTransform[0])
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)
	gl.BindVertexArray(0)
}
