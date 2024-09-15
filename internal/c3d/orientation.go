package c3d

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Orientation provides data and utilities for manipulating position and Euler
// angle rotations on the X and Y axis.
type Orientation struct {
	Position mgl32.Vec3 // Position
	pitch    float32    // Pitch in radians
	yaw      float32    // Yaw in radians
	roll     float32    // Roll in radians
	quat     mgl32.Quat // Rotation quaternion
}

// OrientationZero is a global convenience pointer to empty orientation.
var OrientationZero *Orientation = &Orientation{}

// FacingToOrientation is the facing to orientation table.
var FacingToOrientation = [6]*Orientation{
	NewOrientation(mgl32.Vec3{0, 0, 0}, 0, 0, 0),
	NewOrientation(mgl32.Vec3{0, 0, 0}, 0, math.Pi, 0),
	NewOrientation(mgl32.Vec3{0, 0, 0}, 0, -math.Pi*0.5, 0),
	NewOrientation(mgl32.Vec3{0, 0, 0}, 0, math.Pi*0.5, 0),
	NewOrientation(mgl32.Vec3{0, 0, 0}, math.Pi*0.5, 0, 0),
	NewOrientation(mgl32.Vec3{0, 0, 0}, -math.Pi*0.5, 0, 0),
}

// NewOrientation creates a new orientation with the given rotations.
func NewOrientation(position mgl32.Vec3, pitch, yaw, roll float32) *Orientation {
	return &Orientation{
		Position: position,
		pitch:    pitch,
		yaw:      yaw,
		roll:     roll,
		quat:     mgl32.QuatIdent(),
	}
}

// GetPitch returns the current pitch angle.
func (o *Orientation) GetPitch() float32 { return o.pitch }

// GetYaw returns the current yaw angle.
func (o *Orientation) GetYaw() float32 { return o.yaw }

// GetRoll returns the current roll angle.
func (o *Orientation) GetRoll() float32 { return o.roll }

// Pitch rotations the orientation about the X axis by r radians.
func (o *Orientation) Pitch(r float32) {
	o.pitch = float32(math.Mod(float64(o.pitch+r), math.Pi*2))
}

// Yaw rotations the orientation about the Y axis by r radians.
func (o *Orientation) Yaw(r float32) {
	o.yaw = float32(math.Mod(float64(o.yaw+r), math.Pi*2))
}

// Roll rotations the orientation about the Z axis by r radians.
func (o *Orientation) Roll(r float32) {
	o.roll = float32(math.Mod(float64(o.roll+r), math.Pi*2))
}

// RotationMatrix returns the rotation matrix for this orientation.
func (o *Orientation) RotationMatrix() mgl32.Mat4 {
	return mgl32.AnglesToQuat(o.pitch, o.yaw, o.roll, mgl32.XYZ).Mat4()
}

// TranslationMatrix returns the translation matrix for this orientation.
func (o *Orientation) TranslationMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(o.Position[0], o.Position[1], o.Position[2])
}

// TransformMatrix returns the full transform matrix for this orientation.
func (o *Orientation) TransformMatrix() mgl32.Mat4 {
	return o.TranslationMatrix().Mul4(o.RotationMatrix())
}
