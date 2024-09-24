package vox

import "github.com/go-gl/mathgl/mgl32"

// AABB models an axis-aligned bounding box.
type AABB struct {
	Min mgl32.Vec3 // Minimum extents, point of the bottom north west corner
	Max mgl32.Vec3 // Maximum extents, point of the top south east corner
}

// Contains returns true if the AABB contains point p.
func (b AABB) Contains(p mgl32.Vec3) bool {
	if p[0] < b.Min[0] || p[0] > b.Max[0] ||
		p[1] < b.Min[1] || p[1] > b.Max[1] ||
		p[2] < b.Min[2] || p[2] > b.Max[2] {
		return false
	}
	return true
}

// Lines returns pairs of points describing all of the lines of the wire frame
// for this AABB.
func (b AABB) Lines() []mgl32.Vec3 {
	return []mgl32.Vec3{
		// Bottom square
		{b.Min[0], b.Min[1], b.Min[2]},
		{b.Min[0], b.Min[1], b.Max[2]},
		{b.Min[0], b.Min[1], b.Min[2]},
		{b.Max[0], b.Min[1], b.Min[2]},
		{b.Max[0], b.Min[1], b.Max[2]},
		{b.Min[0], b.Min[1], b.Max[2]},
		{b.Max[0], b.Min[1], b.Max[2]},
		{b.Max[0], b.Min[1], b.Min[2]},
		// Top square
		{b.Min[0], b.Max[1], b.Min[2]},
		{b.Min[0], b.Max[1], b.Max[2]},
		{b.Min[0], b.Max[1], b.Min[2]},
		{b.Max[0], b.Max[1], b.Min[2]},
		{b.Max[0], b.Max[1], b.Max[2]},
		{b.Min[0], b.Max[1], b.Max[2]},
		{b.Max[0], b.Max[1], b.Max[2]},
		{b.Max[0], b.Max[1], b.Min[2]},
		// Uprights
		{b.Min[0], b.Min[1], b.Min[2]},
		{b.Min[0], b.Max[1], b.Min[2]},
		{b.Max[0], b.Min[1], b.Min[2]},
		{b.Max[0], b.Max[1], b.Min[2]},
		{b.Min[0], b.Min[1], b.Max[2]},
		{b.Min[0], b.Max[1], b.Max[2]},
		{b.Max[0], b.Min[1], b.Max[2]},
		{b.Max[0], b.Max[1], b.Max[2]},
	}
}
