package vox

// Chunk manages a 16x16x16 dense matrix of values.
type Chunk[T comparable] struct {
	x, y, z int             // Volume bottom-north-west corner
	v       [16 * 16 * 16]T // Matrix of values
	inv     T               // Invalid value
}

// NewChunk returns a new chunk filled with the given fill value.
func NewChunk[T comparable](x, y, z int, fill, invalid T) *Chunk[T] {
	ret := &Chunk[T]{
		x:   x,
		y:   y,
		z:   z,
		inv: invalid,
	}
	ret.Fill(fill)
	return ret
}

// Fill fills the chunk with the given fill value.
func (c *Chunk[T]) Fill(f T) {
	for i := range c.v {
		c.v[i] = f
	}
}

// Get returns the value at the given location, or the invalid value if out of
// bounds.
func (c *Chunk[T]) Get(x, y, z int) T {
	return c.GetRelative(x-c.x, y-c.y, z-c.z)
}

// GetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the chunk.
func (c *Chunk[T]) GetRelative(x, y, z int) T {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return c.inv
	}
	return c.v[(z*16*16)+(y*16)+x]
}

// Set sets the value at the given location. If the location is out of bounds
// this is a no-op. True is returned if the contents of the chunk were altered.
func (c *Chunk[T]) Set(x, y, z int, f T) bool {
	return c.SetRelative(x-c.x, y-c.y, z-c.z, f)
}

// SetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the chunk.
func (c *Chunk[T]) SetRelative(x, y, z int, f T) bool {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return false
	}
	if c.v[(z*16*16)+(y*16)+x] == f {
		return false
	}
	c.v[(z*16*16)+(y*16)+x] = f
	return true
}
