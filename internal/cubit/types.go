package cubit

// Position represents a cube position within the world.
type Position struct {
	X, Y, Z int
}

// Position offsets per c3d.Facing value index.
var PositionOffsets = [6]Position{
	{0, 0, -1},
	{0, 0, 1},
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
}

// Pos creates a new position from X, Y, and Z coordinates.
func Pos(x, y, z int) Position {
	return Position{
		X: x,
		Y: y,
		Z: z,
	}
}

// Add adds r to this position and returns the result.
func (p Position) Add(r Position) Position {
	return Position{
		X: p.X + r.X,
		Y: p.Y + r.Y,
		Z: p.Z + r.Z,
	}
}

// Sub subtracts r from this position and returns the result.
func (p Position) Sub(r Position) Position {
	return Position{
		X: p.X - r.X,
		Y: p.Y - r.Y,
		Z: p.Z + r.Z,
	}
}

// Mul multiplies r and this position and returns the result.
func (p Position) Mul(r Position) Position {
	return Position{
		X: p.X * r.X,
		Y: p.Y * r.Y,
		Z: p.Z * r.Z,
	}
}

// Div divides this position by r and returns the result.
func (p Position) Div(r Position) Position {
	return Position{
		X: p.X / r.X,
		Y: p.Y / r.Y,
		Z: p.Z / r.Z,
	}
}

// Mod returns the result of this position modulo r.
func (p Position) Mod(r Position) Position {
	return Position{
		X: p.X % r.X,
		Y: p.Y % r.Y,
		Z: p.Z % r.Z,
	}
}
