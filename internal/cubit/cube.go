package cubit

import (
	"fmt"

	"github.com/qbradq/cubit/internal/c3d"
)

// Faces is the global face atlas.
var Faces *c3d.FaceAtlas = c3d.NewFaceAtlas()

// CubeRef is a reference to a single Cube object that holds all of the static
// data for the cube.
type CubeRef uint16

// CubeRefInvalid is the invalid value for CubeRef types.
const CubeRefInvalid CubeRef = 0xFFFF

// GetCubeDef returns the CubeRef assigned to the given id.
func GetCubeDef(id string) CubeRef {
	c, found := cubeDefsById[id]
	if !found {
		return CubeRefInvalid
	}
	return c.Ref
}

// Cube represents one cubic meter of the world.
type Cube struct {
	Ref         CubeRef          `json:"-"`           // CubeRef value assigned to the cube definition
	ID          string           `json:"-"`           // Unique ID
	Name        string           `json:"name"`        // Descriptive name
	Faces       [6]c3d.FaceIndex `json:"faces"`       // Face graphic to use for each face of the cube.
	Transparent bool             `json:"transparent"` // If true this cube can be seen through
}

// CubeInvalid is the invalid cube definition.
var CubeInvalid = &Cube{
	Ref:  0xFFFF,
	ID:   "",
	Name: "invalid",
	Faces: [6]c3d.FaceIndex{
		0xFFFF,
		0xFFFF,
		0xFFFF,
		0xFFFF,
		0xFFFF,
		0xFFFF,
	},
	Transparent: true,
}

// registerCube registers a cube definition by name.
func registerCube(c *Cube) error {
	if _, duplicate := cubeDefsById[c.ID]; duplicate {
		return fmt.Errorf("duplicate cube id %s", c.ID)
	}
	c.Ref = CubeRef(len(cubeDefs))
	cubeDefsById[c.ID] = c
	cubeDefs = append(cubeDefs, c)
	return nil
}

// cubeDefs holds all of the cube definitions loaded from mods by ID.
var cubeDefsById = map[string]*Cube{}

// cubeDefs is the CubeRef to Cube lookup table.
var cubeDefs []*Cube
