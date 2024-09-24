package cubit

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/vox"
)

const ChunkWidth int = 16
const ChunkHeight int = 16
const ChunkDepth int = 16

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
	aabb       vox.AABB          // Bounding box of the chunk
	o          c3d.Orientation   // Cached orientation value for the chunk
	cubes      vox.Chunk[Cell]   // Cells of the chunk
	mesh       *c3d.CubeMesh     // Mesh object used to build the vbo data for the cube mesh
	vox        []*Vox            // List of all vox models to draw
	vos        []c3d.Orientation // List of orientations to use when drawing vox
	cubesDirty bool              // When true, the cube array has changed since the last call to compile()
}

// NewChunk creates a new Chunk ready for use.
func NewChunk(p Position) *Chunk {
	ret := &Chunk{
		pos:        p,
		cubesDirty: true,
	}
	ret.pos.X *= ChunkWidth
	ret.pos.Y *= ChunkHeight
	ret.pos.Z *= ChunkDepth
	ret.cubes = *vox.NewChunk(ret.pos.X, ret.pos.Y, ret.pos.Z, CellInvalid,
		CellInvalid)
	ret.aabb = vox.AABB{
		Min: mgl32.Vec3{
			float32(ret.pos.X),
			float32(ret.pos.Y),
			float32(ret.pos.Z),
		},
		Max: mgl32.Vec3{
			float32(ret.pos.X + ChunkWidth),
			float32(ret.pos.Y + ChunkHeight),
			float32(ret.pos.Z + ChunkDepth),
		},
	}
	ret.o = *c3d.NewOrientation(mgl32.Vec3{
		float32(ret.pos.X) + float32(ChunkWidth/2),
		float32(ret.pos.Y) + float32(ChunkHeight/2),
		float32(ret.pos.Z) + float32(ChunkDepth/2),
	}, 0, 0, 0)
	return ret
}

// SetRelative sets a cube by reference and facing.
func (c *Chunk) SetRelative(p Position, r Cell) {
	if c.cubes.SetRelative(p.X, p.Y, p.Z, r) {
		c.cubesDirty = true
	}
}

// Fill fills the entire chunk with the given cube reference.
func (c *Chunk) Fill(r Cell) {
	c.cubes.Fill(r)
	c.cubesDirty = true
}

// GetRelative returns the packed cube reference at the given location.
func (c *Chunk) GetRelative(p Position) Cell {
	return c.cubes.GetRelative(p.X, p.Y, p.Z)
}

// At returns a pointer to the cube definition or voxel model definition at the
// given cell.
func (c *Chunk) At(p Position) (*Cube, *Vox, c3d.Facing) {
	r := c.GetRelative(p)
	if r == CellInvalid {
		return nil, nil, c3d.North
	}
	return r.Decompose()
}

// compile compiles the chunk's mesh.
func (c *Chunk) compile() {
	var p Position
	var cube *Cube
	var vox *Vox
	var f c3d.Facing
	face := func(side c3d.Facing) {
		np := p.Add(PositionOffsets[side])
		cell := c.cubes.GetRelative(np.X, np.Y, np.Z)
		nc, _, _ := cell.Decompose()
		if nc != nil && !nc.Transparent {
			return
		}
		c.mesh.AddFace(byte(p.X), byte(p.Y), byte(p.Z), side, cube.Faces[side])
	}
	if c.mesh == nil {
		c.mesh = c3d.NewCubeMesh(c.aabb)
	}
	c.mesh.Reset()
	c.vox = c.vox[:0]
	c.vos = c.vos[:0]
	for p.Y = 0; p.Y < ChunkHeight; p.Y++ {
		for p.Z = 0; p.Z < ChunkDepth; p.Z++ {
			for p.X = 0; p.X < ChunkWidth; p.X++ {
				cell := c.GetRelative(p)
				cube, vox, f = cell.Decompose()
				if vox != nil {
					c.vox = append(c.vox, vox)
					c.vos = append(c.vos, *c3d.NewOrientation(
						mgl32.Vec3{
							float32(p.X) + 0.5,
							float32(p.Y) + 0.5,
							float32(p.Z) + 0.5,
						},
						c3d.FacingToOrientation[f].GetPitch(),
						c3d.FacingToOrientation[f].GetYaw(),
						c3d.FacingToOrientation[f].GetRoll(),
					))
				}
				if cube == nil || cube.Transparent {
					continue
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

// Add adds the chunk to the app for rendering.
func (c *Chunk) Add(app *c3d.App) {
	if c.cubesDirty {
		c.compile()
		c.cubesDirty = false
	}
	app.AddCubeMesh(c.mesh, c3d.OrientationZero)
	for i := range c.vox {
		app.AddVoxelMesh(c.vox[i].mesh, &c.vos[i])
	}
}

// // Draw draws the chunk using the provided c3d program.
// func (c *Chunk) Draw(prg *c3d.Program) {
// 	if c.cubesDirty {
// 		c.compile()
// 		c.cubesDirty = false
// 	}
// 	prg.DrawCubeMesh(c.mesh, c3d.OrientationZero)
// 	for i, v := range c.vox {
// 		prg.DrawCubeMesh(v.mesh, &c.vos[i])
// 	}
// }
