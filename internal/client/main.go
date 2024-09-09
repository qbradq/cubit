package client

import (
	"runtime"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/data"
	"github.com/qbradq/cubit/internal/c3d"
	"github.com/qbradq/cubit/internal/cubit"
)

/*
	mesh := prg.NewMesh([]float32{
		-0.5, 0.5, 0.5, 0, 0, 0, // Top left
		0.5, 0.5, 0.5, 1, 0, 0, // Top right
		-0.5, -0.5, 0.5, 0, 1, 0, // Bottom left
		-0.5, -0.5, 0.5, 0, 1, 0, // Bottom left
		0.5, 0.5, 0.5, 1, 0, 0, // Top right
		0.5, -0.5, 0.5, 1, 1, 0, // Bottom right
	})
*/

// Configuration variables
var mouseSensitivity float32 = 0.15
var walkSpeed float32 = 2.5

// Super globals
var debugTexture *c3d.Texture // Debugging texture
var dt float32                // Delta time for the current frame
var screenWidth int = 1280    // Width of the screen in pixels
var screenHeight int = 720    // Height of the screen in pixels

func init() {
	runtime.LockOSThread()
}

func Main() {
	// Window creation
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	win, err := glfw.CreateWindow(screenWidth, screenHeight, "Cubit", nil, nil)
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
	prg, err := glInit()
	if err != nil {
		panic(err)
	}
	defer prg.Delete()
	prg.Use()
	input := cubit.NewInput(win, mgl32.Vec2{float32(screenWidth), float32(screenHeight)})
	// Main loop
	cubeMesh := prg.NewCubeMesh()
	cubeMesh.AddFace(mgl32.Vec3{0, 0, 0}, c3d.North, 1)
	cubeMesh.AddFace(mgl32.Vec3{0, 0, 0}, c3d.South, 1)
	cubeMesh.AddFace(mgl32.Vec3{0, 0, 0}, c3d.East, 1)
	cubeMesh.AddFace(mgl32.Vec3{0, 0, 0}, c3d.West, 1)
	cubeMesh.AddFace(mgl32.Vec3{0, 0, 0}, c3d.Top, 1)
	cubeMesh.AddFace(mgl32.Vec3{0, 0, 0}, c3d.Bottom, 1)
	cubeMesh.Upload()
	cmo := c3d.NewOrientation(mgl32.Vec3{0, 0, 0}, 0, 0, 0)
	cam := c3d.NewCamera(mgl32.Vec3{0, 0, 3})
	lastFrame := glfw.GetTime()
	for !win.ShouldClose() {
		// Update state
		glfw.PollEvents()
		input.PollEvents()
		currentFrame := glfw.GetTime()
		dt = float32(currentFrame - lastFrame)
		lastFrame = currentFrame
		// TODO Remove Update camera
		speed := walkSpeed * dt
		if input.IsPressed("forward") {
			cam.Position = cam.Position.Add(cam.Front.Mul(speed))
		}
		if input.IsPressed("backward") {
			cam.Position = cam.Position.Sub(cam.Front.Mul(speed))
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
		if input.IsPressed("cancel") {
			win.SetShouldClose(true)
		}
		cam.Yaw += input.CursorDelta[0] * mouseSensitivity
		cam.Pitch += input.CursorDelta[1] * mouseSensitivity
		if cam.Pitch > 89 {
			cam.Pitch = 89
		}
		if cam.Pitch < -89 {
			cam.Pitch = -89
		}
		// Setup the frame buffer
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// Set projection and camera matrixes
		pMat := mgl32.Perspective(mgl32.DegToRad(60), float32(1280)/float32(720), 0.1, 100.0)
		gl.UniformMatrix4fv(int32(prg.GetUniformLocation("project")), 1, false, &pMat[0])
		cMat := cam.TransformMatrix()
		gl.UniformMatrix4fv(int32(prg.GetUniformLocation("camera")), 1, false, &cMat[0])
		// Draw
		cmo.Yaw(dt)
		cmo.Pitch(dt * 0.5)
		cmo.Roll(dt * 0.25)
		prg.DrawCubeMesh(cubeMesh, cmo)
		// Finish the frame
		win.SwapBuffers()
	}
}

func glInit() (*c3d.Program, error) {
	// Global initialization
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepthf(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	// Load debug texture
	if d, err := data.FS.ReadFile("textures/debug.png"); err != nil {
		return nil, err
	} else {
		debugTexture = c3d.NewTexture(d)
	}
	// Create shader program
	var vs *c3d.Shader
	if d, err := data.FS.ReadFile("shaders/vertex.glsl"); err != nil {
		return nil, err
	} else {
		if vs, err = c3d.NewShader(string(d), c3d.ShaderTypeVertex); err != nil {
			return nil, err
		}
	}
	var fs *c3d.Shader
	if d, err := data.FS.ReadFile("shaders/fragment.glsl"); err != nil {
		return nil, err
	} else {
		if fs, err = c3d.NewShader(string(d), c3d.ShaderTypeFragment); err != nil {
			return nil, err
		}
	}
	var p *c3d.Program
	var err error
	if p, err = c3d.NewProgram(vs, fs, debugTexture, cubit.Faces); err != nil {
		return nil, err
	}
	p.Use()
	return p, nil
}
