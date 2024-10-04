package t

// Cell encodes the cube or vox contained in a single Cell along with the
// orientation.
type Cell uint32

// CellInvalid is the invalid valid for cells.
const CellInvalid Cell = 0xFFFFFFFF

// CellForCube returns the cell value with the given values encoded.
func CellForCube(r CubeRef, f Facing) Cell {
	return Cell(r) | (Cell(f) << 16)
}

// CellForVox returns the cell value with the given values encoded.
func CellForVox(r VoxRef, f Facing) Cell {
	return Cell(r) | (Cell(f) << 16) | 0x80000000
}

// IsCube returns true if this is a cube reference.
func (l Cell) IsCube() bool {
	if l == CellInvalid {
		return false
	}
	if l&0x80000000 != 0 {
		return false
	}
	c := CubeRef(l & 0xFFFF)
	return c != CubeRefInvalid
}

// IsVox returns true if this is a voxel model reference.
func (l Cell) IsVox() bool {
	if l == CellInvalid {
		return false
	}
	if l&0x80000000 == 0 {
		return false
	}
	c := VoxRef(l & 0xFFFF)
	return c != VoxRefInvalid
}

// Decompose returns the portions of the cell value broken out.
func (l Cell) Decompose() (c CubeRef, v VoxRef, f Facing) {
	f = Facing((l >> 16) & 0x7)
	if f > Bottom {
		f = 0
	}
	c = CubeRefInvalid
	v = VoxRefInvalid
	if l&0x80000000 == 0 {
		c = CubeRef(l & 0xFFFF)
	} else {
		v = VoxRef(l & 0xFFFF)
	}
	return c, v, f
}
