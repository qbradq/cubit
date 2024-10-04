package mod

import (
	"fmt"

	"github.com/qbradq/cubit/internal/c3d"
)

// partsMeshMap is the mapping of resource paths to part meshes.
var partsMeshMap = map[string]*c3d.VoxelMesh{}

// registerPartMesh registers a part mesh by resource path.
func registerPartMesh(p string, m *c3d.VoxelMesh) error {
	if _, duplicate := partsMeshMap[p]; duplicate {
		return fmt.Errorf("duplicate part mesh %s", p)
	}
	partsMeshMap[p] = m
	return nil
}

// GetPartMesh returns the part mesh by resource path.
func GetPartMesh(path string) *c3d.VoxelMesh {
	return partsMeshMap[path]
}
