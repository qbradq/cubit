package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
)

const pixelScale = float32(1) / float32(16)
const invalidVAO = 0xFFFFFFFF
const invalidVBO = 0xFFFFFFFF

// CubeMesh is a utility struct that builds cube-based meshes.
type CubeMesh struct {
	vao        uint32  // Vertex Array Object ID
	vbo        uint32  // Vertex Buffer Object ID
	count      int32   // Vertex count
	vboCurrent bool    // If false, the VBO needs to be reuploaded.
	d          []byte  // Raw mesh data
	vbuf       [8]byte // Vertex data buffer
}

// NewCubeMesh constructs a new CubeMesh object ready for use.
func NewCubeMesh() *CubeMesh {
	ret := &CubeMesh{
		vao: invalidVAO,
		vbo: invalidVBO,
	}
	return ret
}

// vert adds a single vertex to the data buffer.
func (m *CubeMesh) vert(x, y, z, u, v byte, f Facing) {
	d := m.vbuf[:]
	d[0] = x
	d[1] = y
	d[2] = z
	d[3] = u
	d[4] = v
	d[5] = facingNormalCompressed[f][0]
	d[6] = facingNormalCompressed[f][1]
	d[7] = facingNormalCompressed[f][2]
	m.d = append(m.d, d...)
	m.count++
}

// AddFace adds a face at the given position with the given normal. The face
// will be placed 0.5 units from the position in the direction of the normal.
func (m *CubeMesh) AddFace(x, y, z byte, f Facing, face FaceIndex) {
	tl, br := face.toCompressedUV()
	o := cubeFacingOffsetsI[f]
	// X          Y             Z             S      T      Normal XYZ
	m.vert(x+o[0][0], y+o[0][1], z+o[0][2], tl[0], tl[1], f) // Top left
	m.vert(x+o[1][0], y+o[1][1], z+o[1][2], br[0], tl[1], f) // Top right
	m.vert(x+o[2][0], y+o[2][1], z+o[2][2], tl[0], br[1], f) // Bottom left
	m.vert(x+o[3][0], y+o[3][1], z+o[3][2], tl[0], br[1], f) // Bottom left
	m.vert(x+o[4][0], y+o[4][1], z+o[4][2], br[0], tl[1], f) // Top right
	m.vert(x+o[5][0], y+o[5][1], z+o[5][2], br[0], br[1], f) // Bottom right
	m.count = int32(len(m.d))
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
		var stride int32 = 3*1 + 2*1 + 3*1
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexPosition")),
			3, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexPosition")))
		offset += 3 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexUV")),
			2, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexUV")))
		offset += 2 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexNormal")),
			3, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexNormal")))
		offset += 3 * 1
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
	gl.BindVertexArray(0)
}
