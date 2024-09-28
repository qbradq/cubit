package c3d

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/util"
)

// RGBAVolume is a 16x16x16 volume of RGB pixels.
type RGBAVolume struct {
	Position mgl32.Vec2  // Position of
	d        []byte      // Vertex buffer data
	count    int32       // Vertex count
	vao      uint32      // Vertex buffer array ID
	vbo      uint32      // Vertex buffer object ID
	vbuf     [6 * 6]byte // Vertex buffer
	dirty    bool        // If true, the volume has changed
	vd       []uint8     // Voxel data
	tid      uint32      // Texture ID of the volume texture
}

// NewRGBVolume creates a new RGB volume ready for use.
func NewRGBVolume() *RGBAVolume {
	// Init
	ret := &RGBAVolume{
		dirty: true,
		vd:    make([]uint8, 16*16*16*4),
	}
	// Texture array setup
	gl.CreateTextures(gl.TEXTURE_2D_ARRAY, 1, &ret.tid)
	gl.TextureParameteri(ret.tid, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(ret.tid, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(ret.tid, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(ret.tid, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TextureParameteri(ret.tid, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TextureStorage3D(ret.tid, 1, gl.RGBA8, 16, 16, 16)
	// VAO setup
	gl.CreateBuffers(1, &ret.vbo)
	gl.CreateVertexArrays(1, &ret.vao)
	gl.VertexArrayVertexBuffer(ret.vao, 0, ret.vbo, 0, 3*1+3*1)
	gl.EnableVertexArrayAttrib(ret.vao, 0)
	gl.EnableVertexArrayAttrib(ret.vao, 1)
	gl.VertexArrayAttribIFormat(ret.vao, 0, 3, gl.UNSIGNED_BYTE, 0)
	gl.VertexArrayAttribIFormat(ret.vao, 1, 3, gl.UNSIGNED_BYTE, 3*1)
	gl.VertexArrayAttribBinding(ret.vao, 0, 0)
	gl.VertexArrayAttribBinding(ret.vao, 1, 0)
	return ret
}

// SetData sets the raw voxel data for the volume.
func (v *RGBAVolume) SetData(vox *util.Vox) {
	if vox.Width != 16 || vox.Height != 16 || vox.Depth != 16 {
		panic("RGBVolume.SetData called with wrong dimensions")
	}
	for iz := 0; iz < 16; iz++ {
		for iy := 0; iy < 16; iy++ {
			for ix := 0; ix < 16; ix++ {
				v.Set(ix, iy, iz, vox.Get(ix, iy, iz))
			}
		}
	}
	v.dirty = true
}

// Set sets the voxel at the given location.
func (v *RGBAVolume) Set(x, y, z int, c [4]uint8) {
	idx := (z*16*16 + y*16 + x) * 4
	v.vd[idx+0] = c[0]
	v.vd[idx+1] = c[1]
	v.vd[idx+2] = c[2]
	v.vd[idx+3] = c[3]
}

// layer lays down a single layer into the volume's quad lattice.
func (v *RGBAVolume) layer(f Facing, layer int) {
	d := v.vbuf[:]
	o := CubeFacingOffsets[f]
	a := layerAxis[f]
	m := layerAxisMask[f]
	x := a[0] * layer
	y := a[1] * layer
	z := a[2] * layer
	lm := facingXYZLayerMask[f]
	fo := facingXYZOffsets[f]
	for i := 0; i < 6; i++ {
		d[i*6+0] = byte((x + o[i][0]) * m[0])
		d[i*6+1] = byte((y + o[i][1]) * m[1])
		d[i*6+2] = byte((z + o[i][2]) * m[2])
		d[i*6+3] = byte(fo[i][0] + layer*lm[0])
		d[i*6+4] = byte(fo[i][1] + layer*lm[1])
		d[i*6+5] = byte(fo[i][2] + layer*lm[2])
	}
	v.d = append(v.d, d[:]...)
	v.count += 6
}

// update builds the optimized voxel slice lattice for rendering.
func (v *RGBAVolume) update() {
	face := func(x, y, z int, f Facing) bool {
		if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
			return false
		}
		if v.vd[(z*16*16+y*16+x)*4+3] < 255 {
			return false
		}
		dx := FacingOffsets[f][0] + x
		dy := FacingOffsets[f][1] + y
		dz := FacingOffsets[f][2] + z
		if dx < 0 || dx > 15 || dy < 0 || dy > 15 || dz < 0 || dz > 15 {
			return false
		}
		return v.vd[(dz*16*16+dy*16+dx)*4+3] == 255
	}
	min := [3]int{0, 0, 0}
	max := [3]int{0, 0, 0}
	for f := North; f <= Bottom; f++ {
		fplm := facingXYZLayerMask[f]
		fsm := facingSweepMax[f]
		for l := 0; l < 16; l++ {
			min[0] = fplm[0] * l
			min[1] = fplm[1] * l
			min[2] = fplm[2] * l
			max[0] = min[0] + fsm[0]
			max[1] = min[1] + fsm[1]
			max[2] = min[2] + fsm[2]
			for iy := min[1]; iy <= max[1]; iy++ {
				for iz := min[2]; iz <= max[2]; iz++ {
					for ix := min[0]; ix <= max[0]; ix++ {
						if face(ix, iy, iz, f) {
							goto doLayer
						}
					}
				}
			}
			continue
		doLayer:
			v.layer(f, l)
		}
	}
}

func (v *RGBAVolume) draw() {
	if v.dirty {
		v.update()
		if len(v.d) > 0 {
			gl.NamedBufferData(v.vbo, len(v.d), gl.Ptr(v.d), gl.STATIC_DRAW)
		}
		gl.TextureSubImage3D(v.tid, 0, 0, 0, 0, 16, 16, 16, gl.RGBA,
			gl.UNSIGNED_BYTE, gl.Ptr(v.vd))
		v.dirty = false
	}
	if v.vao != invalidVAO {
		gl.BindTextureUnit(0, v.tid)
		gl.BindVertexArray(v.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, v.count)
	}
}
