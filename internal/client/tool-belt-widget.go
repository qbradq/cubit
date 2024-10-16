package client

import (
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/mod"
	"github.com/qbradq/cubit/internal/t"
)

// toolBeltWidget implements a widget that presents 10 usable items bound to
// hotkeys 1-0.
type toolBeltWidget struct {
	baseWidget
	items    [10]t.Cell                      // Cells in the tool belt
	cds      [10]*c3d.CubeMeshDrawDescriptor // Cube meshes for the tool belt cells
	selected int                             // Index of the selected slot
	hover    int                             // Index of the slot the mouse is hovering over
	dirty    bool                            // If true, the UI needs to be rebuilt
}

// newToolBeltWidget returns a new toolBeltWidget ready for use.
func newToolBeltWidget(app *c3d.App) *toolBeltWidget {
	ret := &toolBeltWidget{
		baseWidget: *newBaseWidget(app),
		dirty:      true,
	}
	ret.UIMesh.Layer = layerToolBelt
	x := (t.VirtualScreenGlyphsWide - 40) / 2 * t.VirtualScreenGlyphSize
	y := (t.VirtualScreenGlyphsHigh - 4) * t.VirtualScreenGlyphSize
	x += t.VirtualScreenGlyphSize * 2
	y += t.VirtualScreenGlyphSize * 2
	for i := range ret.items {
		ret.items[i] = t.CellInvalid
		ret.cds[i] = ret.CubeMeshIcon(x, y, 8, 8, 8, nil, t.North, mod.CubeDefs)
		ret.UIMesh.Cubes = append(ret.UIMesh.Cubes, ret.cds[i])
		x += t.VirtualScreenGlyphSize * 4
	}
	return ret
}

// getSelectedCell returns the cell describing the item in the currently
// selected slot.
func (w *toolBeltWidget) getSelectedCell() t.Cell {
	return w.items[w.selected]
}

// updateMesh builds the UI mesh.
func (w *toolBeltWidget) updateMesh() {
	w.Reset(true, true)
	x := (t.VirtualScreenGlyphsWide - 40) / 2
	y := t.VirtualScreenGlyphsHigh - 4
	for i := range w.items {
		np := npWindow
		if i == w.selected || i == w.hover {
			np = npHighlight
		}
		w.UIMesh.NinePatch(
			x*t.VirtualScreenGlyphSize,
			y*t.VirtualScreenGlyphSize,
			4*t.VirtualScreenGlyphSize,
			4*t.VirtualScreenGlyphSize,
			np,
		)
		x += 4
	}
}

// updateCube builds the cube mesh for the given index.
func (w *toolBeltWidget) updateCube(i int) {
	if i < 0 || i >= len(w.items) {
		return
	}
	c := w.items[i]
	cube, _, f := c.Decompose()
	if cube == t.CubeRefInvalid || int(cube) >= len(mod.CubeDefs) {
		return
	}
	d := uint8(t.VirtualScreenGlyphSize * 2)
	w.cds[i].Mesh.Reset()
	c3d.AddCube(
		[3]uint8{0, 0, 0},
		[3]uint8{d, d, d},
		d, f,
		mod.CubeDefs[cube],
		w.cds[i].Mesh,
	)
}

// update updates the widget.
func (w *toolBeltWidget) update() {
	if w.dirty {
		w.updateMesh()
		w.dirty = false
	}
}

// input handles tool belt input.
func (w *toolBeltWidget) input() {
	s := w.selected
	h := w.hover
	l := (t.VirtualScreenGlyphsWide - 40) / 2
	t := t.VirtualScreenGlyphsHigh - 4
	pos := input.CursorGlyph
	if pos[0] >= l && pos[0] < l+40 && pos[1] >= t && pos[1] < t+4 {
		pos[0] -= l
		pos[1] -= t
		w.hover = pos[0] / 4
		if input.ButtonPressed(0) {
			w.selected = w.hover
		}
	} else {
		w.hover = -1
	}
	if input.WasPressed("tool-belt-1") {
		w.selected = 0
	}
	if input.WasPressed("tool-belt-2") {
		w.selected = 1
	}
	if input.WasPressed("tool-belt-3") {
		w.selected = 2
	}
	if input.WasPressed("tool-belt-4") {
		w.selected = 3
	}
	if input.WasPressed("tool-belt-5") {
		w.selected = 4
	}
	if input.WasPressed("tool-belt-6") {
		w.selected = 5
	}
	if input.WasPressed("tool-belt-7") {
		w.selected = 6
	}
	if input.WasPressed("tool-belt-8") {
		w.selected = 7
	}
	if input.WasPressed("tool-belt-9") {
		w.selected = 8
	}
	if input.WasPressed("tool-belt-0") {
		w.selected = 9
	}
	if s != w.selected || h != w.hover {
		w.dirty = true
	}
}

// setSelectedItem sets the contents of the selected slot on the tool belt.
func (w *toolBeltWidget) setSelectedItem(c t.Cell) {
	w.setItem(c, w.selected)
}

// setItem sets the contents of one slot on the tool belt.
func (w *toolBeltWidget) setItem(c t.Cell, i int) {
	if i < 0 || i >= len(w.items) {
		return
	}
	if w.items[i] == c {
		return
	}
	w.items[i] = c
	w.updateCube(i)
	w.dirty = true
}
