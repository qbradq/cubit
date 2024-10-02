package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
)

const minVoxFileVersion uint32 = 150

// Vox represents the contents of a .vox file in a generic format. The internal
// voxels are laid out with Z being the up axis and with the Y axis facing
// forward.
type Vox struct {
	Width, Height, Depth int
	Voxels               [][4]uint8
}

// Get implements the c3d.VoxelSource interface.
func (v *Vox) Get(x, y, z int) [4]uint8 {
	sx := x
	sy := v.Height - (z + 1)
	sz := y
	if sx < 0 || sx >= v.Width || sy < 0 || sy >= v.Height || sz < 0 ||
		sz >= v.Depth {
		return [4]uint8{0, 0, 0, 0}
	}
	return v.Voxels[sz*v.Width*v.Height+sy*v.Width+sx]
}

// Dimensions implements the c3d.VoxelSource interface.
func (v *Vox) Dimensions() (w, h, d int) {
	return v.Width, v.Depth, v.Height
}

// NewVoxFromReader returns a new Vox structure with the contents loaded from
// a MagicaVoxel .vox file.
func NewVoxFromReader(r io.Reader) (*Vox, error) {
	type chunk struct {
		ct        string
		data      []byte
		childData []byte
	}
	newChunk := func(r io.Reader) *chunk {
		ret := &chunk{}
		var buf = []byte{0, 0, 0, 0}
		if _, err := r.Read(buf); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
		}
		ret.ct = string(buf)
		nData := GetUint32(r)
		nChild := GetUint32(r)
		ret.data = make([]byte, nData)
		ret.childData = make([]byte, nChild)
		r.Read(ret.data)
		r.Read(ret.childData)
		return ret
	}
	var buf = []byte{0, 0, 0, 0}
	ret := &Vox{}
	r.Read(buf)
	version := GetUint32(r)
	if version < minVoxFileVersion {
		return nil, fmt.Errorf(
			".vox files must be at least version %d, found version %d",
			minVoxFileVersion, version)
	}
	main := newChunk(r)
	if main == nil || main.ct != "MAIN" {
		return nil, errors.New("did not find MAIN chunk in .vox file")
	}
	r = bytes.NewReader(main.childData)
	sizeSeen := false
	voxSeen := false
	var xyziBuf []uint8
	var pal = [256][4]uint8{}
	for {
		c := newChunk(r)
		if c == nil {
			break
		}
		switch c.ct {
		case "SIZE":
			if sizeSeen {
				return nil, errors.New("multiple SIZE chunks found")
			}
			sizeSeen = true
			cr := bytes.NewReader(c.data)
			ret.Width = int(GetUint32(cr))
			ret.Height = int(GetUint32(cr))
			ret.Depth = int(GetUint32(cr))
		case "XYZI":
			if voxSeen {
				return nil, errors.New("multiple XYZI chunks found")
			}
			voxSeen = true
			cr := bytes.NewReader(c.data)
			n := int(GetUint32(cr))
			xyziBuf = make([]uint8, ret.Width*ret.Height*ret.Depth)
			for i := 0; i < n; i++ {
				x := int(GetByte(cr))
				y := int(GetByte(cr))
				z := int(GetByte(cr))
				idx := uint8(GetByte(cr))
				xyziBuf[z*ret.Width*ret.Height+y*ret.Width+x] = idx
			}
		case "RGBA":
			cr := bytes.NewReader(c.data)
			for i := 0; i < 255; i++ {
				pal[i+1][0] = GetByte(cr)
				pal[i+1][1] = GetByte(cr)
				pal[i+1][2] = GetByte(cr)
				pal[i+1][3] = GetByte(cr)
				if pal[i+1][3] < 255 {
					log.Println(pal[i+1][3])
				}
			}
		}
	}
	// Compile voxel volume from the XYZI buffer
	ret.Voxels = make([][4]uint8, ret.Width*ret.Height*ret.Depth)
	for i, ci := range xyziBuf {
		ret.Voxels[i] = pal[ci]
	}
	return ret, nil
}
