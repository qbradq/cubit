package c3d

import (
	"encoding/binary"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// vsTileDims is the dimensions of a tile in virtual screen units.
const vsTileDims int = 4

// UIMesh is a drawable 2D orthographic mode object that can include ui tiles
// and text.
type UIMesh struct {
	Text     *TextMesh  // Text layer
	Position mgl32.Vec2 // Position of the element in virtual screen units
	Layer    uint16     // Layer index, the higher the value the higher priority
	d        []byte     // Vertex buffer data
	count    int32      // Vertex count
	vao      uint32     // Vertex buffer array ID
	vbo      uint32     // Vertex buffer object ID
	vboDirty bool       // If true, the VBO needs to be updated
	vbuf     [6]byte    // Vertex buffer
}

// newUIMesh creates a new UIMesh ready for use.
func newUIMesh(f *fontManager) *UIMesh {
	// Init
	ret := &UIMesh{
		Text: newTextMesh(f),
	}
	gl.CreateBuffers(1, &ret.vbo)
	gl.CreateVertexArrays(1, &ret.vao)
	gl.VertexArrayVertexBuffer(ret.vao, 0, ret.vbo, 0, 2*2+2*1)
	gl.EnableVertexArrayAttrib(ret.vao, 0)
	gl.EnableVertexArrayAttrib(ret.vao, 1)
	gl.VertexArrayAttribFormat(ret.vao, 0, 2, gl.SHORT, false, 0)
	gl.VertexArrayAttribFormat(ret.vao, 1, 2, gl.BYTE, false, 2*2)
	gl.VertexArrayAttribBinding(ret.vao, 0, 0)
	gl.VertexArrayAttribBinding(ret.vao, 1, 0)
	return ret
}

// Reset resets the text object to blank. Internal memory buffers are retained
// to reduce allocations.
func (e *UIMesh) Reset() {
	e.d = e.d[:0]
	e.count = 0
	e.vboDirty = true
	e.Text.Reset()
}

// vert packs one vertex into the mesh.
func (e *UIMesh) vert(x, y, u, v int) {
	d := e.vbuf[:]
	binary.LittleEndian.PutUint16(d[0:2], uint16(int16(x)))
	binary.LittleEndian.PutUint16(d[2:4], uint16(int16(y)))
	d[4] = byte(u)
	d[5] = byte(v)
	e.d = append(e.d, d...)
	e.count++
}

// Tile draws the given UI tile at the position given in virtual screen units.
func (e *UIMesh) Tile(x, y int, i FaceIndex) {
	t := VirtualScreenHeight - y // Invert Y
	b := t - vsTileDims
	l := x
	r := l + vsTileDims
	u, v := i.ToAtlasXY()
	//XY  U   V
	e.vert(l, t, u, v)     // TL
	e.vert(r, t, u+1, v)   // TR
	e.vert(l, b, u, v+1)   // BL
	e.vert(l, b, u, v+1)   // BL
	e.vert(r, t, u+1, v)   // TR
	e.vert(r, b, u+1, v+1) // BR
	e.vboDirty = true
}

// Scaled draws the given UI tile at the position given in virtual screen units
// with the given width and height in virtual screen units.
func (e *UIMesh) Scaled(x, y, w, h int, i FaceIndex) {
	t := VirtualScreenHeight - y // Invert Y
	b := t - h
	l := x
	r := l + w
	u, v := i.ToAtlasXY()
	//XY  U   V
	e.vert(l, t, u, v)     // TL
	e.vert(r, t, u+1, v)   // TR
	e.vert(l, b, u, v+1)   // BL
	e.vert(l, b, u, v+1)   // BL
	e.vert(r, t, u+1, v)   // TR
	e.vert(r, b, u+1, v+1) // BR
	e.vboDirty = true
}

// NinePatch draws the specified area with the passed nine patch.
func (e *UIMesh) NinePatch(x, y, w, h int, n NinePatch) {
	d := vsTileDims
	tl := [2]int{x, y}
	t := [2][2]int{{x + d, y}, {w - d*2, d}}
	tr := [2]int{x + w - d, y}
	l := [2][2]int{{x, y + d}, {d, h - d*2}}
	c := [2][2]int{{x + d, y + d}, {w - d*2, h - d*2}}
	r := [2][2]int{{x + w - d, y + d}, {d, h - d*2}}
	bl := [2]int{x, y + h - d}
	b := [2][2]int{{x + d, y + h - d}, {w - d*2, d}}
	br := [2]int{x + w - d, y + h - d}
	e.Tile(tl[0], tl[1], n[0])
	e.Tile(tr[0], tr[1], n[2])
	e.Tile(bl[0], bl[1], n[6])
	e.Tile(br[0], br[1], n[8])
	e.Scaled(t[0][0], t[0][1], t[1][0], t[1][1], n[1])
	e.Scaled(b[0][0], b[0][1], b[1][0], b[1][1], n[7])
	e.Scaled(l[0][0], l[0][1], l[1][0], l[1][1], n[3])
	e.Scaled(r[0][0], r[0][1], r[1][0], r[1][1], n[5])
	e.Scaled(c[0][0], c[0][1], c[1][0], c[1][1], n[4])
}

func (e *UIMesh) draw() {
	if e.vboDirty {
		if len(e.d) > 0 {
			gl.NamedBufferData(e.vbo, len(e.d), gl.Ptr(e.d), gl.STATIC_DRAW)
		}
		e.vboDirty = false
	}
	if e.vao != invalidVAO {
		gl.BindVertexArray(e.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, e.count)
	}
}
