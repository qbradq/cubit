package client

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/cubit"
)

// ChunkLocationRef is a relative index into the chunk.
type ChunkLocationRef uint16

// NewChunkLocationRef returns a new chunk location reference for the given
// position.
func NewChunkLocationRef(p cubit.Position) ChunkLocationRef {
	return ChunkLocationRef(
		(p.X & 0xF) |
			((p.Y & 0xF) << 4) |
			((p.Z & 0xF) << 8),
	)
}

// ToPosition returns the chunk relative position encoded into the reference.
func (r ChunkLocationRef) ToPosition() cubit.Position {
	return cubit.Pos(
		int(r&0xF),
		int((r&0xF0)>>4),
		int((r&0xF00)>>8),
	)
}

// Chunk represents a 16x16x16 chunk of space.
type Chunk struct {
	p   cubit.Position           // Chunk position in world coordinates
	c   cubit.Chunk              // The chunk we are modeling
	cdd *c3d.ChunkDrawDescriptor // Draw descriptor holding all 3D assets
	lcr uint32                   // Last compiled revision of the chunk data
}

// NewChunk creates a new Chunk ready for use.
func NewChunk(p cubit.Position) *Chunk {
	ref := cubit.NewChunkRefForWorldPosition(p)
	ret := &Chunk{
		p: p,
		c: *world.GetChunkByRef(ref),
		cdd: &c3d.ChunkDrawDescriptor{
			ID: uint32(ref),
			CubeDD: c3d.CubeMeshDrawDescriptor{
				ID:   uint32(ref),
				Mesh: c3d.NewCubeMesh(),
				Position: mgl32.Vec3{
					float32(p.X),
					float32(p.Y),
					float32(p.Z),
				},
			},
			VoxelDDs: []*c3d.VoxelMeshDrawDescriptor{},
		},
	}
	return ret
}

// compile compiles the chunk's mesh.
func (c *Chunk) compile() {
	var p cubit.Position
	var cube *cubit.Cube
	var vox *cubit.Vox
	var f c3d.Facing
	face := func(side c3d.Facing) {
		np := p.Add(cubit.PositionOffsets[side]).Add(c.p)
		cell := world.GetCell(np)
		nc, _, _ := cell.Decompose()
		if nc != nil && !nc.Transparent {
			return
		}
		c.cdd.CubeDD.Mesh.AddFace(byte(p.X), byte(p.Y), byte(p.Z), side,
			cube.Faces[side])
	}
	c.cdd.CubeDD.Mesh.Reset()
	c.cdd.VoxelDDs = c.cdd.VoxelDDs[:0]
	for p.Y = 0; p.Y < 16; p.Y++ {
		for p.Z = 0; p.Z < 16; p.Z++ {
			for p.X = 0; p.X < 16; p.X++ {
				gp := p.Add(c.p)
				cell := world.GetCell(gp)
				cube, vox, f = cell.Decompose()
				if vox != nil {
					vp := c.p.Add(p)
					vmdd := &c3d.VoxelMeshDrawDescriptor{
						ID:          uint32(NewChunkLocationRef(p)),
						Mesh:        vox.Mesh,
						CenterPoint: mgl32.Vec3{8, 8, 8},
						Position: mgl32.Vec3{
							float32(vp.X),
							float32(vp.Y),
							float32(vp.Z),
						},
						Facing: f,
					}
					c.cdd.VoxelDDs = append(c.cdd.VoxelDDs, vmdd)
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

// Update does periodic updates on the chunk for client-side things like chunk
// compilation.
func (c *Chunk) Update() {
	if c.lcr < c.c.Revision {
		c.compile()
		c.lcr = c.c.Revision
	}
}
