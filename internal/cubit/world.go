package cubit

import (
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/vox"
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

// World manages the state of the entire world.
type World struct {
	sm *vox.SparseMatrix[Cell] // Matrix of all cells
}

// NewWorld returns a new World object read for use.
func NewWorld() *World {
	return &World{
		sm: vox.NewSparseMatrix(3, CellInvalid),
	}
}

// TODO DEBUG REMOVE
func (w *World) TestGen() {
	rGrass := GetCubeDef("/cubit/grass")
	rect := func(min, max [3]int, r Cell) {
		for iy := min[1]; iy <= max[1]; iy++ {
			for iz := min[2]; iz <= max[2]; iz++ {
				for ix := min[0]; ix <= max[0]; ix++ {
					w.SetCell(Pos(ix, iy, iz), r)
				}
			}
		}
	}
	rect([3]int{0, 0, 0}, [3]int{15, 0, 15}, CellForCube(rGrass, c3d.North))
}

// SetCell sets the cube and facing at the given position in the world.
func (w *World) SetCell(p Position, r Cell) {
	w.sm.Set(p.X, p.Y, p.Z, r)
}

// GetCell returns the cell value at the given position in the world.
func (w *World) GetCell(p Position) Cell {
	return w.sm.Get(p.X, p.Y, p.Z)
}
