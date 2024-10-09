package c3d

import (
	"fmt"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/t"
)

// NinePatch describes the resources used to generate an arbitrarily sized
// rectangle skinned with nine tiles stretched over the area.
type NinePatch [9]t.FaceIndex

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
	Facing   t.Facing   // Facing of the mesh
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
	ID          uint32        // ID
	Bounds      AABB          // Bounds of the model
	Origin      mgl32.Vec3    // Origin of the model in model space
	Orientation t.Orientation // Origin of the model
	Root        *Part         // Root part that everything else is attached to
}

// LineMeshDrawDescriptor describes how and where to render a line mesh.
type LineMeshDrawDescriptor struct {
	ID          uint32        // ID
	Orientation t.Orientation // Origin of the model
	Mesh        *LineMesh     // Mesh
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

// AABB represents an axis-aligned bounding box with an optional bounds mesh.
type AABB struct {
	Bounds t.AABB
	mesh   *LineMesh
}

// draw draws the bounds of the line.
func (b *AABB) draw(prg *program) {
	if b.mesh == nil {
		b.mesh = NewLineMesh()
		c := [4]uint8{0, 255, 0, 255}
		bnw := mgl32.Vec3{b.Bounds[0][0], b.Bounds[0][1], b.Bounds[0][2]}
		bne := mgl32.Vec3{b.Bounds[1][0], b.Bounds[0][1], b.Bounds[0][2]}
		bsw := mgl32.Vec3{b.Bounds[0][0], b.Bounds[0][1], b.Bounds[1][2]}
		bse := mgl32.Vec3{b.Bounds[1][0], b.Bounds[0][1], b.Bounds[1][2]}
		tnw := mgl32.Vec3{b.Bounds[0][0], b.Bounds[1][1], b.Bounds[0][2]}
		tne := mgl32.Vec3{b.Bounds[1][0], b.Bounds[1][1], b.Bounds[0][2]}
		tsw := mgl32.Vec3{b.Bounds[0][0], b.Bounds[1][1], b.Bounds[1][2]}
		tse := mgl32.Vec3{b.Bounds[1][0], b.Bounds[1][1], b.Bounds[1][2]}
		b.mesh.Line(bnw, bne, c)
		b.mesh.Line(bne, bse, c)
		b.mesh.Line(bse, bsw, c)
		b.mesh.Line(bsw, bnw, c)
		b.mesh.Line(tnw, tne, c)
		b.mesh.Line(tne, tse, c)
		b.mesh.Line(tse, tsw, c)
		b.mesh.Line(tsw, tnw, c)
		b.mesh.Line(bnw, tnw, c)
		b.mesh.Line(bne, tne, c)
		b.mesh.Line(bsw, tsw, c)
		b.mesh.Line(bse, tse, c)
	}
	b.mesh.draw(prg)
}
