package c3d

var bvmSliceDims = [6][3]int{
	{1, 1, 0},
	{1, 1, 0},
	{0, 1, 1},
	{0, 1, 1},
	{1, 0, 1},
	{1, 0, 1},
}

// VoxelSource is the interface that must be implemented by voxel data providers
// to the BuildVoxelMesh function.
type VoxelSource interface {
	// Get returns the value at the given position within the voxel volume.
	Get(x, y, z int) [4]uint8
	// Dimensions returns the dimensions of the voxel volume.
	Dimensions() (w, h, d int)
}

// BuildVoxelMesh builds a VoxelMesh object from the passed voxel data.
func BuildVoxelMesh(v VoxelSource) *VoxelMesh {
	ret := NewVoxelMesh()
	width, height, depth := v.Dimensions()
	face := func(pos [3]int, f Facing, c [4]uint8) {
		np := [3]int{}
		np[0] = pos[0] + FacingOffsets[f][0]
		np[1] = pos[1] + FacingOffsets[f][1]
		np[2] = pos[2] + FacingOffsets[f][2]
		if v.Get(np[0], np[1], np[2])[3] < 255 {
			p := [3]uint8{uint8(pos[0]), uint8(pos[1]), uint8(pos[2])}
			ret.AddFace(p, f, c)
		}
	}
	// Visible face only build
	for iy := 0; iy < height; iy++ {
		for iz := 0; iz < depth; iz++ {
			for ix := 0; ix < width; ix++ {
				c := v.Get(ix, iy, iz)
				if c[3] < 255 {
					continue
				}
				for i := North; i <= Bottom; i++ {
					face([3]int{ix, iy, iz}, i, c)
				}
			}
		}
	}
	return ret
}
