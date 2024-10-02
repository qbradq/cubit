package cubit

import (
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/util"
)

// VoxRef encodes a reference to a vox model.
type VoxRef uint16

// VoxRefInvalid is the invalid valid value for VoxRef variables.
const VoxRefInvalid VoxRef = 0xFFFF

// Vox manages an RGBA voxel image.
type Vox struct {
	Ref                  VoxRef         // Voxel reference
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
	v.Ref = VoxRef(len(voxDefs))
	voxIndex[p] = v
	voxDefs = append(voxDefs, v)
}

// voxIndex is the global registry of vox models.
var voxIndex = map[string]*Vox{}

// voxDefs is the list of voxel definitions by internal ID.
var voxDefs = []*Vox{}

// NewVox creates a new Vox object ready to use.
func NewVox(v *util.Vox) *Vox {
	return &Vox{
		Mesh:   c3d.BuildVoxelMesh(v),
		width:  v.Width,
		height: v.Height,
		depth:  v.Depth,
		voxels: v.Voxels,
	}
}
