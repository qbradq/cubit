package c3d

import (
	"fmt"
	"image"
	"image/draw"
	"strconv"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// atlasTextureDims are the X and Y dimensions of the face atlas texture in
// pixels.
const atlasTextureDims = 2048

// FaceDims are the X and Y dimensions of a face in pixels.
const FaceDims = 16

// atlasDims are the X and Y dimensions of the face atlas in faces.
const atlasDims = atlasTextureDims / FaceDims

// Stepping value for face U and V offsets.
const pageStep float32 = float32(1) / float32(atlasDims)

// FaceIndex indexes a tile in the face atlas. The 4 least significant bits
// encode the X face position with the source page image, and the next 4 bits
// encode the Y. The 8 most significant bits of the index encode the source
// face page.
type FaceIndex uint16

// FaceIndexInvalid is the invalid face index value.
const FaceIndexInvalid FaceIndex = 0xFFFF

func (t *FaceIndex) UnmarshalJSON(b []byte) error {
	*t = FaceIndexInvalid
	if len(b) < 1 {
		return nil
	}
	if b[0] == '"' {
		if len(b) < 3 {
			return nil
		}
		b = b[1 : len(b)-1]
	}
	v, err := strconv.ParseInt(string(b), 0, 32)
	if err != nil {
		return err
	}
	*t = FaceIndex(v)
	return nil
}

// ToXY returns the X, Y, and Z coordinates of the tile in tile source sets.
func (t FaceIndex) ToXYZ() (x int, y int, z int) {
	x = int(t & 0x000F)
	y = int(t&0x00F0) >> 4
	z = int(t&0xFF00) >> 8
	return
}

// ToXY returns the X and Y coordinates of the tile in atlas texture.
func (t FaceIndex) ToAtlasXY() (x int, y int) {
	x = int(t) % atlasDims
	y = int(t) / atlasDims
	return
}

// ToUV returns the UV (texture coordinates) that should be assigned to the
// top-left and bottom-right vertexes of a face for the given FaceIndex value.
func (t FaceIndex) ToUV() (mgl32.Vec2, mgl32.Vec2) {
	x, y := t.ToAtlasXY()
	return mgl32.Vec2{
			float32(x) * pageStep,
			float32(y) * pageStep,
		},
		mgl32.Vec2{
			float32(x+1) * pageStep,
			float32(y+1) * pageStep,
		}
}

// toCompressedUV returns the UV (texture coordinates) that should be assigned
// to the top-left and bottom-right vertexes of a face for the given FaceIndex
// value, represented as the compressed face UV.
func (t FaceIndex) toCompressedUV() ([2]byte, [2]byte) {
	x, y := t.ToAtlasXY()
	return [2]byte{byte(x), byte(y)}, [2]byte{byte(x + 1), byte(y + 1)}
}

// FaceIndexFromXYZ returns the FaceIndex value for the given coordinates, where
// X and Y range 0-7 and Z ranges 0-255.
func FaceIndexFromXYZ(x, y, z int) FaceIndex {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 255 {
		return FaceIndexInvalid
	}
	return FaceIndex((x) | (y << 4) | (z << 8))
}

// FaceAtlas manages
type FaceAtlas struct {
	textureID uint32      // Texture ID in OpenGL for the atlas
	nextIndex FaceIndex   // The next FaceIndex value to be assigned
	img       *image.RGBA // The atlas texture image while building
}

// NewFaceAtlas creates a new FaceAtlas object ready for use.
func NewFaceAtlas() *FaceAtlas {
	ret := &FaceAtlas{
		img: image.NewRGBA(image.Rect(0, 0, atlasTextureDims,
			atlasTextureDims)),
	}
	return ret
}

// freeMemory frees internal memory pages used during image rendering. This is
// not needed after a call to upload().
func (a *FaceAtlas) freeMemory() {
	a.img = nil
}

// AddFace adds a single face graphic to the atlas and returns the FaceIndex.
func (a *FaceAtlas) AddFace(img *image.RGBA) FaceIndex {
	if img.Bounds().Max.X != FaceDims || img.Bounds().Max.Y != FaceDims {
		panic(fmt.Errorf("all face images must be %dx%d pixels", FaceDims,
			FaceDims))
	}
	x, y := a.nextIndex.ToAtlasXY()
	draw.Draw(
		a.img,
		image.Rect(x*FaceDims, y*FaceDims, (x+1)*FaceDims, (y+1)*FaceDims),
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
		atlasTextureDims, atlasTextureDims, 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(a.img.Pix))
}

// bind binds the texture to the 3D texture unit.
func (a *FaceAtlas) bind(prg *program) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, a.textureID)
	gl.Uniform1i(prg.uni("uAtlas"), 0)
}
