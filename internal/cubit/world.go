package cubit

import (
	"github.com/qbradq/cubit/internal/c3d"
)

// Cell encodes the cube or vox contained in a single Cell along with the
// orientation.
type Cell uint32

// CellInvalid is the invalid valid for cells.
const CellInvalid Cell = 0xFFFFFFFF

// CellForCube returns the cell value with the given values encoded.
func CellForCube(r CubeRef, f c3d.Facing) Cell {
	return Cell(r) | (Cell(f) << 16)
}

// CellForVox returns the cell value with the given values encoded.
func CellForVox(r VoxRef, f c3d.Facing) Cell {
	return Cell(r) | (Cell(f) << 16) | 0x80000000
}

// Decompose returns the portions of the cell value broken out.
func (l Cell) Decompose() (c *Cube, v *Vox, f c3d.Facing) {
	f = c3d.Facing((l >> 16) & 0x7)
	if f > c3d.Bottom {
		f = 0
	}
	c = nil
	v = nil
	if l&0x80000000 == 0 {
		ref := CubeRef(l & 0xFFFF)
		if ref != CubeRefInvalid {
			c = cubeDefs[ref]
		}
	} else {
		ref := VoxRef(l & 0xFFFF)
		if ref != VoxRefInvalid {
			v = voxDefs[ref]
		}
	}
	return
}

// ChunkRef references a single chunk within the world.
type ChunkRef uint32

// InvalidChunkRef is the invalid value for ChunkRef.
const InvalidChunkRef ChunkRef = 0xFFFFFFFF

// NewChunkRef creates a new chunk reference with the given chunk coordinates.
func NewChunkRef(p Position) ChunkRef {
	x := p.X + (2^14)/2
	y := p.Y + (2^14)/2
	z := p.Z + (2^4)/2
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
func NewChunkRefForWorldPosition(p Position) ChunkRef {
	return NewChunkRef(p.Div(Pos(16, 16, 16)))
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

// TODO DEBUG REMOVE
func (w *World) TestGen() {
	rStone := GetCubeDef("/cubit/stone")
	rGrass := GetCubeDef("/cubit/grass")
	vWindow := GetVoxByPath("/cubit/window0")
	rect := func(min, max Position, r Cell) {
		for iy := min.Y; iy <= max.Y; iy++ {
			for iz := min.Z; iz <= max.Z; iz++ {
				for ix := min.X; ix <= max.X; ix++ {
					w.SetCell(Pos(ix, iy, iz), r)
				}
			}
		}
	}
	// Ground
	rect(Pos(0, 0, 0), Pos(15, 0, 15), CellForCube(rGrass, c3d.North))
	// Walls
	rect(Pos(4, 1, 4), Pos(4, 3, 10), CellForCube(rStone, c3d.North))
	rect(Pos(10, 1, 4), Pos(10, 3, 10), CellForCube(rStone, c3d.North))
	rect(Pos(4, 1, 4), Pos(10, 3, 4), CellForCube(rStone, c3d.North))
	rect(Pos(4, 1, 10), Pos(10, 3, 10), CellForCube(rStone, c3d.North))
	// Ceiling
	rect(Pos(4, 4, 4), Pos(10, 4, 10), CellForCube(rStone, c3d.North))
	// Window and doorway
	w.SetCell(Pos(6, 1, 10), CellInvalid)
	w.SetCell(Pos(6, 2, 10), CellInvalid)
	w.SetCell(Pos(8, 2, 10), CellForVox(vWindow.Ref, c3d.North))
}

// SetCell sets the cube and facing at the given position in the world. Returns
// true if the voxel was changed.
func (w *World) SetCell(p Position, v Cell) bool {
	cp := p.Div(Pos(16, 16, 16))
	cr := NewChunkRef(cp)
	c := w.chunks[cr]
	if c == nil {
		c = NewChunk(p, CellInvalid)
		w.chunks[cr] = c
	}
	return c.Set(p, v)
}

// GetCell returns the cell value at the given position in the world.
func (w *World) GetCell(p Position) Cell {
	cp := p.Div(Pos(16, 16, 16))
	cr := NewChunkRef(cp)
	c := w.chunks[cr]
	if c == nil {
		return CellInvalid
	}
	return c.Get(p)
}
