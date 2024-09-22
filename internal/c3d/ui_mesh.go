package c3d

import (
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
	d        []float32  // Vertex buffer data
	vao      uint32     // Vertex buffer array ID
	vbo      uint32     // Vertex buffer object ID
	vboDirty bool       // If true, the VBO needs to be updated
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
	var stride int32 = 2*4 + 2*4
	var offset int = 0
	gl.BindVertexArray(ret.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, ret.vbo)
	gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexPosition")),
		2, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(prg.attr("aVertexPosition")))
	offset += 2 * 4
	gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexUV")),
		2, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(prg.attr("aVertexUV")))
	offset += 2 * 4
	return ret
}

// Reset resets the text object to blank. Internal memory buffers are retained
// to reduce allocations.
func (e *UIMesh) Reset() {
	e.d = e.d[:0]
	e.vboDirty = true
	e.Text.Reset()
}

// Tile draws the given UI tile at the position given in virtual screen units.
func (e *UIMesh) Tile(p mgl32.Vec2, i FaceIndex) {
	t := float32(VirtualScreenHeight) - p[1] // Invert Y
	b := t - float32(vsTileDims)
	l := p[0]
	r := l + float32(vsTileDims)
	utl, ubr := i.ToUV()
	e.d = append(e.d, []float32{
		//XY  U   V
		l, t, utl[0], utl[1], // TL
		r, t, ubr[0], utl[1], // TR
		l, b, utl[0], ubr[1], // BL
		l, b, utl[0], ubr[1], // BL
		r, t, ubr[0], utl[1], // TR
		r, b, ubr[0], ubr[1], // BR
	}...)
	e.vboDirty = true
}

// Scaled draws the given UI tile at the position given in virtual screen units
// with the given width and height in virtual screen units.
func (e *UIMesh) Scaled(p, d mgl32.Vec2, i FaceIndex) {
	t := float32(VirtualScreenHeight) - p[1] // Invert Y
	b := t - float32(d[1])
	l := p[0]
	r := l + float32(d[0])
	utl, ubr := i.ToUV()
	e.d = append(e.d, []float32{
		//XY  U   V
		l, t, utl[0], utl[1], // TL
		r, t, ubr[0], utl[1], // TR
		l, b, utl[0], ubr[1], // BL
		l, b, utl[0], ubr[1], // BL
		r, t, ubr[0], utl[1], // TR
		r, b, ubr[0], ubr[1], // BR
	}...)
	e.vboDirty = true
}

// NinePatch draws the specified area with the passed nine patch.
func (e *UIMesh) NinePatch(p, d mgl32.Vec2, n NinePatch) {
	// p[1] = float32(VirtualScreenHeight) - p[1] // Invert Y
	w := float32(vsTileDims)
	tl := mgl32.Vec2{p[0], p[1]}
	t := [2]mgl32.Vec2{{p[0] + w, p[1]}, {d[0] - w*2, w}}
	tr := mgl32.Vec2{p[0] + d[0] - w, p[1]}
	l := [2]mgl32.Vec2{{p[0], p[1] + w}, {w, d[1] - w*2}}
	c := [2]mgl32.Vec2{{p[0] + w, p[1] + w}, {d[0] - w*2, d[1] - w*2}}
	r := [2]mgl32.Vec2{{p[0] + d[0] - w, p[1] + w}, {w, d[1] - w*2}}
	bl := mgl32.Vec2{p[0], p[1] + d[1] - w}
	b := [2]mgl32.Vec2{{p[0] + w, p[1] + d[1] - w}, {d[0] - w*2, w}}
	br := mgl32.Vec2{p[0] + d[0] - w, p[1] + d[1] - w}
	e.Tile(tl, n[0])
	e.Tile(tr, n[2])
	e.Tile(bl, n[6])
	e.Tile(br, n[8])
	e.Scaled(t[0], t[1], n[1])
	e.Scaled(b[0], b[1], n[7])
	e.Scaled(l[0], l[1], n[3])
	e.Scaled(r[0], r[1], n[5])
	e.Scaled(c[0], c[1], n[4])
}

// upload uploads the VBO data.
func (e *UIMesh) upload() {
	if e.Text.vboDirty {
		e.Text.upload()
	}
	e.vboDirty = false
	if len(e.d) == 0 {
		return
	}
	gl.BindVertexArray(e.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, e.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(e.d)*4, gl.Ptr(e.d), gl.STATIC_DRAW)
}

func (e *UIMesh) draw() {
	if e.vboDirty {
		e.upload()
	}
	gl.BindVertexArray(e.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(e.d)))
}
