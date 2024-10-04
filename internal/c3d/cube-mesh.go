package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/qbradq/cubit/internal/t"
)

const invalidVAO = 0xFFFFFFFF
const invalidVBO = 0xFFFFFFFF

var facingLightLevels = [6]byte{223, 223, 191, 191, 255, 127}

// CubeMesh is a utility struct that builds cube-based meshes.
type CubeMesh struct {
	vao        uint32    // Vertex Array Object ID
	vbo        uint32    // Vertex Buffer Object ID
	count      int32     // Vertex count
	vboCurrent bool      // If false, the VBO needs to be reuploaded.
	d          []byte    // Raw mesh data
	defs       []*t.Cube // List of cube definitions
}

// NewCubeMesh constructs a new CubeMesh object ready for use.
func NewCubeMesh(defs []*t.Cube) *CubeMesh {
	ret := &CubeMesh{
		vao:  invalidVAO,
		vbo:  invalidVBO,
		defs: defs,
	}
	return ret
}

// vert adds a single vertex to the data buffer.
func (m *CubeMesh) vert(x, y, z, u, v uint8, i int, c t.Cell, f t.Facing) {
	cube, _, _ := c.Decompose()
	if int(cube) >= len(m.defs) {
		return
	}
	cd := m.defs[cube]
	face := cd.Faces[f]
	fx, fy := face.ToAtlasXY()
	m.d = append(m.d, x, y, z, uint8(fx), uint8(fy), u, v, facingLightLevels[f])
	m.count++
	m.vboCurrent = false
}

// Reset rests the mesh builder state.
func (m *CubeMesh) Reset() {
	m.d = m.d[:0]
	m.count = 0
	m.vboCurrent = true
}

// draw draws the cube mesh.
func (m *CubeMesh) draw(p *program) {
	if m.vao == invalidVAO {
		gl.GenVertexArrays(1, &m.vao)
	}
	if m.vbo == invalidVBO {
		var stride int32 = 3*1 + 2*1 + 2*1 + 1*1
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexPosition")),
			3, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexPosition")))
		offset += 3 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aAtlasXY")),
			2, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aAtlasXY")))
		offset += 2 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexUV")),
			2, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexUV")))
		offset += 2 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexLightLevel")),
			1, gl.UNSIGNED_BYTE, true, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexLightLevel")))
		offset += 1 * 1
	}
	if !m.vboCurrent {
		if len(m.d) > 0 {
			gl.BindVertexArray(m.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
			gl.BufferData(gl.ARRAY_BUFFER, len(m.d), gl.Ptr(m.d), gl.STATIC_DRAW)
		}
		m.vboCurrent = true
	}
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)
}
