package c3d

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/golang/freetype/truetype"
	"github.com/qbradq/cubit/data"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const faDims int = 2048
const faGlyphSize int = 32
const faGlyphDims int = 64
const faGlyphOfs int = (faGlyphDims - faGlyphSize) / 2
const faGlyphsWide int = faDims / faGlyphDims
const faAtlasStep float32 = float32(faGlyphDims) / float32(faDims)

// glyph describes the attributes of a single glyph in the atlas.
type glyph struct {
	width int     // Width of the glyph in pixels
	u, v  float32 // UV coordinates for this glyph in the texture atlas.
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
	ret.img = image.NewRGBA(image.Rect(0, 0, faDims, faDims))
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
	// TODO REMOVE DEBUG
	f, err := os.Create("font.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, ret.img)
	if err != nil {
		f.Close()
		panic(err)
	}
	f.Close()
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
		int32(faDims), int32(faDims), 0,
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
		int32(faDims), int32(faDims), 0,
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
		u:     float32(i%faGlyphsWide) * faAtlasStep,
		v:     float32(i/faGlyphsWide) * faAtlasStep,
	}
	t.glyphMap[r] = uint16(i)
	t.glyphs = append(t.glyphs, g)
	// Update the backing image
	t.d.Dot = fixed.P(
		(i%faGlyphsWide)*faGlyphDims+faGlyphOfs,
		(i/faGlyphsWide)*faGlyphDims+(faGlyphDims-faGlyphOfs),
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

// bind binds the texture to the 3D texture unit.
func (f *fontManager) bind(prg *program) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.t)
	gl.Uniform1i(prg.uni("uFont"), 0)
}
