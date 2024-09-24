package vox

import "github.com/go-gl/mathgl/mgl32"

// Ray models a ray in 3D space.
type Ray struct {
	Origin    mgl32.Vec3 // Origin of the ray
	Direction mgl32.Vec3 // Normalized direction vector of the ray
}

// IntersectsAABB returns true if the ray is within or intersects the AABB.
func (r Ray) IntersectsAABB(b AABB) bool {
	if b.Contains(r.Origin) {
		return true
	}
	return false
}
