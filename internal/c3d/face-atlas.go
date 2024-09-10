package c3d

import (
	"image"
	"image/draw"
	"strconv"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// Stepping value for face S offset.
const pageSStep float32 = float32(1) / float32(64)

// Stepping value for face T offset.
const pageTStep float32 = float32(1) / float32(64)

// FaceIndex indexes a tile in the face atlas. The 10 most significant bits
// encode the atlas page containing the tile. The next 3 most significant bits
// encode the Y position on the texture atlas in tiles. The file 3 least
// significant bits encode the X position on the texture atlas in tiles.
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
	x = int(t >> 6)
	y = int((t >> 3) & 0x7)
	z = int(t & 0x07)
	return
}

// ToXY returns the X and Y coordinates of the tile in atlas texture.
func (t FaceIndex) ToAtlasXY() (x int, y int) {
	x = int(t % 64)
	y = int(t / 64)
	return
}

// ToST returns the ST (texture coordinates) that should be assigned to the
// top-left and bottom-right vertexes of a face for the given FaceIndex value.
func (t FaceIndex) ToST() (mgl32.Vec3, mgl32.Vec3) {
	x, y := t.ToAtlasXY()
	return mgl32.Vec3{
			float32(x) * pageSStep,
			float32(y) * pageTStep,
			0,
		},
		mgl32.Vec3{
			float32(x+1) * pageSStep,
			float32(y+1) * pageTStep,
			0,
		}
}

// FaceIndexFromXYZ returns the FaceIndex value for the given coordinates, where
// X and Y range 0-7 and Z ranges 0-255.
func FaceIndexFromXYZ(x, y, z int) FaceIndex {
	if x < 0 || x > 7 || y < 0 || y > 7 || z < 0 || z > 255 {
		return FaceIndexInvalid
	}
	return FaceIndex((x) | (y << 3) | (z << 6))
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
		img: image.NewRGBA(image.Rect(0, 0, 2048, 2048)),
	}
	for i := range ret.img.Pix {
		ret.img.Pix[i] = 0xFF
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
	if img.Bounds().Max.X != 32 || img.Bounds().Max.Y != 32 {
		panic("all face images must be 32x32 pixels")
	}
	x, y := a.nextIndex.ToAtlasXY()
	draw.Draw(
		a.img,
		image.Rect(x*32, y*32, (x+1)*32, (y+1)*32),
		img,
		image.Pt(0, 0),
		draw.Src,
	)
	ret := a.nextIndex
	a.nextIndex++
	return ret
}

// upload uploads the entire face atlas to the GPU as a 2D texture array.
func (a *FaceAtlas) upload() {
	gl.GenTextures(1, &a.textureID)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, a.textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		2048, 2048, 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(a.img.Pix))
}

// bind binds the texture to the 3D texture unit.
func (a *FaceAtlas) bind(u int32) {
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, a.textureID)
	gl.Uniform1i(u, 1)
}

// unbind un-binds the texture from the 3D texture unit.
func (a *FaceAtlas) unbind() {
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
