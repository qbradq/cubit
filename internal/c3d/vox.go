package c3d

// VoxBuilder manages memory buffers used for building optimized vox images.
type VoxBuilder struct {
}

// NewVoxBuilder returns a new VoxBuilder ready for use.
func NewVoxBuilder() *VoxBuilder {
	return &VoxBuilder{}
}

// BuildCubeMesh builds a CubeMesh object from the passed voxel data.
func (g *VoxBuilder) BuildCubeMesh(voxels [][4]uint8,
	width, height, depth int) *CubeMesh {
	ret := NewCubeMesh(true)
	face := func(pos [3]int, f Facing, c [4]uint8) {
		transparent := false
		np := [3]int{}
		np[0] = pos[0] + FacingOffsets[f][0]
		np[1] = pos[1] + FacingOffsets[f][1]
		np[2] = pos[2] + FacingOffsets[f][2]
		if np[0] < 0 || np[0] >= width ||
			np[1] < 0 || np[1] >= height ||
			np[2] < 0 || np[2] >= depth {
			transparent = true
		}
		if !transparent {
			idx := np[1] * width * depth
			idx += np[2] * width
			idx += np[0]
			transparent = voxels[idx][3] < 255
		}
		if transparent {
			ret.AddVoxelFace(pos, f, c)
		}
	}
	// Visible face only build
	for iy := 0; iy < height; iy++ {
		for iz := 0; iz < depth; iz++ {
			for ix := 0; ix < width; ix++ {
				idx := iy * width * depth
				idx += iz * width
				idx += ix
				if voxels[idx][3] == 0 {
					continue
				}
				for i := North; i <= Bottom; i++ {
					face([3]int{ix, iy, iz}, i, voxels[idx])
				}
			}
		}
	}
	return ret
}
