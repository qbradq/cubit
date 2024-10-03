package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
)

var voxelLightLevels = []float32{}

func init() {
	n := float32(223.0 / 255.0)
	s := float32(223.0 / 255.0)
	e := float32(191.0 / 255.0)
	w := float32(191.0 / 255.0)
	t := float32(1.0)
	b := float32(127.0 / 255.0)
	voxelLightLevels = []float32{
		n, s, e, w, t, b, // North
		s, n, w, e, t, b, // South
		w, e, n, s, t, b, // East
		e, w, s, n, t, b, // West
		t, b, e, w, n, s, // Top
		b, t, e, w, s, n, // Bottom
	}
}

// VoxelMesh is a utility struct that builds voxel-based meshes.
type VoxelMesh struct {
	vao        uint32 // Vertex Array Object ID
	vbo        uint32 // Vertex Buffer Object ID
	count      int32  // Vertex count
	vboCurrent bool   // If false, the VBO needs to be reuploaded.
	d          []byte // Raw mesh data
}

// NewVoxelMesh constructs a new CubeMesh object ready for use.
func NewVoxelMesh() *VoxelMesh {
	ret := &VoxelMesh{
		vao: invalidVAO,
		vbo: invalidVBO,
	}
	return ret
}

// vert adds a vertex with the given attributes.
func (m *VoxelMesh) vert(x, y, z uint8, c [4]uint8, f Facing) {
	m.d = append(m.d,
		x, y, z,
		c[0], c[1], c[2],
		byte(f),
	)
	m.count++
	m.vboCurrent = false
}

// // AddFace adds a face at the given voxel position with the given facing.
// // The face is scaled down to the size of one world voxel. Pos is the mesh-
// // relative voxel coordinate. Note that the alpha channel is ignored by this
// // function.
// func (m *VoxelMesh) AddFace(p, d [3]uint8, f Facing, c [4]uint8) {
// 	d[0] -= 1
// 	d[1] -= 1
// 	d[2] -= 1
// 	switch f {
// 	case North:
// 		m.vert(p[0]+d[0]+1, p[1]+d[1]+1, p[2], c, f) // TL
// 		m.vert(p[0], p[1]+d[1]+1, p[2], c, f)        // TR
// 		m.vert(p[0]+d[0]+1, p[1], p[2], c, f)        // BL
// 		m.vert(p[0]+d[0]+1, p[1], p[2], c, f)        // BL
// 		m.vert(p[0], p[1]+d[1]+1, p[2], c, f)        // TR
// 		m.vert(p[0], p[1], p[2], c, f)               // BR
// 	case South:
// 		m.vert(p[0], p[1]+d[1]+1, p[2]+1, c, f)        // TL
// 		m.vert(p[0]+d[0]+1, p[1]+d[1]+1, p[2]+1, c, f) // TR
// 		m.vert(p[0], p[1], p[2]+1, c, f)               // BL
// 		m.vert(p[0], p[1], p[2]+1, c, f)               // BL
// 		m.vert(p[0]+d[0]+1, p[1]+d[1]+1, p[2]+1, c, f) // TR
// 		m.vert(p[0]+d[0]+1, p[1], p[2]+1, c, f)        // BR
// 	case East:
// 		m.vert(p[0]+1, p[1]+d[1]+1, p[2]+d[2]+1, c, f) // TL
// 		m.vert(p[0]+1, p[1]+d[1]+1, p[2], c, f)        // TR
// 		m.vert(p[0]+1, p[1], p[2]+d[2]+1, c, f)        // BL
// 		m.vert(p[0]+1, p[1], p[2]+d[2]+1, c, f)        // BL
// 		m.vert(p[0]+1, p[1]+d[1]+1, p[2], c, f)        // TR
// 		m.vert(p[0]+1, p[1], p[2], c, f)               // BR
// 	case West:
// 		m.vert(p[0], p[1]+d[1]+1, p[2], c, f)        // TL
// 		m.vert(p[0], p[1]+d[1]+1, p[2]+d[2]+1, c, f) // TR
// 		m.vert(p[0], p[1], p[2], c, f)               // BL
// 		m.vert(p[0], p[1], p[2], c, f)               // BL
// 		m.vert(p[0], p[1]+d[1]+1, p[2]+d[2]+1, c, f) // TR
// 		m.vert(p[0], p[1], p[2]+d[2]+1, c, f)        // BR
// 	case Top:
// 		m.vert(p[0], p[1]+1, p[2], c, f)               // TL
// 		m.vert(p[0]+d[0]+1, p[1]+1, p[2], c, f)        // TR
// 		m.vert(p[0], p[1]+1, p[2]+d[2]+1, c, f)        // BL
// 		m.vert(p[0], p[1]+1, p[2]+d[2]+1, c, f)        // BL
// 		m.vert(p[0]+d[0]+1, p[1]+1, p[2], c, f)        // TR
// 		m.vert(p[0]+d[0]+1, p[1]+1, p[2]+d[2]+1, c, f) // BR
// 	case Bottom:
// 		m.vert(p[0]+d[0]+1, p[1], p[2], c, f)        // TL
// 		m.vert(p[0], p[1], p[2], c, f)               // TR
// 		m.vert(p[0]+d[0]+1, p[1], p[2]+d[2]+1, c, f) // BL
// 		m.vert(p[0]+d[0]+1, p[1], p[2]+d[2]+1, c, f) // BL
// 		m.vert(p[0], p[1], p[2], c, f)               // TR
// 		m.vert(p[0], p[1], p[2]+d[2]+1, c, f)        // BR
// 	}
// 	m.vboCurrent = false
// }

// Reset rests the mesh builder state.
func (m *VoxelMesh) Reset() {
	m.d = m.d[:0]
	m.count = 0
	m.vboCurrent = true
}

// draw draws the voxel mesh.
func (m *VoxelMesh) draw(p *program) {
	if m.vao == invalidVAO {
		gl.GenVertexArrays(1, &m.vao)
	}
	if m.vbo == invalidVBO {
		// Note: we have to do this on-demand because voxel meshes are loaded
		// during the mod loading phase, before the GL is initialized.
		var stride int32 = 3*1 + 3*1 + 1*1
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexPosition")),
			3, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexPosition")))
		offset += 3 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexColor")),
			3, gl.UNSIGNED_BYTE, true, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexColor")))
		offset += 3 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexFacing")),
			1, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexFacing")))
		offset += 1 * 1
	}
	if !m.vboCurrent {
		if len(m.d) != 0 {
			gl.BindVertexArray(m.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
			gl.BufferData(gl.ARRAY_BUFFER, len(m.d), gl.Ptr(m.d), gl.STATIC_DRAW)
		}
		m.vboCurrent = true
	}
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)
}
