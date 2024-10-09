package c3d

import (
	"github.com/go-gl/mathgl/mgl32"
)

// AnimationFrame describes a single frame of the animation.
type AnimationFrame struct {
	Time    float32                // Time T for interpolation to Q
	TTL     float32                // Time left to move
	Joints  map[string]*mgl32.Quat // Map of part ID's to quaternions that describe that part's final rotation
	Targets map[string]mgl32.Quat  // Target rotations for each joint
	Starts  map[string]mgl32.Quat  // Staring rotation of the part in the frame
}

// setRotations sets the joint rotations for the given time t, where t is in the
// open range (0.0-1.0).
func (f *AnimationFrame) setRotations(t float32) {
	for k := range f.Targets {
		q := mgl32.QuatNlerp(
			f.Starts[k],
			f.Targets[k],
			t,
		)
		q = q.Normalize()
		*f.Joints[k] = q
	}
}

// Animation contains the information and functionality to play an animation
// loop on a model. Animations consist of rotations to joints only.
type Animation struct {
	Frames []*AnimationFrame // List of all frames
	i      int               // Current frame index
}

// Play starts the animation playing.
func (a *Animation) Play() {
	a.i = 0
	if len(a.Frames) < 1 {
		return
	}
	a.startFrame()
}

// startFrame starts the current frame.
func (a *Animation) startFrame() {
	f := a.Frames[a.i]
	for k := range f.Joints {
		f.Starts[k] = *f.Joints[k]
	}
	f.TTL = f.Time
}

// Update advances toward the next frame, ensuring that when a frame
// transition occurs the final position is reached by the target joint. The t
// argument is the number of seconds elapsed since the last call to this
// function. Frames are transitioned as appropriate.
func (a *Animation) Update(dt float32) {
	if len(a.Frames) < 1 {
		return
	}
	for {
		f := a.Frames[a.i]
		ft := dt
		if ft > f.TTL {
			ft = f.TTL
		}
		dt -= ft
		f.TTL -= ft
		qt := 1.0 - (f.TTL / f.Time)
		f.setRotations(qt)
		if f.TTL <= 0 {
			a.i++
			if a.i >= len(a.Frames) {
				a.i = 0
			}
			a.startFrame()
		}
		if dt <= 0 {
			break
		}
	}
}
