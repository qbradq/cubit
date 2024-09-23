package c3d

import (
	"encoding/binary"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// vsTileDims is the dimensions of a tile in virtual screen units.
const vsTileDims int = 4

// UIMesh is a drawable 2D orthographic mode object that can include ui tiles
// and text.
type UIMesh struct {
	Text     *TextMesh  // Text layer
	Position mgl32.Vec2 // Position of the element in virtual screen units
	Layer    float32    // Layer
	d        []byte     // Vertex buffer data
	count    int32      // Vertex count
	vao      uint32     // Vertex buffer array ID
	vbo      uint32     // Vertex buffer object ID
	vboDirty bool       // If true, the VBO needs to be updated
	vbuf     [6]byte    // Vertex buffer
}

// newUIMesh creates a new UIMesh ready for use.
func newUIMesh(f *fontManager, prg *program) *UIMesh {
	// Init
	ret := &UIMesh{
		Text: newTextMesh(f, prg),
	}
	gl.GenVertexArrays(1, &ret.vao)
	gl.GenBuffers(1, &ret.vbo)
	// Configure buffer attributes
	var stride int32 = 2*2 + 2*1
	var offset int = 0
	gl.BindVertexArray(ret.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, ret.vbo)
	gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexPosition")),
		2, gl.SHORT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(prg.attr("aVertexPosition")))
	offset += 2 * 2
	gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexUV")),
		2, gl.UNSIGNED_BYTE, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(prg.attr("aVertexUV")))
	offset += 2 * 1
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
			gl.BindVertexArray(e.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, e.vbo)
			gl.BufferData(gl.ARRAY_BUFFER, len(e.d), gl.Ptr(e.d), gl.STATIC_DRAW)
		}
		e.vboDirty = false
	}
	gl.BindVertexArray(e.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, e.count)
}
