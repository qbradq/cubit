package mod

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
)

// ModelPartDescriptor describes how to initialize a part attached to a model.
type ModelPartDescriptor struct {
	Mesh        string                `json:"mesh"`        // Absolute path to the part mesh
	CenterPoint mgl32.Vec3            `json:"centerPoint"` // Center point / rotation point of the part
	Orientation c3d.Orientation       `json:"orientation"` // T-pose orientation
	Children    []ModelPartDescriptor `json:"children"`    // Child parts
}

// ModelDescriptor describes how to initialize a model.
type ModelDescriptor struct {
	Root ModelPartDescriptor `json:"root"` // Root part
}

// Model describes a hierarchy of parts with defined animations.
type Model struct {
	DrawDescriptor *c3d.ModelDrawDescriptor // The draw descriptor this model manages
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
		Mesh:        GetPartMesh(d.Mesh),
		CenterPoint: d.CenterPoint,
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
			ID:   1,
			Root: newPart(&d.Root),
		},
	}
	return ret
}
