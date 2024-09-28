package c3d

import (
	"encoding/binary"
	"math"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// LineMesh provides a simple axis indicator made of full-bright lines.
type LineMesh struct {
	Orientation Orientation // Orientation of the wire frame.
	d           []byte      // Raw vertex data
	count       int32       // Count of vertexes
	vao         uint32      // Vertex Array Object ID
	vbo         uint32      // Vertex Buffer Object ID
	vboDirty    bool        // If true the VBO needs to be updated on the GPU
	vbuf        [30]byte    // Double vertex line buffer
}

// NewLineMesh creates a new axis indicator.
func NewLineMesh() *LineMesh {
	ret := &LineMesh{
		vao:         invalidVAO,
		vbo:         invalidVBO,
		Orientation: *NewOrientation(mgl32.Vec3{0, 0, 0}, 0, 0, 0),
	}
	return ret
}

// Line adds a single line to the mesh.
func (m *LineMesh) Line(a, b mgl32.Vec3, c [4]uint8) {
	d := m.vbuf[:]
	binary.LittleEndian.PutUint32(d[0:4], math.Float32bits(a[0]))
	binary.LittleEndian.PutUint32(d[4:8], math.Float32bits(a[1]))
	binary.LittleEndian.PutUint32(d[8:12], math.Float32bits(a[2]))
	d[12] = c[0]
	d[13] = c[1]
	d[14] = c[2]
	binary.LittleEndian.PutUint32(d[15:19], math.Float32bits(b[0]))
	binary.LittleEndian.PutUint32(d[19:23], math.Float32bits(b[1]))
	binary.LittleEndian.PutUint32(d[23:27], math.Float32bits(b[2]))
	d[27] = c[0]
	d[28] = c[1]
	d[29] = c[2]
	m.d = append(m.d, d...)
	m.count += 2
	m.vboDirty = true
}

// WireFrame draws a wire frame from pairs of points in the slice.
func (m *LineMesh) WireFrame(v []mgl32.Vec3, c [4]uint8) {
	for i := 0; i < len(v); i += 2 {
		m.Line(v[i], v[i+1], c)
	}
}

// Reset resets the mesh to empty.
func (m *LineMesh) Reset() {
	m.d = m.d[:0]
	m.count = 0
	m.vboDirty = false
}

// draw draws the line mesh.
func (m *LineMesh) draw(prg *program) {
	if m.vao == invalidVAO {
		gl.GenVertexArrays(1, &m.vao)
	}
	if m.vbo == invalidVBO {
		var stride int32 = 3*4 + 3*1
		var offset int = 0
		gl.GenBuffers(1, &m.vbo)
		gl.BindVertexArray(m.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexPosition")),
			3, gl.FLOAT, false, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(prg.attr("aVertexPosition")))
		offset += 3 * 4
		gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexColor")),
			3, gl.UNSIGNED_BYTE, true, stride, uintptr(offset))
		gl.EnableVertexAttribArray(uint32(prg.attr("aVertexColor")))
		offset += 3 * 1
	}
	if m.vboDirty {
		if len(m.d) > 0 {
			gl.BindVertexArray(m.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
			gl.BufferData(gl.ARRAY_BUFFER, len(m.d), gl.Ptr(m.d),
				gl.STATIC_DRAW)
		}
		m.vboDirty = false
	}
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.LINES, 0, m.count)
}
