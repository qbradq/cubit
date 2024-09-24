package cubit

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/vox"
)

// World manages the state of the entire world.
type World struct {
	chunks *vox.OctalTree[*Chunk] // Sparse chunks volume
}

// NewWorld returns a new World object read for use.
func NewWorld() *World {
	return &World{
		chunks: vox.NewOctalTreeRoot[*Chunk](4096),
	}
}

// GetChunk returns the chunk in chunk coordinates.
func (w *World) GetChunk(p Position) *Chunk {
	v := mgl32.Vec3{float32(p.X), float32(p.Y), float32(p.Z)}
	c := w.chunks.Get(v)
	if c == nil {
		c = w.generateChunk(p)
		w.chunks.Set(v, c)
	}
	return c
}

// GetChunkForPosition returns the chunk in world coordinates.
func (w *World) GetChunkForPosition(p Position) *Chunk {
	v := mgl32.Vec3{float32(p.X / ChunkWidth), float32(p.Y / ChunkHeight),
		float32(p.Z / ChunkDepth)}
	c := w.chunks.Get(v)
	if c == nil {
		c = w.generateChunk(p)
		w.chunks.Set(v, c)
	}
	return c
}

// generateChunk generates the chunk for the given chunk coordinates and
// returns it.
func (w *World) generateChunk(p Position) *Chunk {
	ret := NewChunk(p)
	n := c3d.North
	// For now we just hard-code the chunk generation.
	rDirt := CubeDefsIndex("/cubit/dirt")
	rGrass := CubeDefsIndex("/cubit/grass")
	if p.Y < 0 {
		ret.Fill(CellForCube(rDirt, n))
	}
	if p.Y > 0 {
		ret.Fill(CellInvalid)
	}
	if p.Y == 0 {
		// Flat
		for iy := 0; iy < ChunkHeight; iy++ {
			for iz := 0; iz < ChunkDepth; iz++ {
				for ix := 0; ix < ChunkWidth; ix++ {
					if iy < 11 {
						ret.SetRelative(Pos(ix, iy, iz), CellForCube(rDirt, n))
					} else if iy == 11 {
						ret.SetRelative(Pos(ix, iy, iz), CellForCube(rGrass, n))
					}
				}
			}
		}
		// Dirt house at origin
		if p.X == 0 && p.Z == 0 {
			for ix := 5; ix <= 11; ix++ {
				ret.SetRelative(Pos(ix, 12, 9), CellForCube(rDirt, c3d.North))
				ret.SetRelative(Pos(ix, 13, 9), CellForCube(rDirt, c3d.North))
				if ix == 9 {
					ret.SetRelative(Pos(ix, 12, 5), CellForCube(rDirt, c3d.North))
					ret.SetRelative(Pos(ix, 13, 5), CellForVox(GetVoxByPath("/cubit/window0").Ref, c3d.North))
				} else if ix != 7 {
					ret.SetRelative(Pos(ix, 12, 5), CellForCube(rDirt, c3d.North))
					ret.SetRelative(Pos(ix, 13, 5), CellForCube(rDirt, c3d.North))
				}
			}
			for iz := 6; iz <= 8; iz++ {
				ret.SetRelative(Pos(5, 12, iz), CellForCube(rDirt, c3d.North))
				ret.SetRelative(Pos(5, 13, iz), CellForCube(rDirt, c3d.North))
				ret.SetRelative(Pos(11, 12, iz), CellForCube(rDirt, c3d.North))
				ret.SetRelative(Pos(11, 13, iz), CellForCube(rDirt, c3d.North))
			}
			for iz := 5; iz <= 9; iz++ {
				for ix := 5; ix <= 11; ix++ {
					ret.SetRelative(Pos(ix, 14, iz), CellForCube(rDirt, c3d.North))
				}
			}
		}
	}
	return ret
}

// SetCell sets the cube and facing at the given position in the world.
func (w *World) SetCell(p Position, c Cell, f c3d.Facing) {
	w.GetChunkForPosition(p).cubes.Set(p.X, p.Y, p.Z, c)
}
