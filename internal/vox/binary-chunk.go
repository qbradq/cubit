package vox

// BinaryChunk holds a dense 16x16x16 matrix of boolean values. The zero value
// is valid and is a matrix of all false values.
type BinaryChunk [4*16 + 1]uint64

// Get returns the value at the given location, or false if out of bounds.
func (c *BinaryChunk) Get(x, y, z int) bool {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return false
	}
	iBit := y*16*16 + z*16 + x
	return c[iBit>>6]>>(iBit&0x3F) != 0
}

// Set sets the value at the given location. This is a no-op if out of bounds.
// If v matches the value currently in the matrix, this is also a no-op.
func (c *BinaryChunk) Set(x, y, z int, v bool) {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return
	}
	iBit := y*16*16 + z*16 + x
	cv := c[iBit>>6]>>(iBit&0x3F) != 0
	if cv == v {
		return
	}
	if v {
		c[iBit>>6] |= 1 << (iBit & 0x3F)
		c[4*16]++
	} else {
		c[iBit>>6] &= ^(1 << (iBit & 0x3F))
		c[4*16]--
	}
}

// IsEmpty returns true if the chunk contains no true values.
func (c *BinaryChunk) IsEmpty() bool {
	return c[4*16] == 0
}
