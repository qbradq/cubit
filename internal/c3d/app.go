package c3d

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// voxelMeshRef refers to a voxel mesh pointer and the orientation it should
// be drawn with.
type voxelMeshRef struct {
	m *VoxelMesh
	o Orientation
}

// cubeMeshRef refers to a cube mesh pointer and the orientation it should
// be drawn with.
type cubeMeshRef struct {
	m *CubeMesh
	o Orientation
}

// App manages a set of drawable objects and the shader programs used to draw
// them.
type App struct {
	cubeMeshes  []cubeMeshRef  // List of cube meshes to draw along with orientation
	voxelMeshes []voxelMeshRef // List of voxel meshes to draw along with orientation
	uiMeshes    []*UIMesh      // List of UI meshes to draw
	axis        *AxisIndicator // Debug axis indicator
	pRGBFB      *program       // RGB with no lighting
	pRGB        *program       // RGB
	pCubeMesh   *program       // Face atlas texturing
	pText       *program       // Text rendering
	pUI         *program       // UI tile rendering
	faces       *FaceAtlas     // Face atlas to use for cube mesh rendering
	tiles       *FaceAtlas     // Face atlas to use for ui tile rendering
	fm          *fontManager   // Font manager for the application
}

// NewApp constructs a new App object with the given resources ready to draw.
func NewApp(faces *FaceAtlas, tiles *FaceAtlas) (*App, error) {
	var err error
	ret := &App{
		faces: faces,
		tiles: tiles,
	}
	// rgb_fullbright.glsl
	ret.pRGBFB, err = loadProgram("rgb_fullbright")
	if err != nil {
		return nil, err
	}
	ret.axis = newAxisIndicator(mgl32.Vec3{0, 0, 0}, ret.pRGBFB)
	// rgb.glsl
	ret.pRGB, err = loadProgram("rgb")
	if err != nil {
		return nil, err
	}
	// cube_mesh.glsl
	ret.pCubeMesh, err = loadProgram("cube_mesh")
	if err != nil {
		return nil, err
	}
	ret.faces.upload(ret.pCubeMesh)
	ret.faces.freeMemory()
	// text.glsl
	ret.pText, err = loadProgram("text")
	if err != nil {
		return nil, err
	}
	ret.fm = newFontManager(ret.pText)
	// ui.glsl
	ret.pUI, err = loadProgram("ui")
	if err != nil {
		return nil, err
	}
	ret.tiles.upload(ret.pUI)
	ret.tiles.freeMemory()
	return ret, nil
}

// Delete removes all memory and GPU resources managed by the app.
func (a *App) Delete() {
	a.pRGB.delete()
	a.pRGBFB.delete()
	a.pCubeMesh.delete()
	a.pText.delete()
}

// AddVoxelMesh adds the voxel mesh to the list to render. The value of o is
// copied internally, so o may be reused after the call to AddVoxelMesh.
func (a *App) AddVoxelMesh(m *VoxelMesh, o *Orientation) {
	a.voxelMeshes = append(a.voxelMeshes, voxelMeshRef{m: m, o: *o})
}

// AddCubeMesh adds the cube mesh to the list to render. The value of o is
// copied internally, so o may be reused after the call to AddCubeMesh.
func (a *App) AddCubeMesh(m *CubeMesh, o *Orientation) {
	a.cubeMeshes = append(a.cubeMeshes, cubeMeshRef{m: m, o: *o})
}

// NewUIMesh creates and returns a new UIMesh.
func (a *App) NewUIMesh() *UIMesh {
	return newUIMesh(a.fm, a.pUI)
}

// AddUIMesh adds the UI mesh to the list to render.
func (a *App) AddUIMesh(m *UIMesh) {
	a.uiMeshes = append(a.uiMeshes, m)
}

// Draw draws everything with the given camera for 3D space..
func (a *App) Draw(c *Camera) {
	// Variable setup
	pMat := mgl32.Perspective(mgl32.DegToRad(60),
		float32(VirtualScreenWidth)/float32(VirtualScreenHeight), 0.1, 1000.0)
	// Frame setup
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// TODO REMOVE Draw debug indicators
	a.pRGBFB.use()
	gl.UniformMatrix4fv(int32(a.pRGBFB.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	mt := c.TransformMatrix().Mul4(a.axis.o.TransformMatrix())
	gl.UniformMatrix4fv(a.pRGBFB.uni("uModelViewMatrix"), 1, false, &mt[0])
	a.axis.draw()
	// Draw chunks
	a.pCubeMesh.use()
	gl.UniformMatrix4fv(int32(a.pCubeMesh.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	a.faces.bind(a.pCubeMesh)
	for _, m := range a.cubeMeshes {
		mt := c.TransformMatrix().Mul4(m.o.TransformMatrix())
		gl.UniformMatrix4fv(a.pCubeMesh.uni("uModelViewMatrix"), 1, false,
			&mt[0])
		nt := mt.Inv().Transpose()
		gl.UniformMatrix4fv(a.pCubeMesh.uni("uNormalMatrix"), 1, false, &nt[0])
		m.m.draw(a.pCubeMesh)
	}
	// Draw voxel models
	a.pRGB.use()
	gl.UniformMatrix4fv(int32(a.pRGB.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	for _, v := range a.voxelMeshes {
		mt := c.TransformMatrix().Mul4(v.o.TransformMatrix())
		gl.UniformMatrix4fv(a.pRGB.uni("uModelViewMatrix"), 1, false, &mt[0])
		nt := mt.Inv().Transpose()
		gl.UniformMatrix4fv(a.pRGB.uni("uNormalMatrix"), 1, false, &nt[0])
		v.m.draw(a.pRGB)
	}
	// Draw UI elements, tiles layer
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	a.pUI.use()
	pMat = mgl32.Ortho2D(0, float32(VirtualScreenWidth), 0,
		float32(VirtualScreenHeight))
	gl.UniformMatrix4fv(a.pUI.uni("uProjectionMatrix"), 1, false, &pMat[0])
	a.tiles.bind(a.pUI)
	for _, m := range a.uiMeshes {
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			m.Layer)
		m.draw()
	}
	// Draw UI elements, text layer
	a.pText.use()
	if a.fm.imgDirty {
		a.fm.updateAtlasTexture()
	}
	pMat = mgl32.Ortho2D(0, float32(VirtualScreenWidth), 0,
		float32(VirtualScreenHeight))
	gl.UniformMatrix4fv(a.pText.uni("uProjectionMatrix"), 1, false, &pMat[0])
	a.fm.bind(a.pText)
	for _, m := range a.uiMeshes {
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			m.Layer)
		m.Text.draw()
	}
}
