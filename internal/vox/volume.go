package vox

// Volume manages an arbitrarily-sized dense matrix of values.
type Volume[T comparable] struct {
	x, y, z int // Volume bottom-north-west corner
	w, h, d int // Volume dimensions
	v       []T // Matrix of values
	inv     T   // Invalid value
}

// NewVolume returns a new Volume of the given dimensions and filled with the
// given fill value.
func NewVolume[T comparable](x, y, z, w, h, d int, fill, invalid T) *Volume[T] {
	ret := &Volume[T]{
		x:   x,
		y:   y,
		z:   z,
		w:   w,
		h:   h,
		d:   d,
		v:   make([]T, w*h*d),
		inv: invalid,
	}
	ret.Fill(fill)
	return ret
}

// Fill fills the volume with the given fill value.
func (v *Volume[T]) Fill(f T) {
	for i := range v.v {
		v.v[i] = f
	}
}

// Get returns the value at the given location, or the invalid value if out of
// bounds.
func (v *Volume[T]) Get(x, y, z int) T {
	return v.GetRelative(x-v.x, y-v.y, z-v.z)
}

// GetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the volume.
func (v *Volume[T]) GetRelative(x, y, z int) T {
	if x < 0 || x >= v.w || y < 0 || y >= v.h || z < 0 || z >= v.d {
		return v.inv
	}
	return v.v[(z*v.w*v.h)+(y*v.w)+x]
}

// Set sets the value at the given location. If the location is out of bounds
// this is a no-op. True is returned if the contents of the volume were altered.
func (v *Volume[T]) Set(x, y, z int, f T) bool {
	return v.SetRelative(x-v.x, y-v.y, z-v.z, f)
}

// SetRelative is like Get(), but the location is relative to the bottom-north-
// west corner of the volume.
func (v *Volume[T]) SetRelative(x, y, z int, f T) bool {
	if x < 0 || x >= v.w || y < 0 || y >= v.h || z < 0 || z >= v.d {
		return false
	}
	if v.v[(z*v.w*v.h)+(y*v.w)+x] == f {
		return false
	}
	v.v[(z*v.w*v.h)+(y*v.w)+x] = f
	return true
}
