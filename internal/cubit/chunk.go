package cubit

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
)

const ChunkWidth int = 16
const ChunkHeight int = 16
const ChunkDepth int = 16

var chunkDimensions = Position{
	X: ChunkWidth,
	Y: ChunkHeight,
	Z: ChunkDepth,
}

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

// Chunk represents a 16x16x16 chunk of space.
type Chunk struct {
	pos        Position          // Position representing the global location of the chunk
	ref        ChunkRef          // Chunk self-reference
	o          c3d.Orientation   // Cached orientation value for the chunk
	solid      Cell              // If cubes is nil, this is the CubeRef the chunk is filled with
	cubes      []Cell            // All references in the chunk
	mesh       *c3d.CubeMesh     // Mesh object used to build the vbo data for the cube mesh
	vox        []*Vox            // List of all vox models to draw
	vos        []c3d.Orientation // List of orientations to use when drawing vox
	cubesDirty bool              // When true, the cube array has changed since the last call to compile()
}

// NewChunk creates a new Chunk ready for use.
func NewChunk(r ChunkRef) *Chunk {
	ret := &Chunk{
		pos:        r.ToPosition(),
		ref:        r,
		solid:      CellInvalid,
		cubesDirty: true,
	}
	ret.pos.X *= ChunkWidth
	ret.pos.Y *= ChunkHeight
	ret.pos.Z *= ChunkDepth
	ret.o = *c3d.NewOrientation(mgl32.Vec3{
		float32(ret.pos.X) + float32(ChunkWidth/2),
		float32(ret.pos.Y) + float32(ChunkHeight/2),
		float32(ret.pos.Z) + float32(ChunkDepth/2),
	}, 0, 0, 0)
	for i := range ret.cubes {
		ret.cubes[i] = CellInvalid
	}
	return ret
}

// Set sets a cube by reference and facing.
func (c *Chunk) Set(p Position, r Cell) {
	if p.X < 0 || p.X >= ChunkWidth ||
		p.Y < 0 || p.Y >= ChunkHeight ||
		p.Z < 0 || p.Z >= ChunkDepth {
		return
	}
	if len(c.cubes) == 0 {
		if r == c.solid {
			return
		}
		c.cubes = make([]Cell, ChunkWidth*ChunkHeight*ChunkDepth)
		for i := range c.cubes {
			c.cubes[i] = c.solid
		}
	}
	ofs := p.Y * ChunkWidth * ChunkDepth
	ofs += p.Z * ChunkDepth
	ofs += p.X
	if c.cubes[ofs] == r {
		return
	}
	c.cubes[ofs] = r
	c.cubesDirty = true
}

// Fill fills the entire chunk with the given cube reference.
func (c *Chunk) Fill(r Cell) {
	if len(c.cubes) == 0 && c.solid == r {
		return
	}
	c.solid = r
	c.cubes = nil
	c.cubesDirty = true
}

// CellAt returns the packed cube reference at the given location.
func (c *Chunk) CellAt(p Position) Cell {
	if p.X < 0 || p.X >= ChunkWidth ||
		p.Y < 0 || p.Y >= ChunkHeight ||
		p.Z < 0 || p.Z >= ChunkDepth {
		return CellInvalid
	}
	if len(c.cubes) == 0 {
		return c.solid
	}
	ofs := p.Y * ChunkWidth * ChunkDepth
	ofs += p.Z * ChunkDepth
	ofs += p.X
	return c.cubes[ofs]
}

// At returns a pointer to the cube definition or voxel model definition at the
// given cell.
func (c *Chunk) At(p Position) (*Cube, *Vox, c3d.Facing) {
	r := c.CellAt(p)
	if r == CellInvalid {
		return nil, nil, c3d.North
	}
	return r.Decompose()
}

// compile compiles the chunk's mesh.
func (c *Chunk) compile() {
	var center mgl32.Vec3
	var p Position
	var cube *Cube
	var vox *Vox
	var f c3d.Facing
	face := func(side c3d.Facing) {
		cell := c.CellAt(p.Add(PositionOffsets[side]))
		nc, _, _ := cell.Decompose()
		if nc != nil && !nc.Transparent {
			return
		}
		c.mesh.AddFace(center, side, cube.Faces[side])
	}
	if c.mesh == nil {
		c.mesh = c3d.NewCubeMesh(false)
	}
	c.mesh.Reset()
	c.vox = c.vox[:0]
	c.vos = c.vos[:0]
	for p.Y = 0; p.Y < ChunkHeight; p.Y++ {
		for p.Z = 0; p.Z < ChunkDepth; p.Z++ {
			for p.X = 0; p.X < ChunkWidth; p.X++ {
				cell := c.CellAt(p)
				cube, vox, f = cell.Decompose()
				if vox != nil {
					c.vox = append(c.vox, vox)
					c.vos = append(c.vos, *c3d.NewOrientation(
						mgl32.Vec3{
							float32(p.X+c.pos.X) + 0.5,
							float32(p.Y+c.pos.Y) + 0.5,
							float32(p.Z+c.pos.Z) + 0.5,
						},
						c3d.FacingToOrientation[f].GetPitch(),
						c3d.FacingToOrientation[f].GetYaw(),
						c3d.FacingToOrientation[f].GetRoll(),
					))
				}
				if cube == nil || cube.Transparent {
					continue
				}
				center = mgl32.Vec3{
					float32(p.X) + 0.5,
					float32(p.Y) + 0.5,
					float32(p.Z) + 0.5,
				}
				face(c3d.North)
				face(c3d.South)
				face(c3d.East)
				face(c3d.West)
				face(c3d.Top)
				face(c3d.Bottom)
			}
		}
	}
}

// Draw draws the chunk using the provided c3d program.
func (c *Chunk) Draw(prg *c3d.Program) {
	if c.cubesDirty {
		c.compile()
		c.cubesDirty = false
	}
	prg.DrawCubeMesh(c.mesh, c3d.OrientationZero)
	for i, v := range c.vox {
		prg.DrawCubeMesh(v.mesh, &c.vos[i])
	}
}
