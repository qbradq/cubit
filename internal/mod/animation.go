package mod

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/c3d"
)

// AFQuat is the JSON structure used to hold X, Y, Z rotations to turn into
// quaternions.
type AFQuat [3]float32

func (q *AFQuat) UnmarshalJSON(d []byte) error {
	f := []float32{}
	if err := json.Unmarshal(d, &f); err != nil {
		return err
	}
	if len(f) != 3 {
		return errors.New("rotations are expressed as an array of three numbers representing the X, Y, and Z rotations")
	}
	q[0] = mgl32.DegToRad(f[0])
	q[1] = mgl32.DegToRad(f[1])
	q[2] = mgl32.DegToRad(f[2])
	return nil
}

// AnimationFrame describes a single frame of the animation.
type AnimationFrame struct {
	Time      float32           `json:"time"`      // Time T for interpolation to Q
	Rotations map[string]AFQuat `json:"rotations"` // Map of part ID's to quaternions that describe that part's final rotation
}

// Animation contains the information and functionality to play an animation
// loop on a model. Animations consist of rotations to joints only.
type Animation []*AnimationFrame

// animationsMap is the map of resource paths to animations.
var animationsMap = map[string]Animation{}

func registerAnimation(p string, a Animation) error {
	if _, duplicate := animationsMap[p]; duplicate {
		return fmt.Errorf("duplicate animation path %s", p)
	}
	animationsMap[p] = a
	return nil
}

// getAnimation constructs a new c3d.Animation object with the animation with
// the given resource path.
func getAnimation(p string, joints map[string]*mgl32.Quat) *c3d.Animation {
	a := animationsMap[p]
	if a == nil {
		log.Printf("warning: unknown animation %s\n", p)
		return nil
	}
	ret := &c3d.Animation{}
	for _, f := range a {
		af := &c3d.AnimationFrame{
			Time:    f.Time,
			TTL:     f.Time,
			Joints:  make(map[string]*mgl32.Quat),
			Targets: make(map[string]mgl32.Quat),
			Starts:  make(map[string]mgl32.Quat),
		}
		for k, f := range f.Rotations {
			af.Targets[k] = mgl32.AnglesToQuat(f[2], f[1], f[0], mgl32.ZYX)
		}
		for k, j := range joints {
			af.Joints[k] = j
		}
		ret.Frames = append(ret.Frames, af)
	}
	return ret
}
