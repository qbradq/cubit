package cubit

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/util"
)

// ftd is the facing to orientation table.
var ftd = [6]*c3d.Orientation{
	c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, 0, 0, 0),
	c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, 0, math.Pi, 0),
	c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, 0, -math.Pi*0.5, 0),
	c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, 0, math.Pi*0.5, 0),
	c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, math.Pi*0.5, 0, 0),
	c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, -math.Pi*0.5, 0, 0),
}

// VoxRef encodes a reference to a vox model.
type VoxRef uint16

// VoxRefInvalid is the invalid valid value for VoxRef variables.
const VoxRefInvalid VoxRef = 0xFFFF

// Vox manages an RGBA voxel image.
type Vox struct {
	Ref  VoxRef          // Voxel reference
	rgba *c3d.RGBAVolume // Voxel volume
	v    *util.Vox       // Voxel model file
}

// GetVoxByPath returns the vox model by mod path.
func GetVoxByPath(p string) *Vox {
	return voxIndex[p]
}

// RegisterVox registers the vox model by path.
func RegisterVox(p string, v *Vox) {
	if _, duplicate := voxIndex[p]; duplicate {
		panic("duplicate vox path " + p)
	}
	v.Ref = VoxRef(len(voxDefs))
	voxIndex[p] = v
	voxDefs = append(voxDefs, v)
}

// voxIndex is the global registry of vox models.
var voxIndex = map[string]*Vox{}

// voxDefs is the list of voxel definitions by internal ID.
var voxDefs = []*Vox{}

// NewVox creates a new Vox object ready to use.
func NewVox(v *util.Vox) *Vox {
	return &Vox{
		v: v,
	}
}

// Add adds this voxel's mesh to the app's list of voxel models to draw.
func (g *Vox) Add(a *c3d.App, pos Position, f c3d.Facing) {
	o := ftd[f]
	o.Position = mgl32.Vec3{
		float32(pos.X),
		float32(pos.Y),
		float32(pos.Z),
	}
	if g.rgba == nil {
		g.rgba = c3d.NewRGBVolume()
		g.rgba.SetData(g.v)
	}
	a.AddRGBAVolume(g.rgba, o)
}
