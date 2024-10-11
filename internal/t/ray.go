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

// WorldIntersection represents a ray intersection with the world.
type WorldIntersection struct {
	Position IVec3   // World coordinate of the cell hit
	Face     Facing  // Face struck
	Cube     CubeRef // The cube intersected, if any
	Vox      VoxRef  // The voxel model intersected, if any
	Facing   Facing  // Facing of the cube or voxel model intersected
}

// IntersectWorld returns a description of the point at which the ray
// intersects w, or nil if there is no intersection. Adopted from
// https://stackoverflow.com/questions/12367071/how-do-i-initialize-the-t-variables-in-a-fast-voxel-traversal-algorithm-for-ray
func (r *Ray) IntersectWorld(w *World) *WorldIntersection {
	SIGN := func(x float32) float32 {
		if x > 0 {
			return 1
		}
		if x < 0 {
			return -1
		}
		return 0
	}
	FRAC0 := func(x float32) float32 {
		return x - float32(math.Floor(float64(x)))
	}
	FRAC1 := func(x float32) float32 {
		return 1.0 - x + float32(math.Floor(float64(x)))
	}

	var tMaxX, tMaxY, tMaxZ, tDeltaX, tDeltaY, tDeltaZ float32
	var voxel IVec3
	var x1, y1, z1 float32 // start point
	var x2, y2, z2 float32 // end point
	start := r.O
	end := r.N.Mul(r.L).Add(r.O)
	x1 = start[0]
	y1 = start[1]
	z1 = start[2]
	x2 = end[0]
	y2 = end[1]
	z2 = end[2]

	dx := SIGN(x2 - x1)
	if dx != 0 {
		tDeltaX = float32(math.Min(float64(dx/(x2-x1)), 10000000.0))
	} else {
		tDeltaX = 10000000.0
	}
	if dx > 0 {
		tMaxX = tDeltaX * FRAC1(x1)
	} else {
		tMaxX = tDeltaX * FRAC0(x1)
	}
	voxel[0] = int(x1)

	dy := SIGN(y2 - y1)
	if dy != 0 {
		tDeltaY = float32(math.Min(float64(dy/(y2-y1)), 10000000.0))
	} else {
		tDeltaY = 10000000.0
	}
	if dy > 0 {
		tMaxY = tDeltaY * FRAC1(y1)
	} else {
		tMaxY = tDeltaY * FRAC0(y1)
	}
	voxel[1] = int(y1)

	dz := SIGN(z2 - z1)
	if dz != 0 {
		tDeltaZ = float32(math.Min(float64(dz/(z2-z1)), 10000000.0))
	} else {
		tDeltaZ = 10000000.0
	}
	if dz > 0 {
		tMaxZ = tDeltaZ * FRAC1(z1)
	} else {
		tMaxZ = tDeltaZ * FRAC0(z1)
	}
	voxel[2] = int(z1)

	var face Facing
	for {
		cell := w.GetCell(voxel)
		cRef, vRef, f := cell.Decompose()
		if cRef != CubeRefInvalid || vRef != VoxRefInvalid {
			return &WorldIntersection{
				Position: voxel,
				Cube:     cRef,
				Vox:      vRef,
				Face:     face,
				Facing:   f,
			}
		}
		if tMaxX < tMaxY {
			if tMaxX < tMaxZ {
				voxel[0] += int(dx)
				tMaxX += tDeltaX
				face = West
				if tDeltaX < 0 {
					face = East
				}
			} else {
				voxel[2] += int(dz)
				tMaxZ += tDeltaZ
				face = North
				if tDeltaZ < 0 {
					face = South
				}
			}
		} else {
			if tMaxY < tMaxZ {
				voxel[1] += int(dy)
				tMaxY += tDeltaY
				face = Bottom
				if tDeltaY < 0 {
					face = Top
				}
			} else {
				voxel[2] += int(dz)
				tMaxZ += tDeltaZ
				face = North
				if tDeltaZ < 0 {
					face = South
				}
			}
		}
		if tMaxX > 1 && tMaxY > 1 && tMaxZ > 1 {
			break
		}
	}
	return nil
}
