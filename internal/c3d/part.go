package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/t"
)

// Part manages a voxel volume that represents one body part of a model.
type Part struct {
	Mesh        *VoxelMesh    // Voxel mesh for the part
	Origin      mgl32.Vec3    // Center point of the part
	Orientation t.Orientation // Current part orientation relative to the parent
	Children    []*Part       // Child parts, if any
}

// draw draws the part relative to the given orientation.
func (p *Part) draw(prg *program, o t.Orientation) {
	po := p.Orientation
	po.Q = o.Q.Mul(po.Q)
	po.P = o.Q.Rotate(po.P).Add(o.P)
	if p.Mesh != nil {
		mm := po.TransformMatrix()
		gl.UniformMatrix4fv(prg.uni("uModelMatrix"), 1, false, &mm[0])
		gl.Uniform3f(prg.uni("uRotationPoint"),
			p.Origin[0], p.Origin[1], p.Origin[2])
		p.Mesh.draw(prg)
	}
	for _, c := range p.Children {
		c.draw(prg, po)
	}
}
