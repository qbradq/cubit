package client

import (
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/mod"
	"github.com/qbradq/cubit/internal/t"
)

// paletteWidget implements
type paletteWidget struct {
	baseWidget
	items  []t.Cell                       // Cell storage
	cds    []*c3d.CubeMeshDrawDescriptor  // Cube dds
	voxDDs []*c3d.VoxelMeshDrawDescriptor // Vox dds
	hover  int                            // Index of the currently hovered grid box
	dirty  bool                           // If true the UI needs to be layed out again
}

// newPaletteWidget creates a new palette widget for use.
func newPaletteWidget() *paletteWidget {
	ret := &paletteWidget{
		baseWidget: *newBaseWidget(app),
		dirty:      true,
	}
	ret.UIMesh.Layer = layerPalette
	ret.items = make([]t.Cell, len(mod.CubeDefs)+len(mod.VoxDefs))
	ret.cds = make([]*c3d.CubeMeshDrawDescriptor, len(ret.items))
	ret.voxDDs = make([]*c3d.VoxelMeshDrawDescriptor, len(ret.items))
	i := 0
	for _, d := range mod.CubeDefs {
		x := t.VirtualScreenGlyphsWide - 5*4
		x += (i%5)*4 + 2
		y := (i/5)*4 + 2
		ret.items[i] = t.CellForCube(d.Ref, t.North)
		ret.cds[i] = ret.CubeMeshIcon(
			x*t.VirtualScreenGlyphSize,
			y*t.VirtualScreenGlyphSize,
			8, 8, 8, d, t.North,
			mod.CubeDefs,
		)
		i++
	}
	return ret
}

// layout lays out the UI mesh.
func (w *paletteWidget) layout() {
	w.Reset(true, false)
	sx := (t.VirtualScreenGlyphsWide - 5*4) * t.VirtualScreenGlyphSize
	y := 0
	d := t.VirtualScreenGlyphSize * 4
	i := 0
	for iy := 0; iy < 10; iy++ {
		x := sx
		for ix := 0; ix < 5; ix++ {
			np := npWindow
			if i == w.hover {
				np = npHighlight
			}
			w.NinePatch(x, y, d, d, np)
			x += d
			i++
		}
		y += d
	}

}

// update processes periodic updates to the widget.
func (w *paletteWidget) update() {
	if w.dirty {
		w.layout()
		w.dirty = false
	}
	w.Hidden = !input.InUIMode
}

// input handles input for the widget.
func (w *paletteWidget) input() {
	h := w.hover
	pos := input.CursorGlyph
	l := t.VirtualScreenGlyphsWide - 20
	if pos[0] >= l && pos[1] < 40 {
		x := (pos[0] - l) / 4
		y := pos[1] / 4
		w.hover = y*5 + x
		if input.ButtonPressed(0) {
			c := t.CellInvalid
			i := w.hover
			if i < len(mod.CubeDefs) {
				c = t.CellForCube(mod.CubeDefs[i].Ref, t.North)
			}
			toolBelt.setSelectedItem(c)
		}
	}
	if h != w.hover {
		w.dirty = true
	}
}
