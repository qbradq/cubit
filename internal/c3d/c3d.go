package c3d

import (
	"fmt"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// VoxelScale is the scale of a single voxel.
const VoxelScale float32 = 1.0 / 16.0

// VirtualScreenWidth is the width of the virtual 2D screen in pixels.
const VirtualScreenWidth int = 320

// VirtualScreenHeight is the width of the virtual 2D screen in pixels.
const VirtualScreenHeight int = 180

// NinePatch describes the resources used to generate an arbitrarily sized
// rectangle skinned with nine tiles stretched over the area.
type NinePatch [9]FaceIndex

// ColoredString is a string associated with an RGB color.
type ColoredString struct {
	String string   // The string
	Color  [3]uint8 // The color
}

// CubeMeshDrawDescriptor describes how and where to draw a cube mesh.
type CubeMeshDrawDescriptor struct {
	ID       uint32     // ID
	Mesh     *CubeMesh  // The mesh to draw
	Position mgl32.Vec3 // Position of the bottom-north-west corner of the mesh
}

// VoxelMeshDrawDescriptor describes how and where to draw a cube mesh.
type VoxelMeshDrawDescriptor struct {
	ID       uint32     // ID
	Mesh     *VoxelMesh // The mesh to draw
	Position mgl32.Vec3 // Position of the center point of the voxel mesh
	Facing   Facing     // Facing of the mesh
}

// ChunkDrawDescriptor describes how and where to render the static portions of
// a scene.
type ChunkDrawDescriptor struct {
	ID       uint32                     // ID
	CubeDD   CubeMeshDrawDescriptor     // The draw descriptor for the cube mesh
	VoxelDDs []*VoxelMeshDrawDescriptor // Draw descriptors for all voxel meshes contained within the chunk
}

// ModelDrawDescriptor describes how and where to render a dynamic model.
type ModelDrawDescriptor struct {
	ID          uint32      // ID
	Orientation Orientation // Origin of the model
	Root        *Part       // Root part that everything else is attached to
}

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

var XAxis mgl32.Vec3 = mgl32.Vec3{1, 0, 0}
var YAxis mgl32.Vec3 = mgl32.Vec3{0, 1, 0}
var ZAxis mgl32.Vec3 = mgl32.Vec3{0, 0, 1}

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

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {
	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)
	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)
		log := "NO LOG"
		if logLength > 0 {
			gls := gl.Str(strings.Repeat("\x00", int(logLength)))
			getObjInfoLogFn(glHandle, logLength, nil, gls)
			log = gl.GoStr(gls)
		}
		return fmt.Errorf("%s: %s", failMsg, log)
	}
	return nil
}

var CubeFacingOffsets = [6][6][3]uint8{
	{ // North
		{1, 1, 0}, // TL
		{0, 1, 0}, // TR
		{1, 0, 0}, // BL
		{1, 0, 0}, // BL
		{0, 1, 0}, // TR
		{0, 0, 0}, // BR
	},
	{ // South
		{0, 1, 1}, // TL
		{1, 1, 1}, // TR
		{0, 0, 1}, // BL
		{0, 0, 1}, // BL
		{1, 1, 1}, // TR
		{1, 0, 1}, // BR
	},
	{ // East
		{1, 1, 1}, // TL
		{1, 1, 0}, // TR
		{1, 0, 1}, // BL
		{1, 0, 1}, // BL
		{1, 1, 0}, // TR
		{1, 0, 0}, // BR
	},
	{ // West
		{0, 1, 0}, // TL
		{0, 1, 1}, // TR
		{0, 0, 0}, // BL
		{0, 0, 0}, // BL
		{0, 1, 1}, // TR
		{0, 0, 1}, // BR
	},
	{ // Top
		{0, 1, 0}, // TL
		{1, 1, 0}, // TR
		{0, 1, 1}, // BL
		{0, 1, 1}, // BL
		{1, 1, 0}, // TR
		{1, 1, 1}, // BR
	},
	{ // Bottom
		{1, 0, 0}, // TL
		{0, 0, 0}, // TR
		{1, 0, 1}, // BL
		{1, 0, 1}, // BL
		{0, 0, 0}, // TR
		{0, 0, 1}, // BR
	},
}

var cubeFacingOffsetsI = [6][6][3]byte{
	{ // North
		{1, 1, 0}, // TL
		{0, 1, 0}, // TR
		{1, 0, 0}, // BL
		{1, 0, 0}, // BL
		{0, 1, 0}, // TR
		{0, 0, 0}, // BR
	},
	{ // South
		{0, 1, 1}, // TL
		{1, 1, 1}, // TR
		{0, 0, 1}, // BL
		{0, 0, 1}, // BL
		{1, 1, 1}, // TR
		{1, 0, 1}, // BR
	},
	{ // East
		{1, 1, 1}, // TL
		{1, 1, 0}, // TR
		{1, 0, 1}, // BL
		{1, 0, 1}, // BL
		{1, 1, 0}, // TR
		{1, 0, 0}, // BR
	},
	{ // West
		{0, 1, 0}, // TL
		{0, 1, 1}, // TR
		{0, 0, 0}, // BL
		{0, 0, 0}, // BL
		{0, 1, 1}, // TR
		{0, 0, 1}, // BR
	},
	{ // Top
		{0, 1, 0}, // TL
		{1, 1, 0}, // TR
		{0, 1, 1}, // BL
		{0, 1, 1}, // BL
		{1, 1, 0}, // TR
		{1, 1, 1}, // BR
	},
	{ // Bottom
		{1, 0, 0}, // TL
		{0, 0, 0}, // TR
		{1, 0, 1}, // BL
		{1, 0, 1}, // BL
		{0, 0, 0}, // TR
		{0, 0, 1}, // BR
	},
}
