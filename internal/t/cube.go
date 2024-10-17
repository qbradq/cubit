package t

// FacingMap maps orientation facing to face index.
var FacingMap = [6][6]Facing{
	{
		North,
		South,
		East,
		West,
		Top,
		Bottom,
	},
	{
		South,
		North,
		West,
		East,
		Top,
		Bottom,
	},
	{
		West,
		East,
		North,
		South,
		Top,
		Bottom,
	},
	{
		East,
		West,
		South,
		North,
		Top,
		Bottom,
	},
	{
		Bottom,
		Top,
		East,
		West,
		North,
		South,
	},
	{
		Top,
		Bottom,
		East,
		West,
		South,
		North,
	},
}

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
