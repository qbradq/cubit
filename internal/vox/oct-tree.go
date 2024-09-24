package vox

import "github.com/go-gl/mathgl/mgl32"

// OctalTree manages a arbitrarily sized sparse matrix of values.
type OctalTree[T any] struct {
	c      mgl32.Vec3       // Center point of the node
	r      int              // Cubic radius
	l      [8]*OctalTree[T] // Child nodes
	isLeaf bool             // If true, this is a leaf node
	v      [8]T             // Values for leaf nodes
}

// NewOctalTreeRoot returns a new OctalTree with the given dimensions. Note that
// parameter r is the cubic radius, not spherical.
func NewOctalTreeRoot[T any](r int) *OctalTree[T] {
	return &OctalTree[T]{
		r: r,
	}
}

// Get returns the value at the given location. If the space there has not
// been partitioned yet, the zero value will be returned and no space will be
// partitioned.
func (n *OctalTree[T]) Get(p mgl32.Vec3) T {
	var zero T
	idx := 0b000
	if p[0] > n.c[0] {
		idx |= 0b001
	}
	if p[1] > n.c[1] {
		idx |= 0b010
	}
	if p[2] > n.c[2] {
		idx |= 0b100
	}
	if n.isLeaf {
		return n.v[idx]
	}
	c := n.l[idx]
	if c == nil {
		return zero
	}
	return c.Get(p)
}

// Set sets the value at the given location. Note that values are Set on whole
// number cell bounds. This means that 0.0, 0.0, 0.0 and 0.9, 0.5, 0.7 both
// refer to the same location.
func (n *OctalTree[T]) Set(p mgl32.Vec3, v T) {
	idx := 0b000
	if p[0] > n.c[0] {
		idx |= 0b001
	}
	if p[1] > n.c[1] {
		idx |= 0b010
	}
	if p[2] > n.c[2] {
		idx |= 0b100
	}
	if n.isLeaf {
		n.v[idx] = v
		return
	}
	c := n.l[idx]
	if c == nil {
		sign := mgl32.Vec3{-1, -1, -1}
		if p[0] > n.c[0] {
			sign[0] = 1
		}
		if p[1] > n.c[1] {
			sign[1] = 1
		}
		if p[2] > n.c[2] {
			sign[2] = 1
		}
		hr := float32(n.r / 2)
		o := mgl32.Vec3{hr * sign[0], hr * sign[1], hr * sign[2]}
		c = &OctalTree[T]{
			c:      n.c.Add(o),
			r:      n.r / 2,
			isLeaf: n.r == 2,
		}
		n.l[idx] = c
	}
	c.Set(p, v)
}
