package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// UIElement is a drawable 2D orthographic mode object that can include ui tiles
// and text.
type UIElement struct {
	Text     *Text     // Text layer
	d        []float32 // Vertex buffer data
	vao      uint32    // Vertex buffer array ID
	vbo      uint32    // Vertex buffer object ID
	vboDirty bool      // If true, the VBO needs to be updated
}

// newUIElement creates a new UIElement ready for use.
func newUIElement(f *fontManager, prg *program) *UIElement {
	// Init
	ret := &UIElement{
		Text: &Text{
			f: f,
		},
	}
	gl.GenVertexArrays(1, &ret.vao)
	gl.GenBuffers(1, &ret.vbo)
	// ret.Set("Hello, OpenGL!")
	ret.Text.Print(mgl32.Vec2{100, 50}, "Hello, Orth2D!")
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
func (e *UIElement) Reset() {
	e.d = e.d[:0]
	e.vboDirty = true
	e.Text.Reset()
}

// // Tile draws the given UI tile at the position given in virtual screen units.
// func (e *UIElement) Tile(p mgl32.Vec2, t FaceIndex) {
// 	p[1] = (float32(VirtualScreenHeight) - p[1]) - 1 // Invert Y
// 	l := float32(-vsGlyphBoarder) + p[0]
// 	b := float32(-vsGlyphBoarder) + p[1]
// 	for _, sr := range s {
// 		if sr == '\n' {
// 			l = float32(-vsGlyphBoarder) + p[0]
// 			b += float32(vsLineSpacing)
// 			continue
// 		}
// 		t := b + float32(vsGlyphWidth+vsGlyphBoarder*2)
// 		r := l + float32(vsGlyphWidth+vsGlyphBoarder*2)
// 		g := text.f.getGlyph(sr)
// 		ut := g.v
// 		ub := ut + faAtlasStep
// 		ul := g.u
// 		ur := ul + faAtlasStep
// 		text.d = append(text.d, []float32{
// 			//XY  U   V
// 			l, t, ul, ut, // TL
// 			r, t, ur, ut, // TR
// 			l, b, ul, ub, // BL
// 			l, b, ul, ub, // BL
// 			r, t, ur, ut, // TR
// 			r, b, ur, ub, // BR
// 		}...)
// 		l += float32(vsGlyphWidth)
// 	}
// 	text.vboDirty = true
// }

// // upload uploads the VBO data.
// func (t *Text) upload() {
// 	t.vboDirty = false
// 	if len(t.d) == 0 {
// 		return
// 	}
// 	gl.BindVertexArray(t.vao)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(t.d)*4, gl.Ptr(t.d), gl.STATIC_DRAW)
// }

// func (t *Text) draw() {
// 	if t.vboDirty {
// 		t.upload()
// 	}
// 	gl.BindVertexArray(t.vao)
// 	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(t.d)))
// }
