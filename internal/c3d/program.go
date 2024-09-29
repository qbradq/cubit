package c3d

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/qbradq/cubit/data"
)

// program manages a single GPU program.
type program struct {
	id      uint32  // ID of the program
	vShader *shader // Vertex shader
	fShader *shader // Fragment shader
}

// loadProgram loads the named GLSL version 1.00 vertex+fragment shader program.
func loadProgram(name string) (*program, error) {
	d, err := data.FS.ReadFile(filepath.Join("glsl", name+".glsl"))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(d), "\n")
	which := 0
	var vSrc, fSrc string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch strings.ToLower(line) {
		case "[vertex]":
			which = 1
		case "[fragment]":
			which = 2
		default:
			switch which {
			case 1:
				vSrc += line + "\n"
			case 2:
				fSrc += line + "\n"
			default:
				return nil, fmt.Errorf(
					"in glsl program %s, no shader target given", name)
			}
		}
	}
	ret := &program{}
	ret.vShader, err = newShader(vSrc, shaderTypeVertex)
	if err != nil {
		return nil, fmt.Errorf("in glsl program %s, vertex shader: %s", name,
			err)
	}
	ret.fShader, err = newShader(fSrc, shaderTypeFragment)
	if err != nil {
		return nil, fmt.Errorf("in glsl program %s, fragment shader: %s", name,
			err)
	}
	ret.id = gl.CreateProgram()
	gl.AttachShader(ret.id, ret.vShader.id)
	gl.AttachShader(ret.id, ret.fShader.id)
	gl.LinkProgram(ret.id)
	if err := getGlError(ret.id, gl.LINK_STATUS, gl.GetProgramiv,
		gl.GetProgramInfoLog, "program:link"); err != nil {
		return nil, err
	}
	return ret, nil
}

// delete deletes the program.
func (p *program) delete() {
	p.vShader.delete()
	p.fShader.delete()
	gl.DeleteProgram(p.id)
}

// use makes this program the active one.
func (p *program) use() {
	gl.UseProgram(p.id)
}

// attr returns the id of a named shader attribute.
func (p *program) attr(name string) int32 {
	id := gl.GetAttribLocation(p.id, gl.Str(name+"\x00"))
	if id < 0 {
		log.Println("warning: unable to locate shader attribute \"" + name + "\"")
	}
	return id
}

// uni returns the id of a named shader uniform.
func (p *program) uni(name string) int32 {
	id := gl.GetUniformLocation(p.id, gl.Str(name+"\x00"))
	if id < 0 {
		log.Println("warning: unable to locate shader uniform \"" + name + "\"")
	}
	return id
}
