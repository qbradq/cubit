package t

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Ray represents a line in 3D space by defining an origin and normalized
// facing direction.
type Ray struct {
	O mgl32.Vec3 // Origin
	N mgl32.Vec3 // Normal
	L float32    // Length
	i mgl32.Vec3 // Inverse normal
}

// NewRay constructs a new Ray ready for use.
func NewRay(origin, dir mgl32.Vec3, length float32) *Ray {
	ret := &Ray{
		O: origin,
		N: dir.Normalize(),
		L: length,
	}
	ret.i[0] = 1.0 / ret.N[0]
	ret.i[1] = 1.0 / ret.N[1]
	ret.i[2] = 1.0 / ret.N[2]
	return ret
}

// IntersectsAABB returns true if the ray intersects the given AABB.
func (r *Ray) IntersectsAABB(b AABB) bool {
	/*
	   double tx1 = (box.min.X - optray.x0.X)*optray.n_inv.X;
	   double tx2 = (box.max.X - optray.x0.X)*optray.n_inv.X;

	   double tmin = dmnsn_min(tx1, tx2);
	   double tmax = dmnsn_max(tx1, tx2);

	   double ty1 = (box.min.Y - optray.x0.Y)*optray.n_inv.Y;
	   double ty2 = (box.max.Y - optray.x0.Y)*optray.n_inv.Y;

	   tmin = dmnsn_max(tmin, dmnsn_min(ty1, ty2));
	   tmax = dmnsn_min(tmax, dmnsn_max(ty1, ty2));

	   double tz1 = (box.min.Z - optray.x0.Z)*optray.n_inv.Z;
	   double tz2 = (box.max.Z - optray.x0.Z)*optray.n_inv.Z;

	   tmin = dmnsn_max(tmin, dmnsn_min(tz1, tz2));
	   tmax = dmnsn_min(tmax, dmnsn_max(tz1, tz2));

	   return tmax >= dmnsn_max(0.0, tmin) && tmin < t;
	*/
	tx1 := float64((b[0][0] - r.O[0]) * r.i[0])
	tx2 := float64((b[1][0] - r.O[0]) * r.i[0])
	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)
	ty1 := float64((b[0][1] - r.O[1]) * r.i[1])
	ty2 := float64((b[1][1] - r.O[1]) * r.i[1])
	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))
	tz1 := float64((b[0][2] - r.O[2]) * r.i[2])
	tz2 := float64((b[1][2] - r.O[2]) * r.i[2])
	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))
	return tmax >= math.Max(0.0, tmin) && tmin < float64(r.L)
}
