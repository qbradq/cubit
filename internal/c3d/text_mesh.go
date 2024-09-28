package c3d

import (
	"encoding/binary"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// VirtualScreenGlyphsWide is the width of the screen in glyphs.
const VirtualScreenGlyphsWide int = 80

// vsGlyphWidth is the width of one glyph in virtual screen units, minus boarder
// width.
const vsGlyphWidth int = VirtualScreenWidth / VirtualScreenGlyphsWide

// CellDimsVS is the dimensions of a cell in screen units.
const CellDimsVS int = vsGlyphWidth

// vsCellHeight is the height of one font atlas cell in virtual screen units.
const vsCellHeight int = (vsGlyphWidth * faCellHeight) / faCellWidth

// vsBaseline is the Y offset for baseline in virtual screen units.
const vsBaseline int = vsCellHeight / 4

// LineSpacingVS is the line spacing used in print commands in virtual screen
// units.
const LineSpacingVS int = (vsGlyphWidth / 2) * 3 // 1.5

// TextMesh is a drawable layer of text. TextMesh may be printed into the layer
// at any point, and may overlap.
type TextMesh struct {
	f        *fontManager // Font manager in use
	d        []byte       // Vertex buffer data
	count    int32        // Vertex count
	vao      uint32       // Vertex buffer array ID
	vbo      uint32       // Vertex buffer object ID
	vboDirty bool         // If true, the VBO needs to be updated
	vbuf     [6]byte      // Vertex buffer
}

// newTextMesh creates a new text mesh with the given contents ready to use.
func newTextMesh(f *fontManager) *TextMesh {
	// Init
	ret := &TextMesh{
		f: f,
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
func (t *TextMesh) Reset() {
	t.d = t.d[:0]
	t.count = 0
	t.vboDirty = true
}

// vert packs one vertex into the mesh.
func (t *TextMesh) vert(x, y, u, v int) {
	d := t.vbuf[:]
	binary.LittleEndian.PutUint16(d[0:2], uint16(int16(x)))
	binary.LittleEndian.PutUint16(d[2:4], uint16(int16(y)))
	d[4] = byte(u)
	d[5] = byte(v)
	t.d = append(t.d, d...)
	t.count++
}

// Print prints a string at the given screen position in virtual screen units.
func (text *TextMesh) Print(x, y int, s string) {
	if len(s) == 0 {
		return
	}
	text.vboDirty = true
	y = VirtualScreenHeight - y // Invert Y
	l := x
	b := y - (vsCellHeight - vsBaseline)
	for _, sr := range s {
		if sr == '\n' {
			l = x
			b += LineSpacingVS
			continue
		}
		t := b + vsCellHeight
		r := l + vsGlyphWidth
		g := text.f.getGlyph(sr)
		ut := g.v
		ub := ut + 1
		ul := g.u
		ur := ul + 1
		//XY  U   V
		text.vert(l, t, ul, ut) // TL
		text.vert(r, t, ur, ut) // TR
		text.vert(l, b, ul, ub) // BL
		text.vert(l, b, ul, ub) // BL
		text.vert(r, t, ur, ut) // TR
		text.vert(r, b, ur, ub) // BR
		l += vsGlyphWidth
	}
}

func (t *TextMesh) draw() {
	if t.vboDirty {
		if len(t.d) > 0 {
			gl.NamedBufferData(t.vbo, len(t.d), gl.Ptr(t.d), gl.STATIC_DRAW)
		}
		t.vboDirty = false
	}
	gl.BindVertexArray(t.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(t.d)))
}
