package client

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/cubit/internal/t"
)

const (
	buttonLeft   int = 0
	buttonMiddle int = 1
	buttonRight  int = 2
)

type keySpec struct {
	key glfw.Key
	mod glfw.ModifierKey
}

var KeyConfig = map[string][]keySpec{
	"cancel":    {{glfw.KeyEscape, 0}},
	"confirm":   {{glfw.KeyEnter, 0}},
	"forward":   {{glfw.KeyW, 0}},
	"backward":  {{glfw.KeyS, 0}},
	"left":      {{glfw.KeyA, 0}},
	"right":     {{glfw.KeyD, 0}},
	"up":        {{glfw.KeyV, 0}},
	"down":      {{glfw.KeyC, 0}},
	"console":   {{glfw.KeyGraveAccent, glfw.ModControl}},
	"backspace": {{glfw.KeyBackspace, 0}},
	"delete":    {{glfw.KeyDelete, 0}},
	"ui-left":   {{glfw.KeyLeft, 0}},
	"ui-right":  {{glfw.KeyRight, 0}},
	"ui-up":     {{glfw.KeyUp, 0}},
	"ui-down":   {{glfw.KeyDown, 0}},
	"debug":     {{glfw.KeyF12, 0}},
}

// Input manages the input and input configuration.
type Input struct {
	CursorPosition mgl32.Vec2         // Current position of the mouse on the screen
	CursorDelta    mgl32.Vec2         // How far the mouse traveled this frame
	CharsThisFrame []rune             // List of runes generated this frame
	Mods           glfw.ModifierKey   // Modifier key mask this frame
	keysPressed    [glfw.KeyLast]bool // Array of all key states
	keysPushed     [glfw.KeyLast]bool // Array of all keys that had their release events this frame
	buttonsPressed [3]bool            // Array of mouse button states
	buttonsPushed  [3]bool            // Array of mouse buttons that had their release events this frame
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
	win.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused {
			win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		} else {
			win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}
	})
	win.SetCursorPosCallback(ret.posCallback)
	win.SetMouseButtonCallback(ret.mouseButtonCallback)
	return ret
}

func (n *Input) mouseButtonCallback(w *glfw.Window, button glfw.MouseButton,
	action glfw.Action, mods glfw.ModifierKey) {
	n.Mods |= mods
	var i int
	switch button {
	case glfw.MouseButtonLeft:
		i = buttonLeft
	case glfw.MouseButtonMiddle:
		i = buttonMiddle
	case glfw.MouseButtonRight:
		i = buttonRight
	default:
		return
	}
	switch action {
	case glfw.Press:
		n.buttonsPressed[i] = true
		n.buttonsPushed[i] = false
	case glfw.Release:
		n.buttonsPressed[i] = false
		n.buttonsPushed[i] = true
	}
}

func (n *Input) charCallback(w *glfw.Window, char rune) {
	n.CharsThisFrame = append(n.CharsThisFrame, char)
}

func (n *Input) keyCallback(win *glfw.Window, key glfw.Key, scanCode int,
	action glfw.Action, mods glfw.ModifierKey) {
	n.Mods |= mods
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
	if !n.seenMousePos {
		n.lastCursorPos = n.CursorPosition
		n.seenMousePos = true
	}
	app.SetCursorPosition(
		n.CursorPosition.Mul(
			float32(t.VirtualScreenWidth) / float32(screenWidth)))
}

func (n *Input) PollEvents() {
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
		if n.keysPressed[key.key] && n.Mods&key.mod == key.mod {
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
		if n.keysPushed[key.key] && n.Mods&key.mod == key.mod {
			return true
		}
	}
	return false
}

// ButtonPressed returns true if the give mouse button is pressed this frame.
func (n *Input) ButtonPressed(b int) bool {
	if b < 0 || b > buttonRight {
		return false
	}
	return n.buttonsPressed[b]
}

// ButtonPushed returns true if the give mouse button released this frame.
func (n *Input) ButtonPushed(b int) bool {
	if b < 0 || b > buttonRight {
		return false
	}
	return n.buttonsPushed[b]
}

// startFrame resets the keysPush array.
func (n *Input) startFrame() {
	n.CharsThisFrame = n.CharsThisFrame[:0]
	n.Mods = 0
	for i := range n.keysPushed {
		n.keysPushed[i] = false
	}
	for i := range n.buttonsPushed {
		n.buttonsPushed[i] = false
	}
}
