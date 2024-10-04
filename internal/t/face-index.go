package t

import (
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
)

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
	x = int(t) % AtlasDims
	y = int(t) / AtlasDims
	return
}

// ToUV returns the UV (texture coordinates) that should be assigned to the
// top-left and bottom-right vertexes of a face for the given FaceIndex value.
func (t FaceIndex) ToUV() (mgl32.Vec2, mgl32.Vec2) {
	x, y := t.ToAtlasXY()
	return mgl32.Vec2{
			float32(x) * PageStep,
			float32(y) * PageStep,
		},
		mgl32.Vec2{
			float32(x+1) * PageStep,
			float32(y+1) * PageStep,
		}
}

// FaceIndexFromXYZ returns the FaceIndex value for the given coordinates, where
// X and Y range 0-7 and Z ranges 0-255.
func FaceIndexFromXYZ(x, y, z int) FaceIndex {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 255 {
		return FaceIndexInvalid
	}
	return FaceIndex((x) | (y << 4) | (z << 8))
}
