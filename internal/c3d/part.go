package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// Part manages a voxel volume that represents one body part of a model.
type Part struct {
	Mesh        *VoxelMesh  // Voxel mesh for the part
	CenterPoint mgl32.Vec3  // Center point of the part
	Orientation Orientation // Current part orientation relative to the parent
	Children    []*Part     // Child parts, if any
}

// draw draws the part relative to the given orientation.
func (p *Part) draw(prg *program, o *Orientation) {
	po := p.Orientation.Add(o)
	if p.Mesh != nil {
		(&po).Translate(p.CenterPoint)
		mm := po.TransformMatrix()
		gl.UniformMatrix4fv(prg.uni("uModelMatrix"), 1, false, &mm[0])
		gl.Uniform3f(prg.uni("uRotationPoint"),
			p.CenterPoint[0],
			p.CenterPoint[1],
			p.CenterPoint[2],
		)
		p.Mesh.draw(prg)
	}
	for _, c := range p.Children {
		c.draw(prg, &po)
	}
}
