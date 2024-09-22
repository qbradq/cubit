package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

const pixelScale = float32(1) / float32(16)
const invalidVAO = 0xFFFFFFFF
const invalidVBO = 0xFFFFFFFF

// CubeMesh is a utility struct that builds cube-based meshes.
type CubeMesh struct {
	vao        uint32    // Vertex Array Object ID
	vbo        uint32    // Vertex Buffer Object ID
	count      int32     // Vertex count
	vboCurrent bool      // If false, the VBO needs to be reuploaded.
	d          []float32 // Raw mesh data
	voxel      bool      // If true, the mesh will be drawn in draw mode 0
}

// NewCubeMesh constructs a new CubeMesh object ready for use. If the voxel
// parameter is true, the cube mesh is rendered as solid color (Draw Mode 0).
func NewCubeMesh(voxel bool) *CubeMesh {
	ret := &CubeMesh{
		vao:   invalidVAO,
		vbo:   invalidVBO,
		voxel: voxel,
	}
	return ret
}

// AddFace adds a face at the given position with the given normal. The face
// will be placed 0.5 units from the position in the direction of the normal.
func (m *CubeMesh) AddFace(p mgl32.Vec3, facing Facing, f FaceIndex) {
	tl, br := f.ToUV()
	n := FaceNormals[facing]
	o := cubeFacingOffsets[facing]
	m.d = append(m.d,
		// X          Y             Z             S      T      Normal XYZ
		p[0]+o[0][0], p[1]+o[0][1], p[2]+o[0][2], tl[0], tl[1], n[0], n[1], n[2], // Top left
		p[0]+o[1][0], p[1]+o[1][1], p[2]+o[1][2], br[0], tl[1], n[0], n[1], n[2], // Top right
		p[0]+o[2][0], p[1]+o[2][1], p[2]+o[2][2], tl[0], br[1], n[0], n[1], n[2], // Bottom left
		p[0]+o[3][0], p[1]+o[3][1], p[2]+o[3][2], tl[0], br[1], n[0], n[1], n[2], // Bottom left
		p[0]+o[4][0], p[1]+o[4][1], p[2]+o[4][2], br[0], tl[1], n[0], n[1], n[2], // Top right
		p[0]+o[5][0], p[1]+o[5][1], p[2]+o[5][2], br[0], br[1], n[0], n[1], n[2], // Bottom right
	)
	m.count = int32(len(m.d))
	m.vboCurrent = false
}

// Reset rests the mesh builder state.
func (m *CubeMesh) Reset() {
	m.d = m.d[:0]
	m.vboCurrent = false
}

// draw draws the cube mesh.
func (m *CubeMesh) draw(p *program) {
	if m.vao == invalidVAO {
		gl.GenVertexArrays(1, &m.vao)
	}
	if m.vbo == invalidVBO {
		var stride int32 = 3*4 + 2*4 + 3*4
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexPosition")),
			3, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexPosition")))
		offset += 3 * 4
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexUV")),
			2, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexUV")))
		offset += 2 * 4
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexNormal")),
			3, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexNormal")))
		offset += 3 * 4
	}
	if !m.vboCurrent {
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(m.d)*4, gl.Ptr(m.d), gl.STATIC_DRAW)
		m.vboCurrent = true
	}
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)
	gl.BindVertexArray(0)
}
