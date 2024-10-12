package client

import (
	"runtime"
	"time"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/mod"
	"github.com/qbradq/cubit/internal/t"
)

// Configuration variables
var mouseSensitivity float32 = 0.15
var walkSpeed float32 = 5

// Super globals
var dt float32                       // Delta time for the current frame
var screenWidth int = 1280           // Width of the screen in pixels
var screenHeight int = 720           // Height of the screen in pixels
var win *glfw.Window                 // GLFW window
var app *c3d.App                     // Graphics application
var console *consoleWidget           // Console widget
var input *Input                     // Input instance
var world *t.World                   // The currently loaded world
var debugVector mgl32.Vec3           // Debug vector
var debugLines *c3d.LineMesh         // Debug lines mesh
var cam *c3d.Camera                  // Player camera
var cubeSelector *c3d.LineMesh       // Cube selection mesh
var csDD *c3d.LineMeshDrawDescriptor // Cube selector draw descriptor

func init() {
	c := [4]uint8{0, 255, 0, 255}
	cubeSelector = c3d.NewLineMesh()
	cs := cubeSelector
	cs.Hidden = true
	d := float32(1.0 / 16.0)
	w := -d
	e := 1 + d
	n := -d
	s := 1 + d
	b := -d
	t := 1 + d
	cs.Line(mgl32.Vec3{w, b, n}, mgl32.Vec3{e, b, n}, c)
	cs.Line(mgl32.Vec3{e, b, n}, mgl32.Vec3{e, b, s}, c)
	cs.Line(mgl32.Vec3{e, b, s}, mgl32.Vec3{w, b, s}, c)
	cs.Line(mgl32.Vec3{w, b, s}, mgl32.Vec3{w, b, n}, c)
	cs.Line(mgl32.Vec3{w, t, n}, mgl32.Vec3{e, t, n}, c)
	cs.Line(mgl32.Vec3{e, t, n}, mgl32.Vec3{e, t, s}, c)
	cs.Line(mgl32.Vec3{e, t, s}, mgl32.Vec3{w, t, s}, c)
	cs.Line(mgl32.Vec3{w, t, s}, mgl32.Vec3{w, t, n}, c)
	cs.Line(mgl32.Vec3{w, b, n}, mgl32.Vec3{w, t, n}, c)
	cs.Line(mgl32.Vec3{e, b, n}, mgl32.Vec3{e, t, n}, c)
	cs.Line(mgl32.Vec3{w, b, s}, mgl32.Vec3{w, t, s}, c)
	cs.Line(mgl32.Vec3{e, b, s}, mgl32.Vec3{e, t, s}, c)
	csDD = &c3d.LineMeshDrawDescriptor{
		ID:   1,
		Mesh: cs,
	}
}

// UI globals
var npWindow c3d.NinePatch

func init() {
	runtime.LockOSThread()
}

func Main() {
	var err error
	// Window creation
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	win, err = glfw.CreateWindow(screenWidth, screenHeight, "Cubit", nil, nil)
	if err != nil {
		panic(err)
	}
	win.MakeContextCurrent()
	win.SetPos(0, 0)
	// Load mods
	if err := mod.ReloadModInfo(); err != nil {
		panic(err)
	}
	if err := mod.LoadMods("cubit", "town"); err != nil {
		panic(err)
	}
	// OpenGL initialization
	app, err = glInit()
	if err != nil {
		panic(err)
	}
	defer app.Delete()
	// Globals init
	npWindow = c3d.NinePatch{
		mod.GetUITile("/cubit/000"),
		mod.GetUITile("/cubit/001"),
		mod.GetUITile("/cubit/002"),
		mod.GetUITile("/cubit/010"),
		mod.GetUITile("/cubit/011"),
		mod.GetUITile("/cubit/012"),
		mod.GetUITile("/cubit/020"),
		mod.GetUITile("/cubit/021"),
		mod.GetUITile("/cubit/022"),
	}
	input = NewInput(win, mgl32.Vec2{float32(screenWidth), float32(screenHeight)})
	console = newConsoleWidget(app)
	console.printf([3]uint8{0, 255, 255},
		"%s: Welcome to Cubit!", time.Now().Format(time.DateTime))
	console.add(app)
	app.SetCrosshair(mod.GetUITile("/cubit/003"), layerCrosshair)
	app.SetCursor(mod.GetUITile("/cubit/004"), layerCursor)
	app.CrosshairVisible = true
	app.CursorVisible = true
	app.WireFramesVisible = true
	app.DebugTextVisible = true
	// World setup
	world = t.NewWorld()
	TestGen(world)
	chunk := NewChunk(t.IVec3{0, 0, 0})
	chunk.Update()
	app.AddChunkDD(chunk.cdd)
	model := mod.NewModel("/cubit/models/characters/brad")
	model.DrawDescriptor.Orientation.P = mgl32.Vec3{6.5, 1.75, 10.5}
	model.DrawDescriptor.Orientation = model.DrawDescriptor.Orientation.Yaw(180)
	model.StartAnimation("/cubit/animations/characters/walk", "legs")
	app.AddModelDD(model.DrawDescriptor)
	// cam = c3d.NewCamera(mgl32.Vec3{2, 2, 5})
	// cam = c3d.NewCamera(mgl32.Vec3{1, 1, 5})
	cam = c3d.NewCamera(mgl32.Vec3{6.5, 2, 7})
	cam.Yaw = 90.001
	debugLines = c3d.NewLineMesh()
	app.AddLineDD(&c3d.LineMeshDrawDescriptor{
		ID:   1,
		Mesh: debugLines,
	})
	// Main loop
	app.AddLineDD(csDD)
	lastFrame := glfw.GetTime()
	for !win.ShouldClose() {
		// Update state
		input.startFrame()
		glfw.PollEvents()
		input.PollEvents()
		currentFrame := glfw.GetTime()
		dt = float32(currentFrame - lastFrame)
		lastFrame = currentFrame
		chunk.Update()
		model.Update(dt)
		console.update()
		// Handle input
		debugInput()
		if console.isFocused() {
			console.input()
		} else {
			cameraInput(cam)
			editInput()
		}
		// TODO REMOVE
		app.AddDebugLine([3]uint8{255, 255, 0}, "Position: X=%d Y=%d Z=%d",
			int(cam.Position[0]),
			int(cam.Position[1]),
			int(cam.Position[2]),
		)
		if wi != nil {
			app.AddDebugLine([3]uint8{0, 255, 0}, "WI: Pos=%v Face=%d",
				wi.Position, wi.Face)
		} else {
			app.AddDebugLine([3]uint8{0, 255, 0}, "WI: nil")
		}
		// Draw
		app.Draw(cam)
		// Finish the frame
		win.SwapBuffers()
	}
}

func glInit() (*c3d.App, error) {
	if err := gl.Init(); err != nil {
		return nil, err
	}
	gl.ClearColor(0, 0.5, 1, 0.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.ClearDepthf(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	return c3d.NewApp(mod.Faces, mod.UITiles)
}

// TODO DEBUG REMOVE
func TestGen(w *t.World) {
	rStone := mod.GetCubeDef("/cubit/cubes/stone")
	rGrass := mod.GetCubeDef("/cubit/cubes/grass")
	vWindow := mod.GetVoxByPath("/cubit/vox/window0")
	rect := func(min, max t.IVec3, r t.Cell) {
		for iy := min[1]; iy <= max[1]; iy++ {
			for iz := min[2]; iz <= max[2]; iz++ {
				for ix := min[0]; ix <= max[0]; ix++ {
					w.SetCell(t.IVec3{ix, iy, iz}, r)
				}
			}
		}
	}
	// Ground
	rect(t.IVec3{0, 0, 0}, t.IVec3{15, 0, 15}, t.CellForCube(rGrass, t.North))
	// Walls
	rect(t.IVec3{4, 1, 4}, t.IVec3{4, 3, 10}, t.CellForCube(rStone, t.North))
	rect(t.IVec3{10, 1, 4}, t.IVec3{10, 3, 10}, t.CellForCube(rStone, t.North))
	rect(t.IVec3{4, 1, 4}, t.IVec3{10, 3, 4}, t.CellForCube(rStone, t.North))
	rect(t.IVec3{4, 1, 10}, t.IVec3{10, 3, 10}, t.CellForCube(rStone, t.North))
	// Ceiling
	rect(t.IVec3{4, 4, 4}, t.IVec3{10, 4, 10}, t.CellForCube(rStone, t.North))
	// Window and doorway
	w.SetCell(t.IVec3{6, 1, 10}, t.CellInvalid)
	w.SetCell(t.IVec3{6, 2, 10}, t.CellInvalid)
	w.SetCell(t.IVec3{8, 2, 10}, t.CellForVox(vWindow.Ref, t.North))
}

var wi *t.WorldIntersection

// editInput is the input handler for editing mode.
func editInput() {
	cubeSelector.Hidden = true
	if input.InUIMode {
		return
	}
	ray := t.NewRay(cam.Position, cam.Front, 8.0)
	wi = ray.IntersectWorld(world)
	if wi == nil {
		return
	}
	cubeSelector.Hidden = false
	csDD.Orientation = t.O().Translate(mgl32.Vec3{
		float32(wi.Position[0]),
		float32(wi.Position[1]),
		float32(wi.Position[2]),
	})
	if input.ButtonPushed(2) {
		world.SetCell(wi.Position, t.CellInvalid)
	}
	if input.ButtonPushed(0) {
		p := t.PositionOffsets[wi.Face].Add(wi.Position)
		world.SetCell(p,
			t.CellForCube(
				mod.GetCubeDef("/cubit/cubes/grass"),
				t.North,
			),
		)
	}
}
