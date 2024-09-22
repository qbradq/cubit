package c3d

import (
	"encoding/binary"
	"math"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

// VoxelMesh is a utility struct that builds voxel-based meshes.
type VoxelMesh struct {
	vao        uint32   // Vertex Array Object ID
	vbo        uint32   // Vertex Buffer Object ID
	count      int32    // Vertex count
	vboCurrent bool     // If false, the VBO needs to be reuploaded.
	d          []byte   // Raw mesh data
	vbuf       [18]byte // Vertex buffer
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
func (m *VoxelMesh) vert(x, y, z float32, r, g, b uint8, f Facing) {
	d := m.vbuf[:]
	binary.LittleEndian.PutUint32(d[0:4], math.Float32bits(x))
	binary.LittleEndian.PutUint32(d[4:8], math.Float32bits(y))
	binary.LittleEndian.PutUint32(d[8:12], math.Float32bits(z))
	d[12] = r
	d[13] = g
	d[14] = b
	d[15] = facingNormalCompressed[f][0]
	d[16] = facingNormalCompressed[f][1]
	d[17] = facingNormalCompressed[f][2]
	m.d = append(m.d, d...)
	m.count++
}

// AddFace adds a face at the given voxel position with the given facing.
// The face is scaled down to the size of one world voxel. Pos is the mesh-
// relative voxel coordinate. Note that the alpha channel is ignored by this
// function.
func (m *VoxelMesh) AddFace(pos [3]int, f Facing, c [4]uint8) {
	s := pixelScale
	p := [3]float32{
		(float32(pos[0]-8) + 0.5) * s,
		(float32(pos[1]-8) + 0.5) * s,
		(float32(pos[2]-8) + 0.5) * s,
	}
	o := cubeFacingOffsets[f]
	// X              Y                 Z           R     G     B          Normal XYZ
	m.vert(p[0]+o[0][0]*s, p[1]+o[0][1]*s, p[2]+o[0][2]*s, c[0], c[1], c[2], f) // Top left
	m.vert(p[0]+o[1][0]*s, p[1]+o[1][1]*s, p[2]+o[1][2]*s, c[0], c[1], c[2], f) // Top right
	m.vert(p[0]+o[2][0]*s, p[1]+o[2][1]*s, p[2]+o[2][2]*s, c[0], c[1], c[2], f) // Bottom left
	m.vert(p[0]+o[3][0]*s, p[1]+o[3][1]*s, p[2]+o[3][2]*s, c[0], c[1], c[2], f) // Bottom left
	m.vert(p[0]+o[4][0]*s, p[1]+o[4][1]*s, p[2]+o[4][2]*s, c[0], c[1], c[2], f) // Top right
	m.vert(p[0]+o[5][0]*s, p[1]+o[5][1]*s, p[2]+o[5][2]*s, c[0], c[1], c[2], f) // Bottom right
	m.vboCurrent = false
}

// reset rests the mesh builder state.
func (m *VoxelMesh) reset() {
	m.d = m.d[:0]
	m.vboCurrent = false
}

// draw draws the voxel mesh.
func (m *VoxelMesh) draw(p *program) {
	if m.vao == invalidVAO {
		gl.GenVertexArrays(1, &m.vao)
	}
	if m.vbo == invalidVBO {
		// Note: we have to do this on-demand because voxel meshes are loaded
		// during the mod loading phase, before the GL is initialized.
		var stride int32 = 3*4 + 3*1 + 3*1
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexPosition")),
			3, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexPosition")))
		offset += 3 * 4
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexColor")),
			3, gl.UNSIGNED_BYTE, true, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexColor")))
		offset += 3 * 1
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexNormal")),
			3, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexNormal")))
		offset += 3 * 1
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
