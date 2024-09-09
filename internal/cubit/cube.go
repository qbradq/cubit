package cubit

import "github.com/qbradq/cubit/internal/c3d"

// CubeRef is a reference to a single Cube object that holds all of the static
// data for the cube. The facing of cubes are packed into the three most
// significant bits within chunks.
type CubeRef uint16

// CubeDefsIndex returns the

// Cube represents one cubic meter of the world.
type Cube struct {
	Name  string           `json:"name"`  // Descriptive name
	Faces [6]c3d.FaceIndex `json:"faces"` // Fixed tile graphic to use for each face of the cube.
}

// CubeDefs holds all of the cube definitions loaded from mods by ID.
var CubeDefs = map[string]*Cube{}
