package mod

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/t"
)

// ModelPartDescriptor describes how to initialize a part attached to a model.
type ModelPartDescriptor struct {
	ID          string                `json:"id"`          // ID of the part for animations
	Mesh        string                `json:"mesh"`        // Absolute path to the part mesh
	Origin      mgl32.Vec3            `json:"origin"`      // Center point / rotation point of the part
	Orientation t.Orientation         `json:"orientation"` // T-pose orientation
	Children    []ModelPartDescriptor `json:"children"`    // Child parts
}

// ModelDescriptor describes how to initialize a model.
type ModelDescriptor struct {
	Bounds t.AABB               `json:"bounds"` // Bounding box of the model
	Root   *ModelPartDescriptor `json:"root"`   // Root part
}

// animationContext is a context for an animation.
type animationContext struct {
	a       *c3d.Animation // Animation to play
	running bool           // If true the animation is running
}

// Model describes a hierarchy of parts with defined animations.
type Model struct {
	DrawDescriptor *c3d.ModelDrawDescriptor    // The draw descriptor this model manages
	joints         map[string]*mgl32.Quat      // Joint name to rotation quaternion mapping
	animations     map[string]animationContext // Animation layer name to running animation map
}

// modelsMap is the mapping of resource paths to model descriptors.
var modelsMap = map[string]*ModelDescriptor{}

// registerModel registers a model descriptor by resource path.
func registerModel(p string, m *ModelDescriptor) error {
	if _, duplicate := modelsMap[p]; duplicate {
		return fmt.Errorf("duplicate model %s", p)
	}
	modelsMap[p] = m
	return nil
}

// newPart creates a new part (and its children) from the part descriptor
// given.
func newPart(d *ModelPartDescriptor) *c3d.Part {
	ret := &c3d.Part{
		ID:     d.ID,
		Mesh:   GetPartMesh(d.Mesh),
		Origin: d.Origin,
		Orientation: t.Orientation{
			P: d.Orientation.P.Mul(t.VoxelScale),
			Q: d.Orientation.Q,
		},
	}
	for _, cd := range d.Children {
		ret.Children = append(ret.Children, newPart(&cd))
	}
	return ret
}

// NewModel constructs a new model from the model descriptor at the resource
// path given.
func NewModel(p string) *Model {
	d := modelsMap[p]
	if d == nil {
		return nil
	}
	ret := &Model{
		DrawDescriptor: &c3d.ModelDrawDescriptor{
			ID: 1,
			Bounds: c3d.AABB{
				Bounds: d.Bounds,
			},
			Orientation: t.O(),
			Root:        newPart(d.Root),
		},
		joints:     make(map[string]*mgl32.Quat),
		animations: make(map[string]animationContext),
	}
	var fn func(p *c3d.Part)
	fn = func(p *c3d.Part) {
		ret.joints[p.ID] = &p.Orientation.Q
		for _, c := range p.Children {
			fn(c)
		}
	}
	fn(ret.DrawDescriptor.Root)
	return ret
}

// StartAnimation starts the animation identified by the resource path on the
// model on the given animation layer.
func (m *Model) StartAnimation(p, l string) {
	var a *c3d.Animation
	if ac, found := m.animations[l]; found {
		if ac.running {
			return
		}
		a = ac.a
	} else {
		a = getAnimation(p, m.joints)
	}
	if a == nil {
		return
	}
	m.animations[l] = animationContext{
		a:       a,
		running: true,
	}
	a.Play()
}

// Update updates the animations of the model.
func (m *Model) Update(dt float32) {
	for _, a := range m.animations {
		if !a.running {
			continue
		}
		a.a.Update(dt)
	}
}
