package t

// ChunkRef references a single chunk within the world.
type ChunkRef uint32

// InvalidChunkRef is the invalid value for ChunkRef.
const InvalidChunkRef ChunkRef = 0xFFFFFFFF

// NewChunkRef creates a new chunk reference with the given chunk coordinates.
func NewChunkRef(p IVec3) ChunkRef {
	x := p[0] + (2^14)/2
	y := p[1] + (2^14)/2
	z := p[2] + (2^4)/2
	if x < 0 || x >= 2^14 || y < 0 || y >= 2^14 || z < 0 || z >= 2^4 {
		return InvalidChunkRef
	}
	return ChunkRef(
		(x & 0b00000000000000000011111111111111) |
			((y & 0b00000000000000000011111111111111) << 14) |
			((y & 0b00000000000000000000000000001111) << 28),
	)
}

// NewChunkRefForWorldPosition creates a new chunk reference with the chunk that
// contains the given world coordinates.
func NewChunkRefForWorldPosition(p IVec3) ChunkRef {
	return NewChunkRef(p.Div(IVec3{16, 16, 16}))
}

// World manages the state of the entire world.
type World struct {
	chunks map[ChunkRef]*Chunk
}

// NewWorld returns a new World object read for use.
func NewWorld() *World {
	return &World{
		chunks: map[ChunkRef]*Chunk{},
	}
}

// SetCell sets the cube and facing at the given position in the world. Returns
// true if the voxel was changed.
func (w *World) SetCell(p IVec3, v Cell) bool {
	cp := p.Div(IVec3{16, 16, 16})
	cr := NewChunkRef(cp)
	c := w.chunks[cr]
	if c == nil {
		c = newChunk(p, CellInvalid)
		w.chunks[cr] = c
	}
	return c.SetCell(p, v)
}

// GetCell returns the cell value at the given position in the world.
func (w *World) GetCell(p IVec3) Cell {
	cp := p.Div(IVec3{16, 16, 16})
	cr := NewChunkRef(cp)
	c := w.chunks[cr]
	if c == nil {
		return CellInvalid
	}
	return c.GetCell(p)
}

// GetChunkByRef returns the chunk for the given chunk reference, or nil if no
// chunk exists.
func (w *World) GetChunkByRef(r ChunkRef) *Chunk {
	return w.chunks[r]
}
