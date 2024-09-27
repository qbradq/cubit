package c3d

var bvmSliceDims = [6][3]int{
	{1, 1, 0},
	{1, 1, 0},
	{0, 1, 1},
	{0, 1, 1},
	{1, 0, 1},
	{1, 0, 1},
}

// BuildVoxelMesh builds a VoxelMesh object from the passed voxel data.
func BuildVoxelMesh(voxels [][4]uint8,
	width, height, depth int) *VoxelMesh {
	ret := NewVoxelMesh()
	z21 := func(a, b, c int) (int, int, int) {
		if a == 0 {
			a = 1
		}
		if b == 0 {
			b = 1
		}
		if c == 0 {
			c = 1
		}
		return a, b, c
	}
	var slices [6][][4]uint8
	for i := range slices {
		sd := bvmSliceDims[i]
		w, h, d := z21(width*sd[0], height*sd[1], depth*sd[2])
		slices[i] = make([][4]uint8, w*h*d)
	}
	addSlice := func(pos [3]int, f Facing, c [4]uint8) {
		p := [3]uint8{uint8(pos[0]), uint8(pos[1]), uint8(pos[2])}
		ret.AddFace(p, f, c)
	}
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
			addSlice(pos, f, c)
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
