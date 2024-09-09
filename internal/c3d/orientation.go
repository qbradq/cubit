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
