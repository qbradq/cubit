package c3d

import gl "github.com/go-gl/gl/v3.1/gles2"

// VoxelMesh is a utility struct that builds voxel-based meshes.
type VoxelMesh struct {
	vao        uint32    // Vertex Array Object ID
	vbo        uint32    // Vertex Buffer Object ID
	count      int32     // Vertex count
	vboCurrent bool      // If false, the VBO needs to be reuploaded.
	d          []float32 // Raw mesh data
}

// NewVoxelMesh constructs a new CubeMesh object ready for use.
func NewVoxelMesh() *VoxelMesh {
	ret := &VoxelMesh{
		vao: invalidVAO,
		vbo: invalidVBO,
		d:   make([]float32, 0, 54*16),
	}
	return ret
}

// AddFace adds a face at the given voxel position with the given facing.
// The face is scaled down to the size of one world voxel. Pos is the mesh-
// relative voxel coordinate. Note that the alpha channel is ignored by this
// function.
func (m *VoxelMesh) AddFace(pos [3]int, facing Facing, color [4]uint8) {
	s := pixelScale
	p := [3]float32{
		(float32(pos[0]-8) + 0.5) * s,
		(float32(pos[1]-8) + 0.5) * s,
		(float32(pos[2]-8) + 0.5) * s,
	}
	n := FaceNormals[facing]
	o := cubeFacingOffsets[facing]
	c := [3]float32{
		float32(color[0]) / float32(0xFF),
		float32(color[1]) / float32(0xFF),
		float32(color[2]) / float32(0xFF),
	}
	m.d = append(m.d,
		// X              Y                 Z           R     G     B          Normal XYZ
		p[0]+o[0][0]*s, p[1]+o[0][1]*s, p[2]+o[0][2]*s, c[0], c[1], c[2], n[0], n[1], n[2], // Top left
		p[0]+o[1][0]*s, p[1]+o[1][1]*s, p[2]+o[1][2]*s, c[0], c[1], c[2], n[0], n[1], n[2], // Top right
		p[0]+o[2][0]*s, p[1]+o[2][1]*s, p[2]+o[2][2]*s, c[0], c[1], c[2], n[0], n[1], n[2], // Bottom left
		p[0]+o[3][0]*s, p[1]+o[3][1]*s, p[2]+o[3][2]*s, c[0], c[1], c[2], n[0], n[1], n[2], // Bottom left
		p[0]+o[4][0]*s, p[1]+o[4][1]*s, p[2]+o[4][2]*s, c[0], c[1], c[2], n[0], n[1], n[2], // Top right
		p[0]+o[5][0]*s, p[1]+o[5][1]*s, p[2]+o[5][2]*s, c[0], c[1], c[2], n[0], n[1], n[2], // Bottom right
	)
	m.count = int32(len(m.d))
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
		var stride int32 = 3*4 + 3*4 + 3*4
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexPosition")),
			3, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexPosition")))
		offset += 3 * 4
		gl.VertexAttribPointerWithOffset(uint32(p.attr("aVertexColor")),
			3, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(p.attr("aVertexColor")))
		offset += 3 * 4
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
}
