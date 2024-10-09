package t

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// XAxis is a unit vector pointing along the positive X axis.
var XAxis mgl32.Vec3 = mgl32.Vec3{1, 0, 0}

// YAxis is a unit vector pointing along the positive Y axis.
var YAxis mgl32.Vec3 = mgl32.Vec3{0, 1, 0}

// ZAxis is a unit vector pointing along the positive Z axis.
var ZAxis mgl32.Vec3 = mgl32.Vec3{0, 0, 1}

// Orientation combines a position and a rotation quaternion and offers
// functions to transform the values.
type Orientation struct {
	P mgl32.Vec3 // Position
	Q mgl32.Quat // Rotation quaternion
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Orientation) UnmarshalJSON(d []byte) error {
	f := []float32{}
	if err := json.Unmarshal(d, &f); err != nil {
		return err
	}
	if len(f) != 6 {
		return errors.New("orientations are represented as six numbers for the X, Y, and Z displacement, then the X, Y, and Z rotations")
	}
	f[3] = mgl32.DegToRad(f[3])
	f[4] = mgl32.DegToRad(f[4])
	f[5] = mgl32.DegToRad(f[5])
	o.P = mgl32.Vec3{f[0], f[1], f[2]}
	o.Q = mgl32.AnglesToQuat(f[5], f[4], f[3], mgl32.ZYX)
	return nil
}

// FacingToOrientation is the facing to orientation table.
var FacingToOrientation = [6]Orientation{
	O(),
	O().Yaw(math.Pi),
	O().Yaw(-math.Pi * 0.5),
	O().Yaw(math.Pi * 0.5),
	O().Pitch(math.Pi * 0.5),
	O().Pitch(-math.Pi * 0.5),
}

// O returns a new orientation.
func O() Orientation {
	return Orientation{
		Q: mgl32.QuatIdent(),
	}
}

// Translate adds offsets to the orientation's current position.
func (o Orientation) Translate(t mgl32.Vec3) Orientation {
	o.P = o.P.Add(t)
	return o
}

// Pitch rotations the orientation about the X axis by r radians.
func (o Orientation) Pitch(r float32) Orientation {
	return o.Accumulate(Orientation{Q: mgl32.AnglesToQuat(0, 0, r, mgl32.ZYX)})
}

// Yaw rotations the orientation about the Y axis by r radians.
func (o Orientation) Yaw(r float32) Orientation {
	return o.Accumulate(Orientation{Q: mgl32.AnglesToQuat(0, r, 0, mgl32.ZYX)})
}

// Roll rotations the orientation about the Z axis by r radians.
func (o Orientation) Roll(r float32) Orientation {
	return o.Accumulate(Orientation{Q: mgl32.AnglesToQuat(r, 0, 0, mgl32.ZYX)})
}

// RotationMatrix returns the rotation matrix for this orientation.
func (o Orientation) RotationMatrix() mgl32.Mat4 {
	return o.Q.Mat4()
}

// TranslationMatrix returns the translation matrix for this orientation.
func (o Orientation) TranslationMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(o.P[0], o.P[1], o.P[2])
}

// TransformMatrix returns the full transform matrix for this orientation.
func (o Orientation) TransformMatrix() mgl32.Mat4 {
	return o.TranslationMatrix().Mul4(o.RotationMatrix())
}

// Accumulate accumulates the orientation and position in model space and
// returns the new orientation.
func (o Orientation) Accumulate(r Orientation) Orientation {
	// if r.Q.W == 0 {
	// 	r.Q.W = 0.00001
	// }
	// if o.Q.W == 0 {
	// 	o.Q.W = 0.00001
	// }
	o.Q = r.Q.Mul(o.Q).Normalize()
	o.P = o.Q.Rotate(o.P)
	o.P = o.P.Add(r.P)
	return o
}
