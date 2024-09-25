package vox

import "math"

// SparseMatrix uses binary chunks to partition a sparse matrix of
// values.
type SparseMatrix[T comparable] struct {
	pos     [3]int // Position of the bottom-north-west corner of the chunk in global coordinate space
	level   int    // Level, if 0 this means the chunk contains leaf nodes rather than further BCPs
	zero    T      // Zero value for Get()
	solid   T      // if isSolid is true, this is the value the chunk is filled with
	isSolid bool
	bcp     []*SparseMatrix[T] // Matrix of sub-chunks, only valid if level > 0
	leaves  []T                // Matrix of values, only valid if level == 0
}

// NewSparseMatrix creates a new sparse matrix. The level parameter indicates
// the number of levels. The 0x0x0 point is at the bottom-north-west corner of
// the space, and coordinates extend in the positive direction only.
// Level	Matrix Modeled
// 0		16^3
// 1		256^3
// 2		4096^3
// 3		65536^3
// n		16^n+1^3
func NewSparseMatrix[T comparable](level int, zero T) *SparseMatrix[T] {
	ret := &SparseMatrix[T]{
		level:   level,
		zero:    zero,
		solid:   zero,
		isSolid: true,
	}
	return ret
}

// IsEmpty returns true if the node contains no values.
func (p *SparseMatrix[T]) IsEmpty() bool {
	if p.level == 0 {
		return p.leaves == nil || (p.isSolid && p.solid == p.zero)
	}
	return p.bcp == nil
}

// Get returns the value at the given global coordinates.
func (p *SparseMatrix[T]) Get(x, y, z int) T {
	ofs := int(math.Pow(16, float64(p.level+1))) / 2
	return p.get(x+ofs, y+ofs, z+ofs)
}

func (p *SparseMatrix[T]) get(x, y, z int) T {
	if p.IsEmpty() {
		return p.zero
	}
	if p.level == 0 && p.isSolid {
		return p.solid
	}
	s := int(math.Pow(16, float64(p.level)))
	sx := (x - p.pos[0]) / s
	sy := (y - p.pos[1]) / s
	sz := (z - p.pos[2]) / s
	sx &= 0xF
	sy &= 0xF
	sz &= 0xF
	i := sy*16*16 + sz*16 + sx
	if p.level == 0 {
		if len(p.leaves) == 0 {
			return p.zero
		}
		return p.leaves[i]
	} else {
		c := p.bcp[i]
		if c == nil {
			return p.zero
		}
		return c.get(x, y, z)
	}
}

// Set sets the value at the given global coordinates.
func (p *SparseMatrix[T]) Set(x, y, z int, v T) {
	ofs := int(math.Pow(16, float64(p.level+1))) / 2
	p.set(x+ofs, y+ofs, z+ofs, v)
}

func (p *SparseMatrix[T]) set(x, y, z int, v T) {
	s := int(math.Pow(16, float64(p.level)))
	sx := (x - p.pos[0]) / s
	sy := (y - p.pos[1]) / s
	sz := (z - p.pos[2]) / s
	sx &= 0xF
	sy &= 0xF
	sz &= 0xF
	i := sy*16*16 + sz*16 + sx
	if p.level == 0 {
		if p.isSolid && p.solid == v {
			return
		}
		if p.leaves == nil {
			p.leaves = make([]T, 16*16*16)
			f := p.zero
			if p.isSolid {
				f = p.solid
			}
			for idx := range p.leaves {
				p.leaves[idx] = f
			}
		}
		p.leaves[i] = v
		p.isSolid = false
		return
	}
	if len(p.bcp) == 0 {
		p.bcp = make([]*SparseMatrix[T], 16*16*16)
	}
	if p.bcp[i] == nil {
		p.bcp[i] = NewSparseMatrix(p.level-1, p.zero)
	}
	p.bcp[i].set(x, y, z, v)
}
