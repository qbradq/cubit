package t

// Chunk manages a 16x16x16 dense matrix of 32-bit values.
type Chunk struct {
	Position    IVec3  // World coordinates of the bottom-north-west corner of the chunk
	Revision    uint32 // Chunk revision number, increments with every change to cell contents
	VoxRevision uint32 // Chunk vox model revision number, increments with every insertion or removal of a vox model
	cells       []Cell // Matrix of values
	isSolid     bool   // If true the chunk is solid and cells is nil
	solid       Cell   // Solid cell value
}

// newChunk returns a new chunk filled with the given fill value.
func newChunk(p IVec3, fill Cell) *Chunk {
	ret := &Chunk{
		Position: p,
		Revision: 1,
		isSolid:  true,
		solid:    fill,
	}
	return ret
}

// Fill fills the chunk with the given fill value.
func (c *Chunk) Fill(f Cell) {
	if c.isSolid && c.solid == f {
		return
	}
	c.cells = nil
	c.isSolid = true
	c.solid = f
	c.Revision++
}

// Get implements the c3d.VoxelSource interface.
func (c *Chunk) Get(x, y, z int) Cell {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return CellInvalid
	}
	return c.cells[(z*16*16)+(y*16)+x]
}

// Dimensions implements the c3d.VoxelSource interface.
func (c *Chunk) Dimensions() (w, h, d int) {
	return 16, 16, 16
}

// IsEmpty implements the c3d.VoxelSource interface.
func (c *Chunk) IsEmpty(v Cell) bool {
	n, _, _ := v.Decompose()
	return n == CubeRefInvalid
}

// GetCell returns the value at the given location, or the invalid value if out of
// bounds.
func (c *Chunk) GetCell(p IVec3) Cell {
	return c.GetRelative(p.Sub(c.Position))
}

// GetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the chunk.
func (c *Chunk) GetRelative(p IVec3) Cell {
	if p[0] < 0 || p[0] > 15 || p[1] < 0 || p[1] > 15 || p[2] < 0 || p[2] > 15 {
		return CellInvalid
	}
	return c.cells[(p[2]*16*16)+(p[1]*16)+p[0]]
}

// SetCell sets the value at the given location. If the location is out of bounds
// this is a no-op. True is returned if the contents of the chunk were altered.
func (c *Chunk) SetCell(p IVec3, v Cell) bool {
	return c.SetRelative(p.Sub(c.Position), v)
}

// SetRelative is like Get(), but the location ip.Zs relative to the bottom-north-
// west corner of the chunk.
func (c *Chunk) SetRelative(p IVec3, v Cell) bool {
	if p[0] < 0 || p[0] > 15 || p[1] < 0 || p[1] > 15 || p[2] < 0 || p[2] > 15 {
		return false
	}
	if c.isSolid {
		if c.solid == v {
			return false
		}
		c.cells = make([]Cell, 16*16*16)
		for i := range c.cells {
			c.cells[i] = c.solid
		}
		c.isSolid = false
	}
	idx := p[2]*16*16 + p[1]*16 + p[0]
	ov := c.cells[idx]
	if ov == v {
		return false
	}
	c.cells[idx] = v
	c.Revision++
	if ov.IsVox() || v.IsVox() {
		c.VoxRevision++
	}
	return true
}
