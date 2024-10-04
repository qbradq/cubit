package c3d

import (
	"encoding/json"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/t"
)

// Orientation provides data and utilities for manipulating position and Euler
// angle rotations on the X and Y axis.
type Orientation struct {
	pos   mgl32.Vec3 // Position
	pitch float32    // Pitch in radians
	yaw   float32    // Yaw in radians
	roll  float32    // Roll in radians
	quat  mgl32.Quat // Rotation quaternion
}

// Unmarshal implements the Unmarshaler interface.
func (o *Orientation) Unmarshal(b []byte) error {
	values := []float32{}
	if err := json.Unmarshal(b, &values); err != nil {
		return err
	}
	o.pos[0] = values[0]
	o.pos[1] = values[1]
	o.pos[2] = values[2]
	o.pitch = values[3]
	o.yaw = values[4]
	o.roll = values[5]
	return nil
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
func NewOrientation(pos mgl32.Vec3, pitch, yaw, roll float32) *Orientation {
	return &Orientation{
		pos:   pos,
		pitch: pitch,
		yaw:   yaw,
		roll:  roll,
		quat:  mgl32.QuatIdent(),
	}
}

// Translate adds offsets to the orientation's current position.
func (o *Orientation) Translate(t mgl32.Vec3) {
	o.pos = o.pos.Add(t)
}

// Position returns the current position of the orientation without rotation.
func (o *Orientation) Position() mgl32.Vec3 { return o.pos }

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
	return mgl32.Translate3D(o.pos[0], o.pos[1], o.pos[2])
}

// VoxelTranslationMatrix returns the translation matrix for this orientation in
// voxel units rather than block units.
func (o *Orientation) VoxelTranslationMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(
		o.pos[0]/t.VoxelScale,
		o.pos[1]/t.VoxelScale,
		o.pos[2]/t.VoxelScale,
	)
}

// TransformMatrix returns the full transform matrix for this orientation.
func (o *Orientation) TransformMatrix() mgl32.Mat4 {
	return o.TranslationMatrix().Mul4(o.RotationMatrix())
}

// VoxelTransformMatrix returns the full transform matrix for this orientation,
// using voxel units for the translation portion.
func (o *Orientation) VoxelTransformMatrix() mgl32.Mat4 {
	return o.VoxelTranslationMatrix().Mul4(o.RotationMatrix())
}

func (o *Orientation) Add(r *Orientation) Orientation {
	ret := Orientation{
		pos: o.pos.Add(r.pos),
	}
	ret.Pitch(o.pitch + r.pitch)
	ret.Yaw(o.yaw + r.yaw)
	ret.Roll(o.roll + r.roll)
	return ret
}
