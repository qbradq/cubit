package client

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/cubit"
)

// voxD is a voxel description encoding position, orientation, and model.
type voxD struct {
	vox *cubit.Vox
	pos cubit.Position
	f   c3d.Facing
}

// Chunk represents a 16x16x16 chunk of space.
type Chunk struct {
	pos   cubit.Position  // Position representing the global location of the chunk
	o     c3d.Orientation // Cached orientation value for the chunk
	mesh  *c3d.CubeMesh   // Mesh object used to build the vbo data for the cube mesh
	vox   []voxD          // List of all vox models to draw
	dirty bool            // When true, the cube array has changed since the last call to compile()
}

// NewChunk creates a new Chunk ready for use.
func NewChunk(p cubit.Position) *Chunk {
	ret := &Chunk{
		pos:   p,
		dirty: true,
	}
	ret.pos.X *= 16
	ret.pos.Y *= 16
	ret.pos.Z *= 16
	ret.o = *c3d.NewOrientation(mgl32.Vec3{
		float32(ret.pos.X) + float32(16/2),
		float32(ret.pos.Y) + float32(16/2),
		float32(ret.pos.Z) + float32(16/2),
	}, 0, 0, 0)
	return ret
}

// compile compiles the chunk's mesh.
func (c *Chunk) compile(w *cubit.World) {
	var p cubit.Position
	var cube *cubit.Cube
	var vox *cubit.Vox
	var f c3d.Facing
	face := func(side c3d.Facing) {
		np := p.Add(cubit.PositionOffsets[side]).Add(c.pos)
		cell := w.GetCell(np)
		nc, _, _ := cell.Decompose()
		if nc != nil && !nc.Transparent {
			return
		}
		c.mesh.AddFace(byte(p.X), byte(p.Y), byte(p.Z), side, cube.Faces[side])
	}
	if c.mesh == nil {
		c.mesh = c3d.NewCubeMesh()
	}
	c.mesh.Reset()
	c.vox = c.vox[:0]
	for p.Y = 0; p.Y < 16; p.Y++ {
		for p.Z = 0; p.Z < 16; p.Z++ {
			for p.X = 0; p.X < 16; p.X++ {
				gp := p.Add(c.pos)
				cell := w.GetCell(gp)
				cube, vox, f = cell.Decompose()
				if vox != nil {
					c.vox = append(c.vox, voxD{
						vox: vox,
						pos: p,
						f:   f,
					})
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
func (c *Chunk) Add(w *cubit.World, app *c3d.App) {
	if c.dirty {
		c.compile(w)
		c.dirty = false
	}
	app.AddCubeMesh(c.mesh, c3d.NewOrientation(
		mgl32.Vec3{
			float32(c.pos.X),
			float32(c.pos.Y),
			float32(c.pos.Z),
		}, 0, 0, 0))
	for _, d := range c.vox {
		d.vox.Add(app, c.pos.Add(d.pos), d.f)
	}
}
