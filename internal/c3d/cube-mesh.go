package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

const pixelScale = float32(1) / float32(16)

// CubeMesh is a utility struct that builds cube-based meshes.
type CubeMesh struct {
	vao        uint32    // Vertex Array Object ID
	vbo        uint32    // Vertex Buffer Object ID
	count      int32     // Vertex count
	uWorld     int32     // World uniform ID
	glStarted  bool      // If true, all OpenGL state changes required have happened
	vboCurrent bool      // If false, the VBO needs to be reuploaded.
	d          []float32 // Raw mesh data
	aPos       int32     // pos attribute location
	aTex       int32     // tex attribute location
	aNorm      int32     // norm attribute location
	voxel      bool      // If true, the mesh will be drawn in draw mode 0
}

// NewCubeMesh constructs a new CubeMesh object ready for use. If the voxel
// parameter is true, the cube mesh is rendered as solid color (Draw Mode 0).
func NewCubeMesh(voxel bool) *CubeMesh {
	return &CubeMesh{
		d:     make([]float32, 0, 54*16),
		voxel: voxel,
	}
}

// AddFace adds a face at the given position with the given normal. The face
// will be placed 0.5 units from the position in the direction of the normal.
func (m *CubeMesh) AddFace(p mgl32.Vec3, facing Facing, f FaceIndex) {
	tl, br := f.ToST()
	n := FaceNormals[facing]
	o := cubeFacingOffsets[facing]
	m.d = append(m.d,
		// X          Y             Z             S      T      None  Normal XYZ
		p[0]+o[0][0], p[1]+o[0][1], p[2]+o[0][2], tl[0], tl[1], 0, n[0], n[1], n[2], // Top left
		p[0]+o[1][0], p[1]+o[1][1], p[2]+o[1][2], br[0], tl[1], 0, n[0], n[1], n[2], // Top right
		p[0]+o[2][0], p[1]+o[2][1], p[2]+o[2][2], tl[0], br[1], 0, n[0], n[1], n[2], // Bottom left
		p[0]+o[3][0], p[1]+o[3][1], p[2]+o[3][2], tl[0], br[1], 0, n[0], n[1], n[2], // Bottom left
		p[0]+o[4][0], p[1]+o[4][1], p[2]+o[4][2], br[0], tl[1], 0, n[0], n[1], n[2], // Top right
		p[0]+o[5][0], p[1]+o[5][1], p[2]+o[5][2], br[0], br[1], 0, n[0], n[1], n[2], // Bottom right
	)
	m.count = int32(len(m.d))
	m.vboCurrent = false
}

// AddPixelFace adds a face at the given voxel position with the given facing.
// The face is scaled down to the size of one world voxel. Pos is the mesh-
// relative voxel coordinate.
func (m *CubeMesh) AddVoxelFace(pos [3]int, facing Facing, color [4]uint8) {
	s := pixelScale
	p := [3]float32{
		(float32(pos[0]-8) + 0.5) * s,
		(float32(pos[1]-8) + 0.5) * s,
		(float32(pos[2]-8) + 0.5) * s,
	}
	n := FaceNormals[facing]
	o := cubeFacingOffsets[facing]
	c := [4]float32{
		float32(color[0]) / float32(0xFF),
		float32(color[1]) / float32(0xFF),
		float32(color[2]) / float32(0xFF),
		float32(color[3]) / float32(0xFF),
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

// Data returns the raw vertex array data for the mesh.
func (m *CubeMesh) Data() []float32 {
	return m.d
}

// Reset rests the mesh builder state.
func (m *CubeMesh) Reset() {
	m.d = m.d[:0]
	m.vboCurrent = false
}

// Upload refreshes the cube mesh on the GPU.
func (m *CubeMesh) Upload() {
	// Prevent panic in unsafe code
	if len(m.d) == 0 {
		m.vboCurrent = true
		return
	}
	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.d)*4, gl.Ptr(m.d), gl.STATIC_DRAW)
	var stride int32 = 3*4 + 3*4 + 3*4
	var offset int = 0
	gl.VertexAttribPointerWithOffset(uint32(m.aPos), 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(m.aPos))
	offset += 3 * 4
	gl.VertexAttribPointerWithOffset(uint32(m.aTex), 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(m.aTex))
	offset += 3 * 4
	gl.VertexAttribPointerWithOffset(uint32(m.aNorm), 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(m.aNorm))
	offset += 3 * 4
	gl.BindVertexArray(0)
	m.vboCurrent = true
}

// draw draws the cube mesh.
func (m *CubeMesh) draw(prg *Program, o *Orientation) {
	if !m.glStarted {
		m.aPos = prg.GetAttributeLocation("pos")
		m.aTex = prg.GetAttributeLocation("texPos")
		m.aNorm = prg.GetAttributeLocation("norm")
		m.uWorld = prg.GetUniformLocation("world")
		gl.GenVertexArrays(1, &m.vao)
		gl.GenBuffers(1, &m.vbo)
		m.glStarted = true
	}
	if !m.vboCurrent {
		m.Upload()
	}
	mTransform := o.TransformMatrix()
	gl.UniformMatrix4fv(m.uWorld, 1, false, &mTransform[0])
	if m.voxel {
		gl.Uniform1i(prg.uDrawMode, drawModeColor)
	} else {
		gl.Uniform1i(prg.uDrawMode, drawModeAtlas)
		prg.atlas.bind(prg.uAtlas)
	}
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)
	gl.BindVertexArray(0)
}
