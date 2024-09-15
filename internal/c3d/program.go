package c3d

import (
	"log"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	drawModeColor   int32 = 0
	drawModeAtlas   int32 = 1
	drawModeTexture int32 = 2
)

// Program manages a GPU program with vertex and pixel shaders.
type Program struct {
	id        uint32     // OpenGL ID of the program
	vShader   *Shader    // Vertex shader
	fShader   *Shader    // Fragment shader
	tex0      *Texture   // Texture bound to uniform tex0
	atlas     *FaceAtlas // Atlas used to render textured cube mesh
	aPos      int32      // pos attribute location
	aTex      int32      // tex attribute location
	uDrawMode int32      // drawMode uniform location
	uWorld    int32      // world uniform location
	uTexture  int32      // tex0 uniform location
	uAtlas    int32      // atlas uniform location
}

// NewProgram creates a new Program ready for use with the given resources.
func NewProgram(v, f *Shader, tex0 *Texture, atlas *FaceAtlas) (*Program, error) {
	ret := &Program{
		id:      gl.CreateProgram(),
		vShader: v,
		fShader: f,
		tex0:    tex0,
		atlas:   atlas,
	}
	gl.AttachShader(ret.id, ret.vShader.id)
	gl.AttachShader(ret.id, ret.fShader.id)
	gl.LinkProgram(ret.id)
	if err := getGlError(ret.id, gl.LINK_STATUS, gl.GetProgramiv,
		gl.GetProgramInfoLog, "program:link"); err != nil {
		return nil, err
	}
	ret.aPos = ret.GetAttributeLocation("pos")
	ret.aTex = ret.GetAttributeLocation("texPos")
	ret.uDrawMode = ret.GetUniformLocation("drawMode")
	ret.uWorld = ret.GetUniformLocation("world")
	ret.uTexture = ret.GetUniformLocation("tex")
	ret.uAtlas = ret.GetUniformLocation("atlas")
	ret.tex0.bind(ret.uTexture)
	ret.atlas.upload()
	ret.atlas.freeMemory()
	ret.atlas.bind(ret.uAtlas)
	return ret, nil
}

// Delete deletes the program.
func (p *Program) Delete() {
	p.tex0.unbind()
	p.atlas.unbind()
	p.vShader.Delete()
	p.fShader.Delete()
	gl.DeleteProgram(p.id)
}

// Use makes this program the active one.
func (p *Program) Use() {
	gl.UseProgram(p.id)
}

// GetAttributeLocation returns the id of a named shader attribute.
func (p *Program) GetAttributeLocation(name string) int32 {
	id := gl.GetAttribLocation(p.id, gl.Str(name+"\x00"))
	if id < 0 {
		log.Println("warning: unable to locate shader attribute \"" + name + "\"")
	}
	return id
}

// GetUniformLocation returns the id of a named shader uniform.
func (p *Program) GetUniformLocation(name string) int32 {
	id := gl.GetUniformLocation(p.id, gl.Str(name+"\x00"))
	if id < 0 {
		log.Println("warning: unable to locate shader uniform \"" + name + "\"")
	}
	return id
}

// DrawCubeMesh draws a cube mesh.
func (p *Program) DrawCubeMesh(m *CubeMesh, o *Orientation) {
	m.draw(p, o)
}

// NewAxisIndicator creates a new axis indicator at the given location.
func (p *Program) NewAxisIndicator(pos mgl32.Vec3) *AxisIndicator {
	return newAxisIndicator(pos, p)
}
