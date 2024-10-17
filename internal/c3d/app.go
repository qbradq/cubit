package c3d

import (
	"fmt"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/t"
)

// App manages a set of drawable objects and the shader programs used to draw
// them.
type App struct {
	CursorVisible     bool                      // If true, draw the cursor
	CrosshairVisible  bool                      // If true, draw the crosshair
	DebugTextVisible  bool                      // If true, draw the debug text
	WireFramesVisible bool                      // If true, draws wire frames
	chunkDDs          []*ChunkDrawDescriptor    // List of chunks to draw
	modelDDs          []*ModelDrawDescriptor    // List of models to draw
	lineDDs           []*LineMeshDrawDescriptor // List of line meshes to draw
	uiMeshes          []*UIMesh                 // List of UI meshes to draw
	cursor            *UIMesh                   // Cursor mesh
	crosshair         *UIMesh                   // Crosshair mesh
	axis              *LineMesh                 // Debug axis indicator
	chunkBounds       AABB                      // Cached chunk bounds wire frame
	pWireFrame        *program                  // RGB with no lighting
	pVoxelMesh        *program                  // RGB voxel meshes
	pModelMesh        *program                  // RGB voxel meshes rigged for animation
	pCubeMesh         *program                  // Face atlas texturing
	pText             *program                  // Text rendering
	pUI               *program                  // UI tile rendering
	pCubeMeshIcon     *program                  // Cube mesh icon rendering
	faces             *FaceAtlas                // Face atlas to use for cube mesh rendering
	tiles             *FaceAtlas                // Face atlas to use for ui tile rendering
	fm                *fontManager              // Font manager for the application
	debugLines        []ColoredString           // Lines for the debug messages
	debugText         *TextMesh                 // Debug text
}

// NewApp constructs a new App object with the given resources ready to draw.
func NewApp(faces *FaceAtlas, tiles *FaceAtlas) (*App, error) {
	var err error
	ret := &App{
		faces: faces,
		tiles: tiles,
	}
	// wireframe.glsl
	ret.pWireFrame, err = loadProgram("wireframe")
	if err != nil {
		return nil, err
	}
	// voxel-mesh.glsl
	ret.pVoxelMesh, err = loadProgram("voxel-mesh")
	if err != nil {
		return nil, err
	}
	gl.Uniform1fv(ret.pVoxelMesh.uni("uLightLevels"),
		int32(len(voxelLightLevels)), &voxelLightLevels[0])
	// model-mesh.glsl
	ret.pModelMesh, err = loadProgram("model-mesh")
	if err != nil {
		return nil, err
	}
	// cube-mesh.glsl
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
	// cube-mesh-icon.glsl
	ret.pCubeMeshIcon, err = loadProgram("cube-mesh-icon")
	if err != nil {
		return nil, err
	}
	// Internal meshes
	ret.genAxis()
	ret.chunkBounds = AABB{
		Bounds: t.AABB{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{16, 16, 16},
		},
	}
	return ret, nil
}

func (a *App) genAxis() {
	a.axis = NewLineMesh()
	c := [4]uint8{255, 255, 255, 255}
	for i := float32(0); i <= 16; i++ {
		a.axis.Line(mgl32.Vec3{i, 0, 0}, mgl32.Vec3{i, 16, 0}, c)
		a.axis.Line(mgl32.Vec3{i, 0, 0}, mgl32.Vec3{i, 0, 16}, c)
		a.axis.Line(mgl32.Vec3{0, i, 0}, mgl32.Vec3{16, i, 0}, c)
		a.axis.Line(mgl32.Vec3{0, i, 0}, mgl32.Vec3{0, i, 16}, c)
		a.axis.Line(mgl32.Vec3{0, 0, i}, mgl32.Vec3{0, 16, i}, c)
		a.axis.Line(mgl32.Vec3{0, 0, i}, mgl32.Vec3{16, 0, i}, c)
	}
}

// Delete removes all memory and GPU resources managed by the app.
func (a *App) Delete() {
	a.pVoxelMesh.delete()
	a.pWireFrame.delete()
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
	y := t.VirtualScreenHeight - len(a.debugLines)*t.LineSpacingVS
	for _, line := range a.debugLines {
		a.debugText.Print(0, y, line.Color, line.String)
		y += t.LineSpacingVS
	}
}

// SetCursor sets the UI tile to use as the mouse cursor.
func (a *App) SetCursor(f t.FaceIndex, l uint16) {
	a.cursor = a.NewUIMesh()
	a.cursor.Scaled(0, 0, t.VSGlyphWidth*2, t.VSGlyphWidth*2, f)
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
func (a *App) SetCrosshair(f t.FaceIndex, l uint16) {
	a.crosshair = a.NewUIMesh()
	a.crosshair.Scaled(0, 0, t.VSGlyphWidth*2, t.VSGlyphWidth*2, f)
	a.crosshair.Position[0] =
		float32(t.VirtualScreenWidth-t.VSGlyphWidth*2) / 2
	a.crosshair.Position[1] =
		float32(t.VirtualScreenHeight-t.VSGlyphWidth*2) / 2
	a.crosshair.Layer = l
}

// NewUIMesh creates and returns a new UIMesh.
func (a *App) NewUIMesh() *UIMesh {
	return newUIMesh(a.fm, a.pText)
}

// AddUIMesh adds the UI mesh to the list to render.
func (a *App) AddUIMesh(m *UIMesh) {
	a.uiMeshes = append(a.uiMeshes, m)
}

// RemoveUIMesh removes the passed UI mesh from the renderer.
func (a *App) RemoveUIMesh(m *UIMesh) {
	for i := 0; i < len(a.uiMeshes); i++ {
		if a.uiMeshes[i] == m {
			a.uiMeshes[i] = a.uiMeshes[len(a.uiMeshes)-1]
			a.uiMeshes = a.uiMeshes[:len(a.uiMeshes)]
			return
		}
	}
}

// AddChunkDD adds the chunk draw descriptor to the list to draw.
func (a *App) AddChunkDD(d *ChunkDrawDescriptor) {
	a.chunkDDs = append(a.chunkDDs, d)
}

// RemoveChunkDD removes the chunk draw descriptor by ID.
func (a *App) RemoveChunkDD(id uint32) {
	for i := 0; i < len(a.chunkDDs); i++ {
		if a.chunkDDs[i].ID == id {
			a.chunkDDs[i] = a.chunkDDs[len(a.chunkDDs)-1]
			a.chunkDDs = a.chunkDDs[:len(a.chunkDDs)]
			return
		}
	}
}

// AddModelDD adds the model draw descriptor to the list to draw.
func (a *App) AddModelDD(d *ModelDrawDescriptor) {
	a.modelDDs = append(a.modelDDs, d)
}

// RemoveModelDD removes the model draw descriptor by ID.
func (a *App) RemoveModelDD(id uint32) {
	for i := 0; i < len(a.modelDDs); i++ {
		if a.modelDDs[i].ID == id {
			a.modelDDs[i] = a.modelDDs[len(a.modelDDs)-1]
			a.modelDDs = a.modelDDs[:len(a.modelDDs)]
			return
		}
	}
}

// AddLineDD adds the line draw descriptor to the list to draw.
func (a *App) AddLineDD(d *LineMeshDrawDescriptor) {
	a.lineDDs = append(a.lineDDs, d)
}

// RemoveLineDD removes the line draw descriptor by ID.
func (a *App) RemoveLineDD(id uint32) {
	for i := 0; i < len(a.lineDDs); i++ {
		if a.lineDDs[i].ID == id {
			a.lineDDs[i] = a.lineDDs[len(a.lineDDs)-1]
			a.lineDDs = a.lineDDs[:len(a.lineDDs)]
			return
		}
	}
}

// Draw draws everything with the given camera for 3D space..
func (a *App) Draw(c *Camera) {
	// Variable setup
	pMat := mgl32.Perspective(mgl32.DegToRad(60),
		float32(t.VirtualScreenWidth)/
			float32(t.VirtualScreenHeight),
		0.1, 1000.0)
	vMat := c.TransformMatrix()
	// Frame setup
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// Draw wire frames
	if a.WireFramesVisible {
		// Debug axis indicator
		a.pWireFrame.use()
		gl.UniformMatrix4fv(int32(a.pWireFrame.uni("uProjectionMatrix")), 1, false,
			&pMat[0])
		mt := c.TransformMatrix().Mul4(a.axis.Orientation.TransformMatrix())
		gl.UniformMatrix4fv(a.pWireFrame.uni("uModelViewMatrix"), 1, false, &mt[0])
		a.axis.draw(a.pWireFrame)
		// Debug line meshes
		for _, d := range a.lineDDs {
			if d.Mesh == nil {
				continue
			}
			mvm := mt.Mul4(d.Orientation.TransformMatrix())
			gl.UniformMatrix4fv(a.pCubeMesh.uni("uModelViewMatrix"), 1, false,
				&mvm[0])
			d.Mesh.draw(a.pWireFrame)
		}
		// Chunk bounds
		for _, d := range a.chunkDDs {
			mvm := mt.Mul4(t.O().Translate(d.CubeDD.Position).TransformMatrix())
			gl.UniformMatrix4fv(a.pCubeMesh.uni("uModelViewMatrix"), 1, false,
				&mvm[0])
			a.chunkBounds.draw(a.pWireFrame)
		}
		// Model bounds
		for _, m := range a.modelDDs {
			mvm := mt.Mul4(t.O().Translate(m.Orientation.P).TransformMatrix())
			gl.UniformMatrix4fv(a.pCubeMesh.uni("uModelViewMatrix"), 1, false,
				&mvm[0])
			m.Bounds.draw(a.pWireFrame)
		}
	}
	// Draw chunks
	a.pCubeMesh.use()
	gl.UniformMatrix4fv(int32(a.pCubeMesh.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	a.faces.bind(a.pCubeMesh)
	for _, d := range a.chunkDDs {
		if d.CubeDD.Mesh == nil {
			continue
		}
		mt := vMat.Mul4(mgl32.Translate3D(
			d.CubeDD.Position[0],
			d.CubeDD.Position[1],
			d.CubeDD.Position[2],
		))
		gl.UniformMatrix4fv(a.pCubeMesh.uni("uModelViewMatrix"), 1, false,
			&mt[0])
		d.CubeDD.Mesh.draw(a.pCubeMesh)
	}
	// Draw voxel cells
	a.pVoxelMesh.use()
	gl.Uniform1fv(a.pVoxelMesh.uni("uLightLevels"),
		int32(len(voxelLightLevels)), &voxelLightLevels[0])
	gl.UniformMatrix4fv(int32(a.pVoxelMesh.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	gl.UniformMatrix4fv(a.pVoxelMesh.uni("uViewMatrix"), 1, false,
		&vMat[0])
	for _, d := range a.chunkDDs {
		for _, v := range d.VoxelDDs {
			if v.Mesh == nil {
				continue
			}
			o := t.FacingToOrientation[v.Facing]
			o.P = v.Position.Add(mgl32.Vec3{8, 8, 8}.Mul(t.VoxelScale))
			mm := o.TransformMatrix()
			gl.UniformMatrix4fv(a.pVoxelMesh.uni("uModelMatrix"), 1, false,
				&mm[0])
			gl.Uniform1i(a.pVoxelMesh.uni("uFacing"), int32(v.Facing))
			v.Mesh.draw(a.pVoxelMesh)
		}
	}
	// Draw voxel models
	a.pModelMesh.use()
	gl.Uniform1fv(a.pModelMesh.uni("uLightLevels"), 6, &voxelLightLevels[0])
	gl.UniformMatrix4fv(int32(a.pModelMesh.uni("uProjectionMatrix")), 1, false,
		&pMat[0])
	gl.UniformMatrix4fv(a.pModelMesh.uni("uViewMatrix"), 1, false,
		&vMat[0])
	for _, d := range a.modelDDs {
		if d.Root != nil {
			d.Root.draw(a.pModelMesh, d.Orientation)
		}
	}
	// Draw UI elements, tiles layer
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	a.pUI.use()
	pMat = mgl32.Ortho2D(0, float32(t.VirtualScreenWidth), 0,
		float32(t.VirtualScreenHeight))
	gl.UniformMatrix4fv(a.pUI.uni("uProjectionMatrix"), 1, false, &pMat[0])
	a.tiles.bind(a.pUI)
	for _, m := range a.uiMeshes {
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			float32(m.Layer)/0xFFFF)
		m.draw(a.pUI)
	}
	// Draw common screen components
	if a.CrosshairVisible && a.crosshair != nil {
		m := a.crosshair
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			float32(m.Layer)/0xFFFF)
		m.draw(a.pUI)
	}
	if a.CursorVisible && a.cursor != nil {
		m := a.cursor
		gl.Uniform3f(a.pUI.uni("uPosition"), m.Position[0], -m.Position[1],
			float32(m.Layer)/0xFFFF)
		m.draw(a.pUI)
	}
	// Draw UI elements, text layer
	a.pText.use()
	if a.fm.imgDirty {
		a.fm.updateAtlasTexture()
	}
	gl.UniformMatrix4fv(a.pText.uni("uProjectionMatrix"), 1, false, &pMat[0])
	a.fm.bind(a.pText)
	for _, m := range a.uiMeshes {
		gl.Uniform3f(a.pText.uni("uPosition"), m.Position[0], -m.Position[1],
			float32(m.Layer)/0xFFFF)
		m.Text.draw()
	}
	// Draw debug text on top of everything
	if a.DebugTextVisible && a.debugText != nil {
		a.updateDebugText()
		m := a.debugText
		gl.Uniform3f(a.pText.uni("uPosition"), 0, 0, 1.0)
		m.draw()
	}
	a.debugLines = a.debugLines[:0]
	// Draw UI elements, cube mesh layer
	a.pCubeMeshIcon.use()
	gl.UniformMatrix4fv(a.pCubeMeshIcon.uni("uProjectionMatrix"), 1, false,
		&pMat[0])
	a.faces.bind(a.pCubeMeshIcon)
	for _, m := range a.uiMeshes {
		for _, d := range m.Cubes {
			if d.Mesh == nil {
				continue
			}
			mt := mgl32.Translate3D(
				d.Position[0],
				-d.Position[1],
				-d.Position[2],
			)
			gl.UniformMatrix4fv(a.pCubeMeshIcon.uni("uModelMatrix"), 1, false,
				&mt[0])
			d.Mesh.draw(a.pCubeMeshIcon)
		}
	}
}
