package c3d

import "github.com/qbradq/cubit/internal/t"

// Mesh is implemented by all mesh types in c3d.
type Mesh[T any] interface {
	// draw asks the mesh to draw itself.
	draw(p *program)
	// vert adds a vertex to the mesh.
	vert(x, y, z, u, v uint8, i int, c T, f t.Facing)
	// Reset resets the vertex data of the mesh.
	Reset()
}

// AddFace adds the given face with the given position and dimensions to the
// mesh.
func AddFace[T any](p, d [3]uint8, f t.Facing, c T, m Mesh[T]) {
	d[0] -= 1
	d[1] -= 1
	d[2] -= 1
	switch f {
	case t.North:
		u := d[0] + 1
		v := d[1] + 1
		m.vert(p[0]+d[0]+1, p[1]+d[1]+1, p[2], 0, 0, 0, c, f) // TL
		m.vert(p[0], p[1]+d[1]+1, p[2], u, 0, 1, c, f)        // TR
		m.vert(p[0]+d[0]+1, p[1], p[2], 0, v, 2, c, f)        // BL
		m.vert(p[0]+d[0]+1, p[1], p[2], 0, v, 2, c, f)        // BL
		m.vert(p[0], p[1]+d[1]+1, p[2], u, 0, 1, c, f)        // TR
		m.vert(p[0], p[1], p[2], u, v, 3, c, f)               // BR
	case t.South:
		u := d[0] + 1
		v := d[1] + 1
		m.vert(p[0], p[1]+d[1]+1, p[2]+1, 0, 0, 0, c, f)        // TL
		m.vert(p[0]+d[0]+1, p[1]+d[1]+1, p[2]+1, u, 0, 1, c, f) // TR
		m.vert(p[0], p[1], p[2]+1, 0, v, 2, c, f)               // BL
		m.vert(p[0], p[1], p[2]+1, 0, v, 2, c, f)               // BL
		m.vert(p[0]+d[0]+1, p[1]+d[1]+1, p[2]+1, u, 0, 1, c, f) // TR
		m.vert(p[0]+d[0]+1, p[1], p[2]+1, u, v, 3, c, f)        // BR
	case t.East:
		u := d[2] + 1
		v := d[1] + 1
		m.vert(p[0]+1, p[1]+d[1]+1, p[2]+d[2]+1, 0, 0, 0, c, f) // TL
		m.vert(p[0]+1, p[1]+d[1]+1, p[2], u, 0, 1, c, f)        // TR
		m.vert(p[0]+1, p[1], p[2]+d[2]+1, 0, v, 2, c, f)        // BL
		m.vert(p[0]+1, p[1], p[2]+d[2]+1, 0, v, 2, c, f)        // BL
		m.vert(p[0]+1, p[1]+d[1]+1, p[2], u, 0, 1, c, f)        // TR
		m.vert(p[0]+1, p[1], p[2], u, v, 3, c, f)               // BR
	case t.West:
		u := d[2] + 1
		v := d[1] + 1
		m.vert(p[0], p[1]+d[1]+1, p[2], 0, 0, 0, c, f)        // TL
		m.vert(p[0], p[1]+d[1]+1, p[2]+d[2]+1, u, 0, 1, c, f) // TR
		m.vert(p[0], p[1], p[2], 0, v, 2, c, f)               // BL
		m.vert(p[0], p[1], p[2], 0, v, 2, c, f)               // BL
		m.vert(p[0], p[1]+d[1]+1, p[2]+d[2]+1, u, 0, 1, c, f) // TR
		m.vert(p[0], p[1], p[2]+d[2]+1, u, v, 3, c, f)        // BR
	case t.Top:
		u := d[0] + 1
		v := d[2] + 1
		m.vert(p[0], p[1]+1, p[2], 0, 0, 0, c, f)               // TL
		m.vert(p[0]+d[0]+1, p[1]+1, p[2], u, 0, 1, c, f)        // TR
		m.vert(p[0], p[1]+1, p[2]+d[2]+1, 0, v, 2, c, f)        // BL
		m.vert(p[0], p[1]+1, p[2]+d[2]+1, 0, v, 2, c, f)        // BL
		m.vert(p[0]+d[0]+1, p[1]+1, p[2], u, 0, 1, c, f)        // TR
		m.vert(p[0]+d[0]+1, p[1]+1, p[2]+d[2]+1, u, v, 3, c, f) // BR
	case t.Bottom:
		u := d[0] + 1
		v := d[2] + 1
		m.vert(p[0]+d[0]+1, p[1], p[2], 0, 0, 0, c, f)        // TL
		m.vert(p[0], p[1], p[2], u, 0, 1, c, f)               // TR
		m.vert(p[0]+d[0]+1, p[1], p[2]+d[2]+1, 0, v, 2, c, f) // BL
		m.vert(p[0]+d[0]+1, p[1], p[2]+d[2]+1, 0, v, 2, c, f) // BL
		m.vert(p[0], p[1], p[2], u, 0, 1, c, f)               // TR
		m.vert(p[0], p[1], p[2]+d[2]+1, u, v, 3, c, f)        // BR
	}
}

// AddCube adds all of the faces of the given cube definition.
func AddCube(p, d [3]uint8, f t.Facing, c *t.Cube, m *CubeMesh) {
	fn := func(face t.Facing) t.Facing {
		return t.FacingMap[f][face]
	}
	n := [3]uint8{p[0], p[1], p[2]}
	s := [3]uint8{p[0], p[1], p[2] + d[2] - 1}
	e := [3]uint8{p[0] + d[0] - 1, p[1], p[2]}
	w := [3]uint8{p[0], p[1], p[2]}
	top := [3]uint8{p[0], p[1] + d[1] - 1, p[2]}
	b := p
	AddFace(n, d, fn(t.North), t.CellForCube(c.Ref, f), m)
	AddFace(s, d, fn(t.South), t.CellForCube(c.Ref, f), m)
	AddFace(e, d, fn(t.East), t.CellForCube(c.Ref, f), m)
	AddFace(w, d, fn(t.West), t.CellForCube(c.Ref, f), m)
	AddFace(top, d, fn(t.Top), t.CellForCube(c.Ref, f), m)
	AddFace(b, d, fn(t.Bottom), t.CellForCube(c.Ref, f), m)
}
