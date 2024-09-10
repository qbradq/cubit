package cubit

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
)

const ChunkWidth int = 32
const ChunkHeight int = 32
const ChunkDepth int = 32

// Chunk represents a 16x16x16 cube of the world.
type Chunk struct {
	c3d.Orientation
	cubes      [ChunkWidth * ChunkHeight * ChunkDepth]CubeRef
	mesh       *c3d.CubeMesh // Mesh object used to build the vbo data
	cubesDirty bool          // When true, the cube array has changed since the last call to compile()
}

// NewChunk creates a new Chunk ready for use.
func NewChunk() *Chunk {
	ret := &Chunk{
		Orientation: *c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, 0, 0, 0),
		cubesDirty:  true,
	}
	for i := range ret.cubes {
		ret.cubes[i] = CubeRefInvalid
	}
	return ret
}

// Set sets a cube by reference and facing.
func (c *Chunk) Set(p Position, r CubeRef, f c3d.Facing) {
	if p.X < 0 || p.X >= ChunkWidth ||
		p.Y < 0 || p.Y >= ChunkHeight ||
		p.Z < 0 || p.Z >= ChunkDepth {
		return
	}
	ofs := p.Y * ChunkWidth * ChunkDepth
	ofs += p.Z * ChunkDepth
	ofs += p.X
	c.cubes[ofs] = cubeRefPacked(r, f)
	c.cubesDirty = true
}

// CubeRefAt returns the packed cube reference at the given location.
func (c *Chunk) CubeRefAt(p Position) CubeRef {
	if p.X < 0 || p.X >= ChunkWidth ||
		p.Y < 0 || p.Y >= ChunkHeight ||
		p.Z < 0 || p.Z >= ChunkDepth {
		return CubeRefInvalid
	}
	ofs := p.Y * ChunkWidth * ChunkDepth
	ofs += p.Z * ChunkDepth
	ofs += p.X
	return c.cubes[ofs]
}

// CubeAt returns the cube definition at the given location.
func (c *Chunk) CubeAt(p Position) *Cube {
	r := c.CubeRefAt(p)
	if r == CubeRefInvalid {
		return CubeInvalid
	}
	return cubeDefs[r]
}

// compile compiles the chunk's mesh.
func (c *Chunk) compile() {
	var center mgl32.Vec3
	var p Position
	var cube *Cube
	face := func(side c3d.Facing) {
		n := c.CubeAt(p.Add(PositionOffsets[side]))
		if !n.Transparent {
			return
		}
		c.mesh.AddFace(center, side, cube.Faces[side])
	}
	if c.mesh == nil {
		c.mesh = c3d.NewCubeMesh()
	}
	c.mesh.Reset()
	for p.Y = 0; p.Y < ChunkHeight; p.Y++ {
		for p.Z = 0; p.Z < ChunkDepth; p.Z++ {
			for p.X = 0; p.X < ChunkWidth; p.X++ {
				r := c.CubeRefAt(p)
				cr, _ := r.unpack()
				if cr == CubeRefInvalid {
					continue
				}
				center = mgl32.Vec3{
					float32(p.X-ChunkWidth/2) - 0.5,
					float32(p.Y-ChunkHeight/2) + 0.5,
					float32(p.Z-ChunkDepth/2) + 0.5,
				}
				cube = cubeDefs[cr]
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
	prg.DrawCubeMesh(c.mesh, &c.Orientation)
}
