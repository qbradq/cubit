package mod

import (
	"fmt"

	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/t"
)

// Faces is the global face atlas.
var Faces *c3d.FaceAtlas = c3d.NewFaceAtlas()

// GetCubeRef returns the CubeRef assigned to the given id.
func GetCubeRef(id string) t.CubeRef {
	c, found := cubeDefsById[id]
	if !found {
		return t.CubeRefInvalid
	}
	return c.Ref
}

// GetCubeDef returns the cube definition.
func GetCubeDef(id string) *t.Cube {
	c, found := cubeDefsById[id]
	if !found {
		return nil
	}
	return c
}

// registerCube registers a cube definition by name.
func registerCube(c *t.Cube) error {
	if _, duplicate := cubeDefsById[c.ID]; duplicate {
		return fmt.Errorf("duplicate cube id %s", c.ID)
	}
	c.Ref = t.CubeRef(len(CubeDefs))
	cubeDefsById[c.ID] = c
	CubeDefs = append(CubeDefs, c)
	return nil
}

// cubeDefs holds all of the cube definitions loaded from mods by ID.
var cubeDefsById = map[string]*t.Cube{}

// CubeDefs is the CubeRef to Cube lookup table.
var CubeDefs []*t.Cube
