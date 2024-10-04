package t

// Cube represents one cubic meter of the world.
type Cube struct {
	Ref         CubeRef      `json:"-"`           // CubeRef value assigned to the cube definition
	ID          string       `json:"-"`           // Unique ID
	Name        string       `json:"name"`        // Descriptive name
	Faces       [6]FaceIndex `json:"faces"`       // Face graphic to use for each face of the cube.
	Transparent bool         `json:"transparent"` // If true this cube can be seen through
}

// CubeInvalid is the invalid cube definition.
var CubeInvalid = &Cube{
	Ref:  0xFFFF,
	ID:   "",
	Name: "invalid",
	Faces: [6]FaceIndex{
		0xFFFF,
		0xFFFF,
		0xFFFF,
		0xFFFF,
		0xFFFF,
		0xFFFF,
	},
	Transparent: true,
}
