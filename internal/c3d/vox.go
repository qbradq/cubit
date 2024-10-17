package c3d

import "github.com/qbradq/cubit/internal/t"

// VoxelSource is the interface that must be implemented by voxel data providers
// to the BuildVoxelMesh function.
type VoxelSource[T any] interface {
	// Get returns the value at the given position within the voxel volume.
	Get(x, y, z int) T
	// Dimensions returns the dimensions of the voxel volume.
	Dimensions() (w, h, d int)
	// IsEmpty returns true if the value is considered empty.
	IsEmpty(v T) bool
}

// voxFace represents one rectangular face representing one or more voxels.
type voxFace[T any] struct {
	x, y, z int // Location of the lower-left corner of the face, in voxel units
	w, h, d int // Dimensions of the face, in voxel units
	v       T   // Voxel value for the face
}

// voxFaceSlice represents one slice of voxel faces.
type voxFaceSlice[T comparable] struct {
	w, h  int           // Dimensions
	e     func(T) bool  // The empty value test method
	f     t.Facing      // t.Facing of all faces in this slice
	faces []*voxFace[T] // Collection of faces to generate
}

// newVoxFaceSlice returns a new voxFaceSlice.
func newVoxFaceSlice[T comparable](w, h int, e func(T) bool, f t.Facing) *voxFaceSlice[T] {
	return &voxFaceSlice[T]{
		w:     w,
		h:     h,
		e:     e,
		f:     f,
		faces: make([]*voxFace[T], w*h),
	}
}

// addFace adds a voxel face to the slice.
func (s *voxFaceSlice[T]) addFace(x, y, z, major, minor int, v T) {
	s.faces[major*s.w+minor] = &voxFace[T]{
		x: x,
		y: y,
		z: z,
		w: 1,
		h: 1,
		d: 1,
		v: v,
	}
}

// greedyMesh uses a greedy meshing algorithm to combine adjacent faces.
func (s *voxFaceSlice[T]) greedyMesh() {
	fSet := []*voxFace[T]{}
	switch s.f {
	case t.North:
		fallthrough
	case t.South:
		for iy := 0; iy < s.h; iy++ {
			for ix := 0; ix < s.w; ix++ {
				face := s.faces[iy*s.w+ix]
				if face == nil {
					continue
				}
				s.faces[iy*s.w+ix] = nil
				fSet = append(fSet, face)
				for tx := ix + 1; tx < s.w; tx++ {
					next := s.faces[iy*s.w+tx]
					if next == nil || next.v != face.v {
						break
					}
					face.w++
					s.faces[iy*s.w+tx] = nil
				}
			nsFaceLoop:
				for ty := face.y + 1; ty < s.h; ty++ {
					for tx := face.x; tx < face.x+face.w; tx++ {
						next := s.faces[ty*s.w+tx]
						if next == nil || next.v != face.v {
							break nsFaceLoop
						}
					}
					face.h++
					for tx := face.x; tx < face.x+face.w; tx++ {
						s.faces[ty*s.w+tx] = nil
					}
				}
			}
		}
	case t.East:
		fallthrough
	case t.West:
		for iy := 0; iy < s.h; iy++ {
			for iz := 0; iz < s.w; iz++ {
				face := s.faces[iy*s.w+iz]
				if face == nil {
					continue
				}
				s.faces[iy*s.w+iz] = nil
				fSet = append(fSet, face)
				for tz := iz + 1; tz < s.w; tz++ {
					next := s.faces[iy*s.w+tz]
					if next == nil || next.v != face.v {
						break
					}
					face.d++
					s.faces[iy*s.w+tz] = nil
				}
			ewFaceLoop:
				for ty := face.y + 1; ty < s.h; ty++ {
					for tz := face.z; tz < face.z+face.d; tz++ {
						next := s.faces[ty*s.w+tz]
						if next == nil || next.v != face.v {
							break ewFaceLoop
						}
					}
					face.h++
					for tz := face.z; tz < face.z+face.d; tz++ {
						s.faces[ty*s.w+tz] = nil
					}
				}
			}
		}
	case t.Top:
		fallthrough
	case t.Bottom:
		for iz := 0; iz < s.h; iz++ {
			for ix := 0; ix < s.w; ix++ {
				face := s.faces[iz*s.w+ix]
				if face == nil {
					continue
				}
				s.faces[iz*s.w+ix] = nil
				fSet = append(fSet, face)
				for tx := ix + 1; tx < s.w; tx++ {
					next := s.faces[iz*s.w+tx]
					if next == nil || next.v != face.v {
						break
					}
					face.w++
					s.faces[iz*s.w+tx] = nil
				}
			tbFaceLoop:
				for tz := face.z + 1; tz < s.h; tz++ {
					for tx := face.x; tx < face.x+face.w; tx++ {
						next := s.faces[tz*s.w+tx]
						if next == nil || next.v != face.v {
							break tbFaceLoop
						}
					}
					face.d++
					for tx := face.x; tx < face.x+face.w; tx++ {
						s.faces[tz*s.w+tx] = nil
					}
				}
			}
		}
	}
	s.faces = fSet
}

// mesh outputs the faces of the slice to the mesh.
func (s *voxFaceSlice[T]) mesh(d Mesh[T]) {
	for _, f := range s.faces {
		if f == nil {
			continue
		}
		AddFace([3]uint8{
			uint8(f.x),
			uint8(f.y),
			uint8(f.z),
		}, [3]uint8{
			uint8(f.w),
			uint8(f.h),
			uint8(f.d),
		}, 1, s.f, f.v, d)
	}
}

// BuildVoxelMesh builds a VoxelMesh object from the passed voxel source and
// constructs faces in the destination mesh. Note that the destination mesh is
// not reset before faces are added.
func BuildVoxelMesh[T comparable](v VoxelSource[T], d Mesh[T]) {
	width, height, depth := v.Dimensions()
	// Determine if a face is required
	face := func(pos [3]int, f t.Facing) bool {
		np := [3]int{}
		np[0] = pos[0] + t.FacingOffsets[f][0]
		np[1] = pos[1] + t.FacingOffsets[f][1]
		np[2] = pos[2] + t.FacingOffsets[f][2]
		vv := v.Get(np[0], np[1], np[2])
		return v.IsEmpty(vv)
	}
	slices := []*voxFaceSlice[T]{}
	// N/S faces sweeps
	for f := t.North; f <= t.South; f++ {
		for iz := 0; iz < depth; iz++ {
			// Build face slice
			s := newVoxFaceSlice(width, height, v.IsEmpty, f)
			for iy := 0; iy < height; iy++ {
				for ix := 0; ix < width; ix++ {
					vv := v.Get(ix, iy, iz)
					if v.IsEmpty(vv) {
						continue
					}
					if face([3]int{ix, iy, iz}, f) {
						s.addFace(ix, iy, iz, iy, ix, vv)
					}
				}
			}
			slices = append(slices, s)
		}
	}
	// E/W faces sweeps
	for f := t.East; f <= t.West; f++ {
		for ix := 0; ix < width; ix++ {
			// Build face slice
			s := newVoxFaceSlice(depth, height, v.IsEmpty, f)
			for iy := 0; iy < height; iy++ {
				for iz := 0; iz < depth; iz++ {
					vv := v.Get(ix, iy, iz)
					if v.IsEmpty(vv) {
						continue
					}
					if face([3]int{ix, iy, iz}, f) {
						s.addFace(ix, iy, iz, iy, iz, vv)
					}
				}
			}
			slices = append(slices, s)
		}
	}
	// Top/Bottom faces sweeps
	for f := t.Top; f <= t.Bottom; f++ {
		for iy := 0; iy < height; iy++ {
			// Build face slice
			s := newVoxFaceSlice(width, depth, v.IsEmpty, f)
			for iz := 0; iz < depth; iz++ {
				for ix := 0; ix < width; ix++ {
					vv := v.Get(ix, iy, iz)
					if v.IsEmpty(vv) {
						continue
					}
					if face([3]int{ix, iy, iz}, f) {
						s.addFace(ix, iy, iz, iz, ix, vv)
					}
				}
			}
			slices = append(slices, s)
		}
	}
	// Build mesh
	for _, s := range slices {
		s.greedyMesh()
		s.mesh(d)
	}
}
