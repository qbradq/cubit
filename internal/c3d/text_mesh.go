package c3d

import (
	"encoding/binary"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/qbradq/cubit/internal/t"
)

// TextMesh is a drawable layer of text. TextMesh may be printed into the layer
// at any point, and may overlap.
type TextMesh struct {
	f        *fontManager // Font manager in use
	d        []byte       // Vertex buffer data
	count    int32        // Vertex count
	vao      uint32       // Vertex buffer array ID
	vbo      uint32       // Vertex buffer object ID
	vboDirty bool         // If true, the VBO needs to be updated
	vbuf     [9]byte      // Vertex buffer
}

// newTextMesh creates a new text mesh with the given contents ready to use.
func newTextMesh(f *fontManager, prg *program) *TextMesh {
	// Init
	ret := &TextMesh{
		f: f,
	}
	gl.GenVertexArrays(1, &ret.vao)
	gl.GenBuffers(1, &ret.vbo)
	// Configure buffer attributes
	var stride int32 = 2*2 + 2*1 + 3*1
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
	gl.VertexAttribPointerWithOffset(uint32(prg.attr("aVertexColor")),
		3, gl.UNSIGNED_BYTE, true, stride, uintptr(offset))
	gl.EnableVertexAttribArray(uint32(prg.attr("aVertexColor")))
	offset += 3 * 1
	return ret
}

// Reset resets the text object to blank. Internal memory buffers are retained
// to reduce allocations.
func (t *TextMesh) Reset() {
	t.d = t.d[:0]
	t.count = 0
	t.vboDirty = true
}

// vert packs one vertex into the mesh.
func (t *TextMesh) vert(x, y, u, v int, c [3]uint8) {
	d := t.vbuf[:]
	binary.LittleEndian.PutUint16(d[0:2], uint16(int16(x)))
	binary.LittleEndian.PutUint16(d[2:4], uint16(int16(y)))
	d[4] = byte(u)
	d[5] = byte(v)
	d[6] = c[0]
	d[7] = c[1]
	d[8] = c[2]
	t.d = append(t.d, d...)
	t.count++
}

// Print prints a string at the given screen position in virtual screen units.
func (text *TextMesh) Print(x, y int, c [3]uint8, s string) {
	if len(s) == 0 {
		return
	}
	text.vboDirty = true
	y = t.VirtualScreenHeight - y // Invert Y
	l := x
	b := y - (t.VSCellHeight - t.VSBaseline)
	for _, sr := range s {
		if sr == '\n' {
			l = x
			b += t.LineSpacingVS
			continue
		}
		p := b + t.VSCellHeight
		r := l + t.VSGlyphWidth
		g := text.f.getGlyph(sr)
		ut := g.v
		ub := ut + 1
		ul := g.u
		ur := ul + 1
		//XY  U   V
		text.vert(l, p, ul, ut, c) // TL
		text.vert(r, p, ur, ut, c) // TR
		text.vert(l, b, ul, ub, c) // BL
		text.vert(l, b, ul, ub, c) // BL
		text.vert(r, p, ur, ut, c) // TR
		text.vert(r, b, ur, ub, c) // BR
		l += t.VSGlyphWidth
	}
}

func (t *TextMesh) draw() {
	if t.vboDirty {
		if len(t.d) > 0 {
			gl.BindVertexArray(t.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
			gl.BufferData(gl.ARRAY_BUFFER, len(t.d), gl.Ptr(t.d), gl.STATIC_DRAW)
		}
		t.vboDirty = false
	}
	gl.BindVertexArray(t.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(t.d)))
}
