package cubit

const CubeWidth int = 16

// Chunk represents a 16x16x16 cube of the world.
type Chunk struct {
	Cubes [CubeWidth * CubeWidth * CubeWidth]CubeRef
}

// NewChunk creates a new Chunk ready for use.
func NewChunk() *Chunk {
	ret := &Chunk{}
	return ret
}
