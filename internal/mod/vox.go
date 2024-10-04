package mod

import (
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/t"
	"github.com/qbradq/cubit/internal/util"
)

// Vox manages an RGBA voxel image.
type Vox struct {
	Ref                  t.VoxRef       // Voxel reference
	Mesh                 *c3d.VoxelMesh // Voxel mesh
	width, height, depth int            // Dimensions
	voxels               [][4]uint8     // RGBA voxels
}

// GetVoxByPath returns the vox model by mod path.
func GetVoxByPath(p string) *Vox {
	return voxIndex[p]
}

// registerVox registers the vox model by path.
func registerVox(p string, v *Vox) {
	if _, duplicate := voxIndex[p]; duplicate {
		panic("duplicate vox path " + p)
	}
	v.Ref = t.VoxRef(len(VoxDefs))
	voxIndex[p] = v
	VoxDefs = append(VoxDefs, v)
}

// voxIndex is the global registry of vox models.
var voxIndex = map[string]*Vox{}

// VoxDefs is the list of voxel definitions by internal ID.
var VoxDefs = []*Vox{}

// NewVox creates a new Vox object ready to use.
func NewVox(v *util.Vox) *Vox {
	mesh := c3d.NewVoxelMesh()
	c3d.BuildVoxelMesh[[4]uint8](v, mesh)
	return &Vox{
		Mesh:   mesh,
		width:  v.Width,
		height: v.Height,
		depth:  v.Depth,
		voxels: v.Voxels,
	}
}
