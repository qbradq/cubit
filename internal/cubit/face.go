package cubit

import (
	"encoding/json"
	"log"

	"github.com/qbradq/cubit/internal/c3d"
)

// FaceRef is a reference code used to index faces.
type FaceRef uint16

// Face contains all of the information about one cube face.
type Face struct {
	Graphic c3d.FaceIndex // Index into the texture atlas
}

// Faces is the c3d.FaceAtlas object loaded with all face graphics.
var Faces = c3d.NewFaceAtlas()

func (f *Face) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	log.Println(s)
	return nil
}
