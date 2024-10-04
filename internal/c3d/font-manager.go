package c3d

import (
	"image"
	"image/draw"
	"log"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/golang/freetype/truetype"
	"github.com/qbradq/cubit/data"
	"github.com/qbradq/cubit/internal/t"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// glyph describes the attributes of a single glyph in the atlas.
type glyph struct {
	u, v, w int
}

// fontManager manages font rendering for an app.
type fontManager struct {
	t        uint32          // Texture ID in OpenGL
	img      *image.RGBA     // Image atlas backing the texture
	imgDirty bool            // If true, img has been updated since the last draw call
	f        *truetype.Font  // Font used for text rendering
	d        *font.Drawer    // Drawer used for rendering glyphs to the atlas texture
	scale    fixed.Int26_6   // Font scale
	nextIdx  uint16          // The next glyph index to be assigned
	glyphMap map[rune]uint16 // Map of runes to glyph index values
	glyphs   []glyph         // Cache of glyphs
}

// newFontManager initializes a new text structure for use.
func newFontManager(prg *program) *fontManager {
	ret := &fontManager{
		glyphMap: map[rune]uint16{},
	}
	ret.img = image.NewRGBA(image.Rect(0, 0, t.FADims, t.FADims))
	draw.Draw(ret.img, ret.img.Bounds(), image.Transparent, image.Pt(0, 0),
		draw.Src)
	// Load font
	d, err := data.FS.ReadFile("mono.ttf")
	if err != nil {
		panic(err)
	}
	ret.f, err = truetype.Parse(d)
	if err != nil {
		panic(err)
	}
	ret.scale = fixed.Int26_6(ret.f.FUnitsPerEm())
	ret.d = &font.Drawer{
		Dst: ret.img,
		Src: image.White,
		Face: truetype.NewFace(ret.f, &truetype.Options{
			Size:    48,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}
	// Initialize texture atlas backing image with the printable ASCII range
	for i := 32; i < 128; i++ {
		ret.cacheGlyph(rune(i))
	}
	// GL setup
	prg.use()
	gl.GenTextures(1, &ret.t)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ret.t)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(t.FADims), int32(t.FADims), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(ret.img.Pix),
	)
	return ret
}

// updateAtlasTexture uploads the backing image to the GPU to replace the
// current glyph atlas.
func (m *fontManager) updateAtlasTexture() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.t)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(t.FADims), int32(t.FADims), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(m.img.Pix),
	)
	m.imgDirty = false
}

// cacheGlyph adds the rune to the texture atlas and returns the index.
func (m *fontManager) cacheGlyph(r rune) glyph {
	// Load the glyph information and add it to the map
	if _, duplicate := m.glyphMap[r]; duplicate {
		log.Printf("warning: duplicate rune to c3d.text.cacheGlyph %s",
			string(r))
		return glyph{}
	}
	i := int(m.nextIdx)
	m.nextIdx++
	gb := &truetype.GlyphBuf{}
	gb.Load(m.f, m.scale, truetype.Index(r), font.HintingNone)
	g := glyph{
		w: gb.AdvanceWidth.Ceil(),
		u: i % t.FACellsWide,
		v: i / t.FACellsWide,
	}
	m.glyphMap[r] = uint16(i)
	m.glyphs = append(m.glyphs, g)
	// Update the backing image
	m.d.Dot = fixed.P(
		g.u*t.FACellWidth+t.FACellXOfs,
		g.v*t.FACellHeight+(t.FACellHeight-t.FACellYOfs),
	)
	m.d.DrawString(string(r))
	m.outlineGlyph(
		g.u*t.FACellWidth,
		g.v*t.FACellHeight,
		t.FACellWidth, t.FACellHeight,
	)
	m.imgDirty = true
	return g
}

// Renders an outline for glyph within the given bounds.
func (t *fontManager) outlineGlyph(x, y, w, h int) {
	for iy := y; iy < y+h; iy++ {
		for ix := x; ix < x+w; ix++ {
			idx := (iy * t.img.Stride) + ix*4
			a := t.img.Pix[idx+3]
			if a > 0 && a < 255 {
				t.img.Pix[idx+0] = 0
				t.img.Pix[idx+1] = 0
				t.img.Pix[idx+2] = 0
				t.img.Pix[idx+3] = 255
			}
		}
	}
}

// getGlyph returns the glyph structure for the given rune, caching it if
// needed.
func (f *fontManager) getGlyph(r rune) glyph {
	if id, found := f.glyphMap[r]; found {
		return f.glyphs[id]
	}
	return f.cacheGlyph(r)
}

// bind binds the texture to the 3D texture unit.
func (f *fontManager) bind(prg *program) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.t)
	gl.Uniform1i(prg.uni("uFont"), 0)
}
