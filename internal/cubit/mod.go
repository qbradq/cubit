package cubit

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/util"
)

// mustRead is a utility function for handling file reading from fs.FS
// implementors.
func mustRead(f fs.File, err error) []byte {
	if err != nil {
		panic(err)
	}
	d, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return d
}

// Mod manages a bundle of cubit engine content in separate external files.
type Mod struct {
	ID          string                          `json:"-"`           // Unique ID of the mod
	Name        string                          `json:"name"`        // Descriptive name of the mod
	Description string                          `json:"description"` // Description of the mod
	faceMap     map[c3d.FaceIndex]c3d.FaceIndex `json:"-"`           // Mapping of mod face indexes to final texture atlas indexes
	f           fs.FS
}

// newMod creates a new Mod object.
func newMod() *Mod {
	return &Mod{
		faceMap: map[c3d.FaceIndex]c3d.FaceIndex{},
	}
}

// Mods is the global map of all mods by ID
var Mods = map[string]*Mod{}

// loadModInfo loads the top-level info for a single mod.
func loadModInfo(name string) error {
	if _, duplicate := Mods[name]; duplicate {
		return fmt.Errorf("duplicate mod ID %s", name)
	}
	mod := newMod()
	f := os.DirFS(filepath.Join("mods", name))
	mod.f = f
	json.Unmarshal(mustRead(mod.f.Open("mod.json")), mod)
	mod.ID = name
	Mods[mod.ID] = mod
	return nil
}

// ReloadModInfo reloads all top-level info for all mods present.
func ReloadModInfo() error {
	Mods = map[string]*Mod{}
	cubeDefsById = map[string]*Cube{}
	cubeDefs = []*Cube{}
	voxIndex = map[string]*Vox{}
	dirs, err := os.ReadDir("mods")
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		// Ignore any top-level files that might be present
		if !dir.IsDir() {
			continue
		}
		// Load all directories as mods
		if err := loadModInfo(dir.Name()); err != nil {
			return err
		}
	}
	return nil
}

// LoadMods loads the named mods in order.
func LoadMods(mods ...string) error {
	ms := []*Mod{}
	for _, modName := range mods {
		mod, found := Mods[modName]
		if !found {
			return fmt.Errorf("mod %s not found", modName)
		}
		ms = append(ms, mod)
	}
	stage := func(fn func(*Mod) error) error {
		for _, mod := range ms {
			if err := fn(mod); err != nil {
				return err
			}
		}
		return nil
	}
	if err := stage(func(m *Mod) error { return m.loadFaces() }); err != nil {
		return err
	}
	if err := stage(func(m *Mod) error { return m.loadCubes() }); err != nil {
		return err
	}
	if err := stage(func(m *Mod) error { return m.loadVox() }); err != nil {
		return err
	}
	return nil
}

// wrap wraps an error for reporting.
func (m *Mod) wrap(where string, err error, args ...any) error {
	if len(args) > 0 {
		where = fmt.Sprintf(where, args)
	}
	return fmt.Errorf(
		"error loading mod %s: %s: %w", m.ID, where, err)
}

// loadFaces loads all faces for the mod.
func (m *Mod) loadFaces() error {
	return fs.WalkDir(m.f, "faces", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return m.wrap("walking faces directory, path=%s", err, path)
		}
		if len(path) < 1 {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ns := strings.ToLower(path)
		ns = filepath.Base(ns)
		ext := filepath.Ext(ns)
		ns = ns[:len(ns)-len(ext)]
		if ext != ".png" {
			return nil
		}
		v, err := strconv.ParseInt(ns, 0, 32)
		if err != nil {
			return m.wrap("parsing face index string %s", err, ns)
		}
		if v > 0xFF {
			return m.wrap("parsing face index",
				errors.New("face pages may range 0-255"))
		}
		f, err := m.f.Open(path)
		if err != nil {
			return m.wrap("reading face file", err)
		}
		if err := m.loadFacePage(uint8(v), f); err != nil {
			return err
		}
		return nil
	})
}

// loadFacePage loads the image as a face page, adding all non-empty (all
// transparent) faces into Faces.
func (m *Mod) loadFacePage(n uint8, r io.Reader) error {
	isEmpty := func(sub *image.RGBA) bool {
		for i := 0; i < len(sub.Pix); i += 4 {
			if sub.Pix[i+3] != 0 {
				return false
			}
		}
		return true
	}
	img, err := png.Decode(r)
	if err != nil {
		return m.wrap("decoding face %d", err, n)
	}
	if img.Bounds().Max.X != 256 || img.Bounds().Max.Y != 256 {
		return m.wrap("decoding face",
			errors.New("face page images must be 256x256 pixels in size"))
	}
	face := image.NewRGBA(image.Rect(0, 0, c3d.FaceDims, c3d.FaceDims))
	for fy := 0; fy < 8; fy++ {
		for fx := 0; fx < 8; fx++ {
			draw.Draw(
				face, image.Rect(0, 0, c3d.FaceDims, c3d.FaceDims),
				img, image.Pt(fx*c3d.FaceDims, fy*c3d.FaceDims),
				draw.Src)
			if isEmpty(face) {
				continue
			}
			m.faceMap[c3d.FaceIndexFromXYZ(fx, fy, int(n))] = Faces.AddFace(face)
		}
	}
	return nil
}

// loadCubes loads all cube definitions.
func (m *Mod) loadCubes() error {
	return fs.WalkDir(m.f, "cubes", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return m.wrap("walking cubes directory, path=%s", err, path)
		}
		if len(path) < 1 {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		ext = strings.ToLower(ext)
		if ext != ".json" {
			return nil
		}
		f, err := m.f.Open(path)
		if err != nil {
			return m.wrap("opening cube file %s", err, path)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return m.wrap("reading cube file %s", err, path)
		}
		cubes := map[string]*Cube{}
		if err := json.Unmarshal(b, &cubes); err != nil {
			return m.wrap("parsing cube file %s", err, path)
		}
		for k, cube := range cubes {
			cube.ID = "/" + m.ID + "/" + k
			// Convert mod-relative face references to global
			for i := range cube.Faces {
				fi, found := m.faceMap[cube.Faces[i]]
				if !found {
					cube.Faces[i] = c3d.FaceIndexInvalid
					continue
				}
				cube.Faces[i] = fi
			}
			if err := registerCube(cube); err != nil {
				return m.wrap("registering cube %s", err, cube.ID)
			}
		}
		return nil
	})
}

// loadVox loads all .vox models from the mod.
func (m *Mod) loadVox() error {
	return fs.WalkDir(m.f, "vox", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return m.wrap("walking vox directory, path=%s", err, path)
		}
		if len(path) < 1 {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		ext = strings.ToLower(ext)
		ns := filepath.Base(path)
		ns = ns[:len(ns)-len(ext)]
		if ext != ".vox" {
			return nil
		}
		modPath := "/" + m.ID + "/" + ns
		f, err := m.f.Open(path)
		if err != nil {
			return m.wrap("opening vox file %s", err, path)
		}
		if voxFile, err := util.NewVoxFromReader(f); err != nil {
			return m.wrap("processing vox file %s", err, path)
		} else {
			RegisterVox(modPath, NewVox(voxFile))
		}
		return nil
	})
}
