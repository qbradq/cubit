package client

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var KeyConfig = map[string][]glfw.Key{
	"cancel":    {glfw.KeyEscape},
	"confirm":   {glfw.KeyEnter},
	"forward":   {glfw.KeyW},
	"backward":  {glfw.KeyS},
	"left":      {glfw.KeyA},
	"right":     {glfw.KeyD},
	"up":        {glfw.KeyV},
	"down":      {glfw.KeyC},
	"console":   {glfw.KeyGraveAccent},
	"backspace": {glfw.KeyBackspace},
	"delete":    {glfw.KeyDelete},
	"ui-left":   {glfw.KeyLeft},
	"ui-right":  {glfw.KeyRight},
	"ui-up":     {glfw.KeyUp},
	"ui-down":   {glfw.KeyDown},
}

// Input manages the input and input configuration.
type Input struct {
	CursorPosition mgl32.Vec2         // Current position of the mouse on the screen
	CursorDelta    mgl32.Vec2         // How far the mouse traveled this frame
	CharsThisFrame []rune             // List of runes generated this frame
	keysPressed    [glfw.KeyLast]bool // Array of all key states
	keysPushed     [glfw.KeyLast]bool // Array of all keys that had their release events this frame
	lastCursorPos  mgl32.Vec2         // Last position of the mouse on the screen
	seenMousePos   bool               // If true we have already seen the mouse position at least once
}

// NewInput returns a new Input object ready for use.
func NewInput(win *glfw.Window, screenSize mgl32.Vec2) *Input {
	ret := &Input{
		lastCursorPos: screenSize.Mul(0.5),
	}
	ret.CursorPosition = ret.lastCursorPos
	win.SetKeyCallback(ret.keyCallback)
	win.SetCharCallback(ret.charCallback)
	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	win.SetCursorPosCallback(ret.posCallback)
	return ret
}

func (n *Input) charCallback(w *glfw.Window, char rune) {
	n.CharsThisFrame = append(n.CharsThisFrame, char)
}

func (n *Input) keyCallback(win *glfw.Window, key glfw.Key, scanCode int,
	action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		fallthrough
	case glfw.Repeat:
		n.keysPressed[key] = true
		n.keysPushed[key] = true
	case glfw.Release:
		n.keysPressed[key] = false
		n.keysPushed[key] = false
	}
}

func (n *Input) posCallback(w *glfw.Window, x, y float64) {
	n.CursorPosition[0] = float32(x)
	n.CursorPosition[1] = float32(y)
}

func (n *Input) PollEvents() {
	if !n.seenMousePos {
		n.lastCursorPos = n.CursorPosition
		n.seenMousePos = true
	}
	n.CursorDelta[0] = n.CursorPosition[0] - n.lastCursorPos[0]
	n.CursorDelta[1] = n.lastCursorPos[1] - n.CursorPosition[1]
	n.lastCursorPos = n.CursorPosition
}

// IsPressed returns true if the key associated with the given action is
// currently pressed down.
func (n *Input) IsPressed(action string) bool {
	keys, found := KeyConfig[action]
	if !found {
		return false
	}
	for _, key := range keys {
		if n.keysPressed[key] {
			return true
		}
	}
	return false
}

// WasPressed returns true if the key associated with the given action was newly
// released this frame.
func (n *Input) WasPressed(action string) bool {
	keys, found := KeyConfig[action]
	if !found {
		return false
	}
	for _, key := range keys {
		if n.keysPushed[key] {
			return true
		}
	}
	return false
}

// startFrame resets the keysPush array.
func (n *Input) startFrame() {
	n.CharsThisFrame = n.CharsThisFrame[:0]
	for i := range n.keysPushed {
		n.keysPushed[i] = false
	}
}
