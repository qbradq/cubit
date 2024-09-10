package cubit

import (
	"fmt"

	"github.com/qbradq/cubit/internal/c3d"
)

// Faces is the global face atlas.
var Faces *c3d.FaceAtlas = c3d.NewFaceAtlas()

// CubeRef is a reference to a single Cube object that holds all of the static
// data for the cube. The facing of cubes are packed into the three most
// significant bits within chunks.
type CubeRef uint16

// CubeRefInvalid is the invalid value for CubeRef types.
const CubeRefInvalid CubeRef = 0b0000111111111111

// cubeRefPacked returns the packed representation of the cube reference and the
// facing.
func cubeRefPacked(r CubeRef, f c3d.Facing) CubeRef {
	return (r & 0b0000111111111111) | ((CubeRef(f) & 0b111) << 13)
}

// unpack unpacks the cube reference into the cube index and facing.
func (r CubeRef) unpack() (ref CubeRef, f c3d.Facing) {
	ref = r & 0b0000111111111111
	f = c3d.Facing(r >> 13)
	return
}

// CubeDefsIndex returns the CubeRef assigned to the given id.
func CubeDefsIndex(id string) CubeRef {
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
	c.Ref = nextCubeRef
	nextCubeRef++
	cubeDefsById[c.ID] = c
	cubeDefs = append(cubeDefs, c)
	return nil
}

// nextCubeRef holds the next value to be assigned registered cubes.
var nextCubeRef CubeRef

// cubeDefs holds all of the cube definitions loaded from mods by ID.
var cubeDefsById = map[string]*Cube{}

// cubeDefs is the CubeRef to Cube lookup table.
var cubeDefs []*Cube
