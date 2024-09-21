package c3d

import gl "github.com/go-gl/gl/v3.1/gles2"

// shaderType indicates the type of shader.
type shaderType int

const (
	shaderTypeVertex   shaderType = 0 // Vertex shader
	shaderTypeFragment shaderType = 1 // Fragment shader
)

// shader manages a single shader.
type shader struct {
	id uint32 // OpenGL ID of the shader
}

// newShader creates a new Shader ready for use.
func newShader(src string, t shaderType) (*shader, error) {
	ts := "vertex"
	if t == shaderTypeFragment {
		ts = "fragment"
	}
	ret := &shader{}
	switch t {
	case shaderTypeVertex:
		ret.id = gl.CreateShader(gl.VERTEX_SHADER)
	case shaderTypeFragment:
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

// delete deletes the shader.
func (s *shader) delete() {
	gl.DeleteProgram(s.id)
}
