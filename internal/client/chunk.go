package client

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/mod"
	"github.com/qbradq/cubit/internal/t"
)

// ChunkLocationRef is a relative index into the chunk.
type ChunkLocationRef uint16

// NewChunkLocationRef returns a new chunk location reference for the given
// position.
func NewChunkLocationRef(p t.IVec3) ChunkLocationRef {
	return ChunkLocationRef(
		(p[0] & 0xF) |
			((p[1] & 0xF) << 4) |
			((p[2] & 0xF) << 8),
	)
}

// ToPosition returns the chunk relative position encoded into the reference.
func (r ChunkLocationRef) ToPosition() t.IVec3 {
	return t.IVec3{
		int(r & 0xF),
		int((r & 0xF0) >> 4),
		int((r & 0xF00) >> 8),
	}
}

// Chunk represents a 16x16x16 chunk of space.
type Chunk struct {
	p   t.IVec3                  // Chunk position in world coordinates
	c   *t.Chunk                 // The chunk we are modeling
	cdd *c3d.ChunkDrawDescriptor // Draw descriptor holding all 3D assets
	lcr uint32                   // Last compiled revision of the chunk data
	lvr uint32                   // Last compiled revision of the chunk vox data
}

// NewChunk creates a new Chunk ready for use.
func NewChunk(p t.IVec3) *Chunk {
	ref := t.NewChunkRefForWorldPosition(p)
	ret := &Chunk{
		p: p,
		c: world.GetChunkByRef(ref),
		cdd: &c3d.ChunkDrawDescriptor{
			ID: uint32(ref),
			CubeDD: c3d.CubeMeshDrawDescriptor{
				ID:   uint32(ref),
				Mesh: c3d.NewCubeMesh(mod.CubeDefs),
				Position: mgl32.Vec3{
					float32(p[0]),
					float32(p[1]),
					float32(p[2]),
				},
			},
			VoxelDDs: []*c3d.VoxelMeshDrawDescriptor{},
		},
	}
	return ret
}

// Update does periodic updates on the chunk for client-side things like chunk
// compilation.
func (c *Chunk) Update() {
	if c.lcr < c.c.Revision {
		c.cdd.CubeDD.Mesh.Reset()
		c3d.BuildVoxelMesh[t.Cell](c.c, c.cdd.CubeDD.Mesh)
		c.lcr = c.c.Revision
	}
	if c.lvr < c.c.VoxRevision {
		c.cdd.VoxelDDs = c.cdd.VoxelDDs[:0]
		for iz := 0; iz < 16; iz++ {
			for iy := 0; iy < 16; iy++ {
				for ix := 0; ix < 16; ix++ {
					cell := c.c.Get(ix, iy, iz)
					_, vr, f := cell.Decompose()
					if vr == t.VoxRefInvalid {
						continue
					}
					c.cdd.VoxelDDs = append(c.cdd.VoxelDDs,
						&c3d.VoxelMeshDrawDescriptor{
							ID:   1,
							Mesh: mod.VoxDefs[vr].Mesh,
							Position: mgl32.Vec3{
								float32(ix),
								float32(iy),
								float32(iz),
							},
							Facing: f,
						})
				}
			}
		}
		c.lvr = c.c.VoxRevision
	}
}
