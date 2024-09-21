package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// vsGlyphWidth is the width of one glyph in virtual screen units, minus boarder
// width.
const vsGlyphWidth int = VirtualScreenWidth / 80

// vsGlyphBoarder is the boarder width in virtual screen units.
const vsGlyphBoarder int = vsGlyphWidth / 2

// vsLineSpacing is the line spacing used in print commands in virtual screen
// units.
const vsLineSpacing int = (vsGlyphWidth / 2) * 3 // 1.5

// Text is a drawable layer of text. Text may be printed into the layer at any
// point, and may overlap.
type Text struct {
	f        *fontManager // Font manager in use
	d        []float32    // Vertex buffer data
	vao      uint32       // Vertex buffer array ID
	vbo      uint32       // Vertex buffer object ID
	vboDirty bool         // If true, the VBO needs to be updated
}

// newText creates a new Text with the given contents ready to use.
func newText(f *fontManager, prg *program) *Text {
	// Init
	ret := &Text{
		f: f,
	}
	gl.GenVertexArrays(1, &ret.vao)
	gl.GenBuffers(1, &ret.vbo)
	// ret.Set("Hello, OpenGL!")
	ret.Print(mgl32.Vec2{0, float32(vsGlyphWidth / 2)}, "@ABCDEFGHIJKLMNO")
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
func (t *Text) Reset() {
	t.d = t.d[:0]
	t.vboDirty = true
}

// Print prints a string at the given screen position in virtual screen units.
// Note that the baseline of the text will appear at the given Y coordinate.
func (text *Text) Print(p mgl32.Vec2, s string) {
	p[1] = (float32(VirtualScreenHeight) - p[1]) - 1 // Invert Y
	l := float32(-vsGlyphBoarder) + p[0]
	b := p[1] - float32(vsGlyphWidth)
	for _, sr := range s {
		if sr == '\n' {
			l = float32(-vsGlyphBoarder) + p[0]
			b += float32(vsLineSpacing)
			continue
		}
		t := b + float32(vsGlyphWidth+vsGlyphBoarder*2)
		r := l + float32(vsGlyphWidth+vsGlyphBoarder*2)
		g := text.f.getGlyph(sr)
		ut := g.v
		ub := ut + faAtlasStep
		ul := g.u
		ur := ul + faAtlasStep
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
	text.vboDirty = true
}

// upload uploads the VBO data.
func (t *Text) upload() {
	t.vboDirty = false
	if len(t.d) == 0 {
		return
	}
	gl.BindVertexArray(t.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(t.d)*4, gl.Ptr(t.d), gl.STATIC_DRAW)
}

func (t *Text) draw() {
	if t.vboDirty {
		t.upload()
	}
	gl.BindVertexArray(t.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(t.d)))
}
