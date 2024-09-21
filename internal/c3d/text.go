package c3d

import (
	"image"
	"log"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/golang/freetype/truetype"
	"github.com/qbradq/cubit/data"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const textAtlasDims int = 1024
const taGlyphDims int = 32
const taGlyphOfs int = 8
const taGlyphsWide int = textAtlasDims / taGlyphDims

// glyph describes the attributes of a single glyph in the atlas.
type glyph struct {
	width int     // Width of the glyph in pixels
	u, v  float32 // UV coordinates for this glyph in the texture atlas.
}

// Text is a drawable string of text.
type Text struct {
	d   []float32 // Vertex buffer data
	vba uint32    // Vertex buffer array ID
	vbo uint32    // Vertex buffer object ID
}

// newText creates a new Text with the given contents ready to use.
func newText(f fontManager, t string) {
	ret := &Text{}
	gl.GenVertexArrays(1, &ret.vba)
	gl.GenBuffers(1, &ret.vbo)
}

// fontManager manages font rendering for a program.
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
func newFontManager() *fontManager {
	ret := &fontManager{
		glyphMap: map[rune]uint16{},
	}
	ret.img = image.NewRGBA(image.Rect(0, 0, textAtlasDims, textAtlasDims))
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
			Size:    24,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}
	// Initialize texture atlas backing image with the printable ASCII range
	for i := 32; i < 128; i++ {
		ret.cacheGlyph(rune(i))
	}
	// GL setup
	gl.GenTextures(1, &ret.t)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ret.t)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(textAtlasDims), int32(textAtlasDims), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(ret.img.Pix),
	)
	return ret
}

// updateAtlasTexture uploads the backing image to the GPU to replace the
// current glyph atlas.
func (t *fontManager) updateAtlasTexture() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.t)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(textAtlasDims), int32(textAtlasDims), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(t.img.Pix),
	)
	t.imgDirty = false
}

// cacheGlyph adds the rune to the texture atlas and returns the index.
func (t *fontManager) cacheGlyph(r rune) glyph {
	// Load the glyph information and add it to the map
	if _, duplicate := t.glyphMap[r]; duplicate {
		log.Printf("warning: duplicate rune to c3d.text.cacheGlyph %s",
			string(r))
		return glyph{}
	}
	i := int(t.nextIdx)
	t.nextIdx++
	gb := &truetype.GlyphBuf{}
	gb.Load(t.f, t.scale, truetype.Index(r), font.HintingNone)
	g := glyph{
		width: gb.AdvanceWidth.Ceil(),
		u:     float32(i%taGlyphsWide) / float32(taGlyphsWide),
		v:     float32(i/taGlyphsWide) / float32(taGlyphsWide),
	}
	t.glyphMap[r] = uint16(i)
	t.glyphs = append(t.glyphs, g)
	// Update the backing image
	t.d.Dot = fixed.P(
		(i%taGlyphsWide)*taGlyphDims+taGlyphOfs,
		(i/taGlyphsWide)*taGlyphDims+(taGlyphDims-taGlyphOfs),
	)
	t.d.DrawString(string(r))
	t.imgDirty = true
	return g
}

// getGlyph returns the glyph structure for the given rune, caching it if
// needed.
func (f *fontManager) getGlyph(r rune) glyph {
	if id, found := f.glyphMap[r]; found {
		return f.glyphs[id]
	}
	return f.cacheGlyph(r)
}
