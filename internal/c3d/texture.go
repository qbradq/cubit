package c3d

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

// Texture represents a single texture.
type Texture struct {
	id     uint32 // Texture ID in OpenGL
	Width  int    // Width of the texture in pixels
	Height int    // Height of the texture in pixels
}

// NewTexture loads a new texture from the provided data.
func NewTexture(data []byte) *Texture {
	// Decode texture image data
	ret := &Texture{}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	ret.Width = img.Bounds().Max.X
	ret.Height = img.Bounds().Max.Y
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	// Set texture options
	gl.GenTextures(1, &ret.id)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ret.id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(ret.Width), int32(ret.Height), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix),
	)
	return ret
}

// bind binds the texture to the 2D texture unit.
func (t *Texture) bind(u int32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
	gl.Uniform1i(u, 0)
}

// unbind un-binds the texture from the 2D texture unit.
func (t *Texture) unbind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
