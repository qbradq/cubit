package t

// IVec3 represents a cube position within the world.
type IVec3 [3]int

// Position offsets per c3d.Facing value index.
var PositionOffsets = [6]IVec3{
	{0, 0, -1},
	{0, 0, 1},
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
}

// Add adds r to this position and returns the result.
func (p IVec3) Add(r IVec3) IVec3 {
	return IVec3{
		p[0] + r[0],
		p[1] + r[1],
		p[2] + r[2],
	}
}

// Sub subtracts r from this position and returns the result.
func (p IVec3) Sub(r IVec3) IVec3 {
	return IVec3{
		p[0] - r[0],
		p[1] - r[1],
		p[2] - r[2],
	}
}

// Mul multiplies r and this position and returns the result.
func (p IVec3) Mul(r IVec3) IVec3 {
	return IVec3{
		p[0] * r[0],
		p[1] * r[1],
		p[2] * r[2],
	}
}

// Div divides this position by r and returns the result.
func (p IVec3) Div(r IVec3) IVec3 {
	return IVec3{
		p[0] / r[0],
		p[1] / r[1],
		p[2] / r[2],
	}
}

// Mod returns the result of this position modulo r.
func (p IVec3) Mod(r IVec3) IVec3 {
	return IVec3{
		p[0] % r[0],
		p[1] % r[1],
		p[2] % r[2],
	}
}
