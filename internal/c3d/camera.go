package c3d

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Camera manages the properties and matrixes of a camera.
type Camera struct {
	Position mgl32.Vec3 // Position of the camera
	Front    mgl32.Vec3 // Camera facing direction
	Up       mgl32.Vec3 // Camera up vector
	Pitch    float32    // Camera rotation about the X axis
	Yaw      float32    // Camera rotation about the Y axis
}

// NewCamera creates a new Camera object ready for use.
func NewCamera(pos mgl32.Vec3) *Camera {
	return &Camera{
		Position: pos,
		Front:    mgl32.Vec3{0, 0, -1},
		Up:       mgl32.Vec3{0, 1, 0},
		Yaw:      -90,
	}
}

// TransformMatrix returns the matrix to transform from world space to camera
// space.
func (c *Camera) TransformMatrix() mgl32.Mat4 {
	var dir mgl32.Vec3
	dir[0] = float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch))))
	dir[1] = float32(math.Sin(float64(mgl32.DegToRad(c.Pitch))))
	dir[2] = float32(math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch))))
	c.Front = dir.Normalize()
	target := c.Position.Add(c.Front)
	return mgl32.LookAt(
		c.Position[0], c.Position[1], c.Position[2],
		target[0], target[1], target[2],
		c.Up[0], c.Up[1], c.Up[2],
	)
}
