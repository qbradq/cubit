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

// Chunk represents a 16x16x16 chunk of space.
type Chunk struct {
	pos        Position        // Position representing the global location of the chunk
	ref        ChunkRef        // Chunk self-reference
	o          c3d.Orientation // Cached orientation value for the chunk
	solid      CubeRef         // If cubes is nil, this is the CubeRef the chunk is filled with
	cubes      []CubeRef       // All cube references in the chunk
	mesh       *c3d.CubeMesh   // Mesh object used to build the vbo data
	cubesDirty bool            // When true, the cube array has changed since the last call to compile()
}

// NewChunk creates a new Chunk ready for use.
func NewChunk(r ChunkRef) *Chunk {
	ret := &Chunk{
		pos:        r.ToPosition(),
		ref:        r,
		solid:      CubeRefInvalid,
		cubesDirty: true,
	}
	ret.o = *c3d.NewOrientation(mgl32.Vec3{
		float32(ret.pos.X) + float32(ChunkWidth/2),
		float32(ret.pos.Y) + float32(ChunkHeight/2),
		float32(ret.pos.Z) + float32(ChunkDepth/2),
	}, 0, 0, 0)
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
	crp := cubeRefPacked(r, f)
	if len(c.cubes) == 0 {
		if crp == c.solid {
			return
		}
		c.cubes = make([]CubeRef, ChunkWidth*ChunkHeight*ChunkDepth)
		for i := range c.cubes {
			c.cubes[i] = c.solid
		}
	}
	ofs := p.Y * ChunkWidth * ChunkDepth
	ofs += p.Z * ChunkDepth
	ofs += p.X
	if c.cubes[ofs] == crp {
		return
	}
	c.cubes[ofs] = crp
	c.cubesDirty = true
}

// Fill fills the entire chunk with the given cube reference.
func (c *Chunk) Fill(r CubeRef, f c3d.Facing) {
	crp := cubeRefPacked(r, f)
	if len(c.cubes) == 0 && c.solid == crp {
		return
	}
	c.solid = crp
	c.cubes = nil
	c.cubesDirty = true
}

// CubeRefAt returns the packed cube reference at the given location.
func (c *Chunk) CubeRefAt(p Position) CubeRef {
	if p.X < 0 || p.X >= ChunkWidth ||
		p.Y < 0 || p.Y >= ChunkHeight ||
		p.Z < 0 || p.Z >= ChunkDepth {
		return CubeRefInvalid
	}
	if len(c.cubes) == 0 {
		return c.solid
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
					float32(p.X) + 0.5,
					float32(p.Y) + 0.5,
					float32(p.Z) + 0.5,
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
	prg.DrawCubeMesh(c.mesh, c3d.OrientationZero)
}
