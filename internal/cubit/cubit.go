package cubit

import "github.com/qbradq/cubit/internal/c3d"

// UITiles is the face atlas for all UI tiles in all mods.
var UITiles *c3d.FaceAtlas = c3d.NewFaceAtlas()

// uiTilesMap is a mapping of mod / ui tile index strings to face indexes in
// UITiles.
var uiTilesMap = map[string]c3d.FaceIndex{}

// GetUITile returns the index of the UI tile at the given path.
func GetUITile(p string) c3d.FaceIndex {
	if idx, found := uiTilesMap[p]; found {
		return idx
	}
	return c3d.FaceIndexInvalid
}
