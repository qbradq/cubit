package client

import (
	"runtime"
	"time"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/cubit"
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
	// Load mods
	if err := cubit.ReloadModInfo(); err != nil {
		panic(err)
	}
	if err := cubit.LoadMods("cubit", "town"); err != nil {
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
		cubit.GetUITile("/cubit/000"),
		cubit.GetUITile("/cubit/001"),
		cubit.GetUITile("/cubit/002"),
		cubit.GetUITile("/cubit/010"),
		cubit.GetUITile("/cubit/011"),
		cubit.GetUITile("/cubit/012"),
		cubit.GetUITile("/cubit/020"),
		cubit.GetUITile("/cubit/021"),
		cubit.GetUITile("/cubit/022"),
	}
	input = NewInput(win, mgl32.Vec2{float32(screenWidth), float32(screenHeight)})
	console = newConsoleWidget(app)
	console.printf("%s: Welcome to Cubit!", time.Now().Format(time.DateTime))
	console.add(app)
	app.SetCrosshair(cubit.GetUITile("/cubit/003"), layerCrosshair)
	app.CrosshairVisible = true
	app.SetCursor(cubit.GetUITile("/cubit/004"), layerCursor)
	app.CursorVisible = true
	app.ChunkBoundsVisible = true
	// World setup
	world := cubit.NewWorld()
	// chunk := world.GetChunk(cubit.Pos(0, 0, 0))
	//chunk.Add(app)
	chunk := world.GetChunk(cubit.Pos(1, 0, 0))
	chunk.Add(app)
	cam := c3d.NewCamera(mgl32.Vec3{9, 13, 8})
	// cam := c3d.NewCamera(mgl32.Vec3{1, 1, 5})
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
		app.AddDebugLine("Position: X=%d Y=%d Z=%d",
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
	return c3d.NewApp(cubit.Faces, cubit.UITiles)
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
