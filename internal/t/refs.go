package t

// CubeRef is a reference to a single Cube object that holds all of the static
// data for the cube.
type CubeRef uint16

// CubeRefInvalid is the invalid value for CubeRef types.
const CubeRefInvalid CubeRef = 0xFFFF

// Facing encodes one of the facing values, North, South, East, West, Up, Down.
type Facing uint8

const (
	North Facing = iota
	South
	East
	West
	Top
	Bottom
)

// FacingOffsets are the offsets from the center voxel to the voxel in the
// direction of the indexed facing.
var FacingOffsets = [6][3]int{
	{0, 0, -1},
	{0, 0, 1},
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
}

// VoxRef encodes a reference to a vox model.
type VoxRef uint16

// VoxRefInvalid is the invalid valid value for VoxRef variables.
const VoxRefInvalid VoxRef = 0xFFFF
