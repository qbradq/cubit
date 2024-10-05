package client

import (
	"math"
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
var dt float32             // Delta time for the current frame
var screenWidth int = 1280 // Width of the screen in pixels
var screenHeight int = 720 // Height of the screen in pixels
var win *glfw.Window       // GLFW window
var app *c3d.App           // Graphics application
var console *consoleWidget // Console widget
var input *Input           // Input instance
var world *t.World         // The currently loaded world

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
	app.CrosshairVisible = true
	app.SetCursor(mod.GetUITile("/cubit/004"), layerCursor)
	app.ChunkBoundsVisible = true
	// World setup
	world = t.NewWorld()
	TestGen(world)
	chunk := NewChunk(t.IVec3{0, 0, 0})
	chunk.Update()
	// app.AddChunkDD(chunk.cdd)
	// app.AddChunkDD(&c3d.ChunkDrawDescriptor{
	// 	ID: 1,
	// 	VoxelDDs: []*c3d.VoxelMeshDrawDescriptor{
	// 		{
	// 			ID:   1,
	// 			Mesh: t.GetVoxByPath("/cubit/vox/debug").Mesh,
	// 		},
	// 	},
	// })
	// model := mod.NewModel("/cubit/models/characters/brad")
	model := mod.NewModel("/cubit/models/characters/test")
	// model.DrawDescriptor.Orientation.Translate(mgl32.Vec3{0, 1, 0})
	// model.DrawDescriptor.Orientation.Yaw(math.Pi * 1.5)
	app.AddModelDD(model.DrawDescriptor)
	cam := c3d.NewCamera(mgl32.Vec3{1, 1, 5})
	// cam := c3d.NewCamera(mgl32.Vec3{7, 2, 7})
	// cam.Yaw = 90.001
	// Main loop
	lastFrame := glfw.GetTime()
	for !win.ShouldClose() {
		// Update state
		input.startFrame()
		glfw.PollEvents()
		input.PollEvents()
		currentFrame := glfw.GetTime()
		dt = float32(currentFrame - lastFrame)
		lastFrame = currentFrame
		console.update()
		model.DrawDescriptor.Orientation.Yaw(math.Pi * dt)
		// Handle input
		if input.WasPressed("debug") {
			app.DebugTextVisible = !app.DebugTextVisible
		}
		if console.isFocused() {
			console.input()
		} else {
			cameraInput(cam)
		}
		// TODO REMOVE
		app.AddDebugLine([3]uint8{255, 255, 0}, "Position: X=%d Y=%d Z=%d",
			int(cam.Position[0]),
			int(cam.Position[1]),
			int(cam.Position[2]),
		)
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

func cameraInput(cam *c3d.Camera) {
	speed := walkSpeed * dt
	if input.IsPressed("forward") {
		dir := cam.Front
		dir[1] = 0
		dir = dir.Normalize().Mul(speed)
		cam.Position = cam.Position.Add(dir)
	}
	if input.IsPressed("backward") {
		dir := cam.Front
		dir[1] = 0
		dir = dir.Normalize().Mul(speed)
		cam.Position = cam.Position.Sub(dir)
	}
	if input.IsPressed("left") {
		cam.Position = cam.Position.Sub(cam.Front.Cross(cam.Up).Normalize().Mul(speed))
	}
	if input.IsPressed("right") {
		cam.Position = cam.Position.Add(cam.Front.Cross(cam.Up).Normalize().Mul(speed))
	}
	if input.IsPressed("up") {
		cam.Position = cam.Position.Add(cam.Up.Mul(speed))
	}
	if input.IsPressed("down") {
		cam.Position = cam.Position.Sub(cam.Up.Mul(speed))
	}
	if input.IsPressed("turn-left") {
		cam.Yaw -= dt * 360.0 / 2.0
	}
	if input.IsPressed("turn-right") {
		cam.Yaw += dt * 360.0 / 2.0
	}
	if input.WasPressed("console") {
		console.stepVisibility()
	}
	cam.Yaw += input.CursorDelta[0] * mouseSensitivity
	cam.Pitch += input.CursorDelta[1] * mouseSensitivity
	if cam.Pitch > 89 {
		cam.Pitch = 89
	}
	if cam.Pitch < -89 {
		cam.Pitch = -89
	}
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
