package cubit

// Chunk manages a 16x16x16 dense matrix of 32-bit values.
type Chunk struct {
	Position Position // World coordinates of the bottom-north-west corner of the chunk
	Revision uint32   // Chunk revision number, increments with every change to chunk contents
	cells    []Cell   // Matrix of values
	isSolid  bool     // If true the chunk is solid and cells is nil
	solid    Cell     // Solid cell value
}

// newChunk returns a new chunk filled with the given fill value.
func newChunk(p Position, fill Cell) *Chunk {
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

// Get returns the value at the given location, or the invalid value if out of
// bounds.
func (c *Chunk) Get(p Position) Cell {
	return c.GetRelative(p.Sub(c.Position))
}

// GetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the chunk.
func (c *Chunk) GetRelative(p Position) Cell {
	if p.X < 0 || p.X > 15 || p.Y < 0 || p.Y > 15 || p.Z < 0 || p.Z > 15 {
		return CellInvalid
	}
	return c.cells[(p.Z*16*16)+(p.Y*16)+p.X]
}

// Set sets the value at the given location. If the location is out of bounds
// this is a no-op. True is returned if the contents of the chunk were altered.
func (c *Chunk) Set(p Position, v Cell) bool {
	return c.SetRelative(p.Sub(c.Position), v)
}

// SetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the chunk.
func (c *Chunk) SetRelative(p Position, v Cell) bool {
	if p.X < 0 || p.X > 15 || p.Y < 0 || p.Y > 15 || p.Z < 0 || p.Z > 15 {
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
	idx := p.Z*16*16 + p.Y*16 + p.X
	if c.cells[idx] == v {
		return false
	}
	c.cells[idx] = v
	c.Revision++
	return true
}
