package cubit

import "github.com/qbradq/cubit/internal/c3d"

// ChunkRef is a reference to a chunk in the sparse world structure. It contains
// the chunk's orientation as signed integers. The least significant 14 bits
// contains the X coordinate. The next 14 bits the Z coordinate. And the final
// 4 bits the Y coordinate.
type ChunkRef uint32

// ToPosition returns the position encoded in the chunk reference.
func (r ChunkRef) ToPosition() Position {
	return Position{
		X: int(r&0x00003FFF) - 0x2000,
		Z: int((r&0x0FFFC000)>>14) - 0x2000,
		Y: int((r&0xF0000000)>>28) - 0x8,
	}
}

// ChunkRefFromCoords returns the position encoded as a chunk reference.
func ChunkRefFromCoords(p Position) ChunkRef {
	return ChunkRef(uint32(p.X+0x2000) |
		uint32(p.Z+0x2000)<<14 |
		uint32(p.Y+0x8)<<28)
}

// World manages the state of the entire world.
type World struct {
	chunks map[ChunkRef]*Chunk // Map of all loaded chunks
}

// NewWorld returns a new World object read for use.
func NewWorld() *World {
	return &World{
		chunks: map[ChunkRef]*Chunk{},
	}
}

// GetChunkByRef returns the chunk by reference.
func (w *World) GetChunkByRef(r ChunkRef) *Chunk {
	c, found := w.chunks[r]
	if !found {
		c = w.generateChunk(r)
		w.chunks[r] = c
	}
	return c
}

// generateChunk generates the for the given chunk reference and returns it.
func (w *World) generateChunk(r ChunkRef) *Chunk {
	ret := NewChunk(r)
	p := r.ToPosition()
	// For now we just hard-code the chunk generation.
	rDirt := CubeDefsIndex("/cubit/dirt")
	rGrass := CubeDefsIndex("/cubit/grass")
	if p.Y < 0 {
		ret.Fill(rDirt, c3d.North)
	}
	if p.Y > 0 {
		ret.Fill(CubeRefInvalid, c3d.North)
	}
	if p.Y == 0 {
		// Flat
		for iy := 0; iy < ChunkHeight; iy++ {
			for iz := 0; iz < ChunkDepth; iz++ {
				for ix := 0; ix < ChunkWidth; ix++ {
					if iy < 11 {
						ret.Set(Pos(ix, iy, iz), rDirt, c3d.North)
					} else if iy == 11 {
						ret.Set(Pos(ix, iy, iz), rGrass, c3d.North)
					}
				}
			}
		}
		// Dirt house at origin
		if p.X == 0 && p.Z == 0 {
			for ix := 5; ix <= 11; ix++ {
				ret.Set(Pos(ix, 12, 9), rDirt, c3d.North)
				ret.Set(Pos(ix, 13, 9), rDirt, c3d.North)
				if ix == 9 {
					ret.Set(Pos(ix, 12, 5), rDirt, c3d.North)
				} else if ix != 7 {
					ret.Set(Pos(ix, 12, 5), rDirt, c3d.North)
					ret.Set(Pos(ix, 13, 5), rDirt, c3d.North)
				}
			}
			for iz := 6; iz <= 8; iz++ {
				ret.Set(Pos(5, 12, iz), rDirt, c3d.North)
				ret.Set(Pos(5, 13, iz), rDirt, c3d.North)
				ret.Set(Pos(11, 12, iz), rDirt, c3d.North)
				ret.Set(Pos(11, 13, iz), rDirt, c3d.North)
			}
			for iz := 5; iz <= 9; iz++ {
				for ix := 5; ix <= 11; ix++ {
					ret.Set(Pos(ix, 14, iz), rDirt, c3d.North)
				}
			}
		}
	}
	return ret
}

// SetCube sets the cube and facing at the given position in the world.
func (w *World) SetCube(p Position, r CubeRef, f c3d.Facing) {
	w.GetChunkByRef(ChunkRefFromCoords(p.Div(chunkDimensions))).Set(
		p.Mod(chunkDimensions), r, f)
}
