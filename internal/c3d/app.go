package c3d

import (
	"fmt"

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
	CursorVisible      bool            // If true, draw the cursor
	CrosshairVisible   bool            // If true, draw the crosshair
	DebugTextVisible   bool            // If true, draw the debug text
	ChunkBoundsVisible bool            // If true, draws the chunk bounds
	cubeMeshes         []cubeMeshRef   // List of cube meshes to draw along with orientation
	voxelMeshes        []voxelMeshRef  // List of voxel meshes to draw along with orientation
	uiMeshes           []*UIMesh       // List of UI meshes to draw
	cursor             *UIMesh         // Cursor mesh
	crosshair          *UIMesh         // Crosshair mesh
	axis               *LineMesh       // Debug axis indicator
	pRGBFB             *program        // RGB with no lighting
	pVoxelMesh         *program        // RGB
	pCubeMesh          *program        // Face atlas texturing
	pText              *program        // Text rendering
	pUI                *program        // UI tile rendering
	faces              *FaceAtlas      // Face atlas to use for cube mesh rendering
	tiles              *FaceAtlas      // Face atlas to use for ui tile rendering
	fm                 *fontManager    // Font manager for the application
	debugLines         []ColoredString // Lines for the debug messages
	debugText          *TextMesh       // Debug text
}

// NewApp constructs a new App object with the given resources ready to draw.
func NewApp(faces *FaceAtlas, tiles *FaceAtlas) (*App, error) {
	var err error
	ret := &App{
		faces: faces,
		tiles: tiles,
	}
	// rgb_fullbright.glsl
	ret.pRGBFB, err = loadProgram("rgb-fullbright")
	if err != nil {
		return nil, err
	}
	// rgb.glsl
	ret.pVoxelMesh, err = loadProgram("voxel-mesh")
	if err != nil {
		return nil, err
	}
	// cube_mesh.glsl
	ret.pCubeMesh, err = loadProgram("cube-mesh")
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
	ret.debugText = newTextMesh(ret.fm, ret.pText)
	// ui.glsl
	ret.pUI, err = loadProgram("ui")
	if err != nil {
		return nil, err
	}
	ret.tiles.upload(ret.pUI)
	ret.tiles.freeMemory()
	// Axis indicator
	ret.axis = NewLineMesh()
	ret.axis.Line(mgl32.Vec3{}, mgl32.Vec3{1, 0, 0}, [4]uint8{255, 0, 0, 255})
	ret.axis.Line(mgl32.Vec3{}, mgl32.Vec3{0, 1, 0}, [4]uint8{0, 255, 0, 255})
	ret.axis.Line(mgl32.Vec3{}, mgl32.Vec3{0, 0, 1}, [4]uint8{0, 0, 255, 255})
	return ret, nil
}

// Delete removes all memory and GPU resources managed by the app.
func (a *App) Delete() {
	a.pVoxelMesh.delete()
	a.pRGBFB.delete()
	a.pCubeMesh.delete()
	a.pText.delete()
}

// AddDebugLine sets the debug text drawn in the bottom-left.
func (a *App) AddDebugLine(c [3]uint8, f string, args ...any) {
	a.debugLines = append(a.debugLines, ColoredString{
		String: fmt.Sprintf(f, args...),
		Color:  c,
	})
}

// updateDebugText updates the debug text mesh.
func (a *App) updateDebugText() {
	a.debugText.Reset()
	y := VirtualScreenHeight - len(a.debugLines)*LineSpacingVS
	for _, line := range a.debugLines {
		a.debugText.Print(0, y, line.Color, line.String)
		y += LineSpacingVS
	}
}

// SetCursor sets the UI tile to use as the mouse cursor.
func (a *App) SetCursor(f FaceIndex, l uint16) {
	a.cursor = a.NewUIMesh()
	a.cursor.Scaled(0, 0, vsGlyphWidth*2, vsGlyphWidth*2, f)
	a.cursor.Layer = l
}

// SetCursorPosition sets the cursor position.
func (a *App) SetCursorPosition(p mgl32.Vec2) {
	if a.cursor == nil {
		return
	}
	a.cursor.Position = p
}

// SetCrosshair sets the UI tile to use as the 3D crosshair.
func (a *App) SetCrosshair(f FaceIndex, l uint16) {
	a.crosshair = a.NewUIMesh()
	a.crosshair.Scaled(0, 0, vsGlyphWidth*2, vsGlyphWidth*2, f)
	a.crosshair.Position[0] = float32(VirtualScreenWidth-vsGlyphWidth*2) / 2
	a.crosshair.Position[1] = float32(VirtualScreenHeight-vsGlyphWidth*2) / 2
	a.crosshair.Layer = l
}

// AddVoxelMesh adds the voxel mesh to the list to render. The value of o is
// copied internally, so o may be reused after the call to AddVoxelMesh.
func (a *App) AddVoxelMesh(m *VoxelMesh, o *Orientation) {
	a.voxelMeshes = append(a.voxelMeshes, voxelMeshRef{m: m, o: *o})
}

// RemoveVoxelMesh removes the passed voxel mesh from the renderer.
func (a *App) RemoveVoxelMesh(m *VoxelMesh) {
	for i := 0; i < len(a.voxelMeshes); i++ {
		if a.voxelMeshes[i].m == m {
			a.voxelMeshes[i] = a.voxelMeshes[len(a.voxelMeshes)-1]
			a.voxelMeshes = a.voxelMeshes[:len(a.voxelMeshes)]
			return
		}
	}
}

// AddCubeMesh adds the cube mesh to the list to render. The value of o is
// copied internally, so o may be reused after the call to AddCubeMesh.
func (a *App) AddCubeMesh(m *CubeMesh, o *Orientation) {
	a.cubeMeshes = append(a.cubeMeshes, cubeMeshRef{m: m, o: *o})
}

// RemoveCubeMesh removes the passed cube mesh from the renderer.
func (a *App) RemoveCubeMesh(m *CubeMesh) {
	for i := 0; i < len(a.cubeMeshes); i++ {
		if a.cubeMeshes[i].m == m {
			a.cubeMeshes[i] = a.cubeMeshes[len(a.cubeMeshes)-1]
			a.cubeMeshes = a.cubeMeshes[:len(a.cubeMeshes)]
			return
		}
	}
}

// NewUIMesh creates and returns a new UIMesh.
func (a *App) NewUIMesh() *UIMesh {
	return newUIMesh(a.fm, a.pUI, a.pText)
}

// AddUIMesh adds the UI mesh to the list to render.
func (a *App) AddUIMesh(m *UIMesh) {
	a.uiMeshes = append(a.uiMeshes, m)
}

// RemoveUIMesh removes the passed cube mesh from the renderer.
func (a *App) RemoveUIMesh(m *UIMesh) {
	for i := 0; i < len(a.uiMeshes); i++ {
		if a.uiMeshes[i] == m {
			a.uiMeshes[i] = a.uiMeshes[len(a.uiMeshes)-1]
			a.uiMeshes = a.uiMeshes[:len(a.uiMeshes)]
			return
		}
	}
}

// Draw draws everything with the given camera for 3D space..
func (a *App) Draw(c *Camera) {
	// Variable setup
	pMat := mgl32.Perspective(mgl32.DegToRad(60),
		float32(VirtualScreenWidth)/float32(VirtualScreenHeight), 0.1, 1000.0)
	// Frame setup
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// Draw wire frames
	a.pRGBFB.use()
	gl.UniformMatrix4fv(int32(a.pRGBFB.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	if a.ChunkBoundsVisible {
		for _, m := range a.cubeMeshes {
			mt := c.TransformMatrix()
			gl.UniformMatrix4fv(a.pRGBFB.uni("uModelViewMatrix"), 1, false,
				&mt[0])
			m.m.drawAABB(a.pRGBFB)
		}
	}
	mt := c.TransformMatrix().Mul4(a.axis.Orientation.TransformMatrix())
	gl.UniformMatrix4fv(a.pRGBFB.uni("uModelViewMatrix"), 1, false, &mt[0])
	a.axis.draw(a.pRGBFB)
	// Draw chunks
	a.pCubeMesh.use()
	gl.UniformMatrix4fv(int32(a.pCubeMesh.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	a.faces.bind(a.pCubeMesh)
	for _, m := range a.cubeMeshes {
		mt := c.TransformMatrix().Mul4(m.o.TransformMatrix())
		gl.UniformMatrix4fv(a.pCubeMesh.uni("uModelViewMatrix"), 1, false,
			&mt[0])
		m.m.draw(a.pCubeMesh)
	}
	// Draw voxel models
	a.pVoxelMesh.use()
	gl.UniformMatrix4fv(int32(a.pVoxelMesh.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	for _, v := range a.voxelMeshes {
		mt := c.TransformMatrix().Mul4(v.o.TransformMatrix())
		gl.UniformMatrix4fv(a.pVoxelMesh.uni("uModelViewMatrix"), 1, false, &mt[0])
		v.m.draw(a.pVoxelMesh)
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
			float32(m.Layer)/0xFFFF)
		m.draw()
	}
	// Draw common screen components
	if a.CrosshairVisible && a.crosshair != nil {
		m := a.crosshair
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			float32(m.Layer)/0xFFFF)
		m.draw()
	}
	if a.CursorVisible && a.cursor != nil {
		m := a.cursor
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			float32(m.Layer)/0xFFFF)
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
			float32(m.Layer)/0xFFFF)
		m.Text.draw()
	}
	// Draw debug text on top of everything
	if a.DebugTextVisible && a.debugText != nil {
		a.updateDebugText()
		m := a.debugText
		gl.Uniform3f(a.pUI.uni("uPosition"), 0, 0, 1.0)
		m.draw()
	}
	a.debugLines = a.debugLines[:0]
}
