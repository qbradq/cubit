package c3d

import (
	"fmt"
	"image"
	"image/draw"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/qbradq/cubit/internal/t"
)

// FaceAtlas manages
type FaceAtlas struct {
	textureID uint32      // Texture ID in OpenGL for the atlas
	nextIndex t.FaceIndex // The next FaceIndex value to be assigned
	img       *image.RGBA // The atlas texture image while building
}

// NewFaceAtlas creates a new FaceAtlas object ready for use.
func NewFaceAtlas() *FaceAtlas {
	ret := &FaceAtlas{
		img: image.NewRGBA(image.Rect(0, 0, t.AtlasTextureDims,
			t.AtlasTextureDims)),
	}
	return ret
}

// freeMemory frees internal memory pages used during image rendering. This is
// not needed after a call to upload().
func (a *FaceAtlas) freeMemory() {
	a.img = nil
}

// AddFace adds a single face graphic to the atlas and returns the FaceIndex.
func (a *FaceAtlas) AddFace(img *image.RGBA) t.FaceIndex {
	if img.Bounds().Max.X != t.FaceDims || img.Bounds().Max.Y != t.FaceDims {
		panic(fmt.Errorf("all face images must be %dx%d pixels", t.FaceDims,
			t.FaceDims))
	}
	x, y := a.nextIndex.ToAtlasXY()
	draw.Draw(
		a.img,
		image.Rect(x*t.FaceDims, y*t.FaceDims, (x+1)*t.FaceDims,
			(y+1)*t.FaceDims),
		img,
		image.Pt(0, 0),
		draw.Src,
	)
	ret := a.nextIndex
	a.nextIndex++
	return ret
}

// upload uploads the entire face atlas to the GPU as a 2D texture array.
func (a *FaceAtlas) upload(prg *program) {
	prg.use()
	gl.GenTextures(1, &a.textureID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, a.textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		t.AtlasTextureDims, t.AtlasTextureDims, 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(a.img.Pix))
}

// bind binds the texture to the 3D texture unit.
func (a *FaceAtlas) bind(prg *program) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, a.textureID)
	gl.Uniform1i(prg.uni("uAtlas"), 0)
}
