package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// vsGlyphWidth is the width of one glyph in virtual screen units, minus boarder
// width.
const vsGlyphWidth int = VirtualScreenWidth / 80

// CellDimsVS is the dimensions of a cell in screen units.
const CellDimsVS int = vsGlyphWidth

// vsCellHeight is the height of one font atlas cell in virtual screen units.
const vsCellHeight int = int(float32(vsGlyphWidth) * (float32(faCellHeight) / float32(faCellWidth)))

// vsBaseline is the Y offset for baseline in virtual screen units.
const vsBaseline int = vsCellHeight / 4

// LineSpacingVS is the line spacing used in print commands in virtual screen
// units.
const LineSpacingVS int = (vsGlyphWidth / 2) * 3 // 1.5

// TextMesh is a drawable layer of text. TextMesh may be printed into the layer
// at any point, and may overlap.
type TextMesh struct {
	f        *fontManager // Font manager in use
	d        []float32    // Vertex buffer data
	vao      uint32       // Vertex buffer array ID
	vbo      uint32       // Vertex buffer object ID
	vboDirty bool         // If true, the VBO needs to be updated
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
func (t *TextMesh) Reset() {
	t.d = t.d[:0]
	t.vboDirty = false
}

// Print prints a string at the given screen position in virtual screen units.
func (text *TextMesh) Print(p mgl32.Vec2, s string) {
	if len(s) == 0 {
		return
	}
	text.vboDirty = true
	p[1] = float32(VirtualScreenHeight) - p[1] // Invert Y
	l := p[0]
	b := p[1] - float32(vsCellHeight-vsBaseline)
	for _, sr := range s {
		if sr == '\n' {
			l = p[0]
			b += float32(LineSpacingVS)
			continue
		}
		t := b + float32(vsCellHeight)
		r := l + float32(vsGlyphWidth)
		g := text.f.getGlyph(sr)
		ut := g.v
		ub := ut + faAtlasStepV
		ul := g.u
		ur := ul + faAtlasStepU
		text.d = append(text.d, []float32{
			//XY  U   V
			l, t, ul, ut, // TL
			r, t, ur, ut, // TR
			l, b, ul, ub, // BL
			l, b, ul, ub, // BL
			r, t, ur, ut, // TR
			r, b, ur, ub, // BR
		}...)
		l += float32(vsGlyphWidth)
	}
}

func (t *TextMesh) draw() {
	if t.vboDirty {
		if len(t.d) > 0 {
			gl.BindVertexArray(t.vao)
			gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
			gl.BufferData(gl.ARRAY_BUFFER, len(t.d)*4, gl.Ptr(t.d), gl.STATIC_DRAW)
		}
		t.vboDirty = false
	}
	gl.BindVertexArray(t.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(t.d)))
}
