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
func (c *Chunk) Set(x, y, z int, r CubeRef, f c3d.Facing) {
	if x < 0 || x >= ChunkWidth ||
		y < 0 || y >= ChunkHeight ||
		z < 0 || z >= ChunkDepth {
		return
	}
	ofs := y * ChunkWidth * ChunkDepth
	ofs += z * ChunkDepth
	ofs += x
	c.cubes[ofs] = cubeRefPacked(r, f)
	c.cubesDirty = true
}

// At returns the packed cube reference at the given location.
func (c *Chunk) At(x, y, z int) CubeRef {
	if x < 0 || x >= ChunkWidth ||
		y < 0 || y >= ChunkHeight ||
		z < 0 || z >= ChunkDepth {
		return CubeRefInvalid
	}
	ofs := y * ChunkWidth * ChunkDepth
	ofs += z * ChunkDepth
	ofs += x
	return c.cubes[ofs]
}

// compile compiles the chunk's mesh.
func (c *Chunk) compile() {
	if c.mesh == nil {
		c.mesh = c3d.NewCubeMesh()
	}
	c.mesh.Reset()
	for iy := 0; iy < ChunkHeight; iy++ {
		for iz := 0; iz < ChunkDepth; iz++ {
			for ix := 0; ix < ChunkWidth; ix++ {
				r := c.At(ix, iy, iz)
				cr, _ := r.unpack()
				if cr == CubeRefInvalid {
					continue
				}
				center := mgl32.Vec3{
					float32(ix-ChunkWidth/2) - 0.5,
					float32(iy-ChunkHeight/2) + 0.5,
					float32(iz-ChunkDepth/2) + 0.5,
				}
				cube := cubeDefs[cr]
				c.mesh.AddFace(center, c3d.North, cube.Faces[0])
				c.mesh.AddFace(center, c3d.South, cube.Faces[1])
				c.mesh.AddFace(center, c3d.East, cube.Faces[2])
				c.mesh.AddFace(center, c3d.West, cube.Faces[3])
				c.mesh.AddFace(center, c3d.Top, cube.Faces[4])
				c.mesh.AddFace(center, c3d.Bottom, cube.Faces[5])
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
