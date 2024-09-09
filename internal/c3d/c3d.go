package c3d

import (
	"fmt"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// Facing encodes one of the facing values, North, South, East, West, Up, Down.
type Facing uint8

const (
	North Facing = iota
	South
	East
	West
	Top
	Bottom
)

var XAxis mgl32.Vec3 = mgl32.Vec3{1, 0, 0}
var YAxis mgl32.Vec3 = mgl32.Vec3{0, 1, 0}
var ZAxis mgl32.Vec3 = mgl32.Vec3{0, 0, 1}

var FaceNormals = [6]mgl32.Vec3{
	{0, 0, -1},
	{0, 0, 1},
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, log)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
	}

	return nil
}

var cubeFacingOffsets = [6][6][3]float32{
	{ // North
		{+0.5, +0.5, -0.5}, // TL
		{-0.5, +0.5, -0.5}, // TR
		{+0.5, -0.5, -0.5}, // BL
		{+0.5, -0.5, -0.5}, // BL
		{-0.5, +0.5, -0.5}, // TR
		{-0.5, -0.5, -0.5}, // BR
	},
	{ // South
		{-0.5, +0.5, +0.5}, // TL
		{+0.5, +0.5, +0.5}, // TR
		{-0.5, -0.5, +0.5}, // BL
		{-0.5, -0.5, +0.5}, // BL
		{+0.5, +0.5, +0.5}, // TR
		{+0.5, -0.5, +0.5}, // BR
	},
	{ // East
		{+0.5, +0.5, +0.5}, // TL
		{+0.5, +0.5, -0.5}, // TR
		{+0.5, -0.5, +0.5}, // BL
		{+0.5, -0.5, +0.5}, // BL
		{+0.5, +0.5, -0.5}, // TR
		{+0.5, -0.5, -0.5}, // BR
	},
	{ // West
		{-0.5, +0.5, -0.5}, // TL
		{-0.5, +0.5, +0.5}, // TR
		{-0.5, -0.5, -0.5}, // BL
		{-0.5, -0.5, -0.5}, // BL
		{-0.5, +0.5, +0.5}, // TR
		{-0.5, -0.5, +0.5}, // BR
	},
	{ // Top
		{-0.5, +0.5, -0.5}, // TL
		{+0.5, +0.5, -0.5}, // TR
		{-0.5, +0.5, +0.5}, // BL
		{-0.5, +0.5, +0.5}, // BL
		{+0.5, +0.5, -0.5}, // TR
		{+0.5, +0.5, +0.5}, // BR
	},
	{ // Bottom
		{+0.5, -0.5, -0.5}, // TL
		{-0.5, -0.5, -0.5}, // TR
		{+0.5, -0.5, +0.5}, // BL
		{+0.5, -0.5, +0.5}, // BL
		{-0.5, -0.5, -0.5}, // TR
		{-0.5, -0.5, +0.5}, // BR
	},
}
