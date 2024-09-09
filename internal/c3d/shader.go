package c3d

import gl "github.com/go-gl/gl/v3.1/gles2"

// ShaderType indicates the type of shader.
type ShaderType int

const (
	ShaderTypeVertex   ShaderType = 0 // Vertex shader
	ShaderTypeFragment ShaderType = 1 // Fragment shader
)

// Shader manages a single shader.
type Shader struct {
	id uint32 // OpenGL ID of the shader
}

// NewShader creates a new Shader ready for use.
func NewShader(src string, t ShaderType) (*Shader, error) {
	ts := "vertex"
	if t == ShaderTypeFragment {
		ts = "fragment"
	}
	ret := &Shader{}
	switch t {
	case ShaderTypeVertex:
		ret.id = gl.CreateShader(gl.VERTEX_SHADER)
	case ShaderTypeFragment:
		ret.id = gl.CreateShader(gl.FRAGMENT_SHADER)
	}
	glStrs, freeFunc := gl.Strs(src + "\x00")
	defer freeFunc()
	gl.ShaderSource(ret.id, 1, glStrs, nil)
	gl.CompileShader(ret.id)
	err := getGlError(ret.id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "shader:compile:"+ts)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Delete deletes the shader.
func (s *Shader) Delete() {
	gl.DeleteProgram(s.id)
}
