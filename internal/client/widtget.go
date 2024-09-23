package client

import "github.com/qbradq/cubit/internal/c3d"

const (
	layerCursor    uint16 = 0xFFFF
	layerConsole   uint16 = 0xFFFE
	layerCrosshair uint16 = 0xFFFD
	layerHighest   uint16 = 0xFFFC
)

// baseWidget is a mixin struct that provides common functionality to all
// widgets.
type baseWidget struct {
	c3d.UIMesh
}

// newBaseWidget returns a new base widget ready for use.
func newBaseWidget(app *c3d.App) *baseWidget {
	ret := &baseWidget{
		UIMesh: *app.NewUIMesh(),
	}
	return ret
}

// add adds the widget to the app to be drawn.
func (w *baseWidget) add(app *c3d.App) {
	app.AddUIMesh(&w.UIMesh)
}
