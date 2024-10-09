package t

import (
	"encoding/json"
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

// AABB represents an axis-aligned bounding box.
type AABB [2]mgl32.Vec3

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *AABB) UnmarshalJSON(d []byte) error {
	f := []float32{}
	if err := json.Unmarshal(d, &f); err != nil {
		return err
	}
	if len(f) != 6 {
		return errors.New("bounding boxes are represented as an array of six numbers representing the X, Y, and Z coordinates of the North-West-Bottom and South-East-Top corners respectively")
	}
	*b = [2]mgl32.Vec3{
		mgl32.Vec3{f[0], f[1], f[2]}.Mul(VoxelScale),
		mgl32.Vec3{f[3], f[4], f[5]}.Mul(VoxelScale),
	}
	return nil
}
