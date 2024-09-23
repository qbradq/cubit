package client

import (
	"fmt"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/mitchellh/go-wordwrap"
	"github.com/qbradq/cubit/internal/c3d"
)

const conLines int = 1024
const conWidth int = 78
const conHeight int = 40

// consoleWidget implements a full screen widget that displays the log history
// and offers an interactive command line.
type consoleWidget struct {
	baseWidget
	lines     []string // Log lines
	lp        int      // Line pointer, points to the next index into lines to be used.
	vis       int      // Visibility code
	textDirty bool     // If true, all text of the widget needs redrawn
	prompt    string   // Prompt string
	po        int      // Offset into the prompt that appears on the left side
	cp        int      // Caret position
}

// newConsoleWidget creates a new console widget and returns it.
func newConsoleWidget(app *c3d.App) *consoleWidget {
	ret := &consoleWidget{
		baseWidget: *newBaseWidget(app),
		lines:      make([]string, conLines),
	}
	ret.cp = len(ret.prompt)
	// Setup UI
	ret.NinePatch(0, 0,
		c3d.VirtualScreenWidth, c3d.VirtualScreenHeight-c3d.CellDimsVS*3,
		npWindow)
	ret.NinePatch(0, c3d.VirtualScreenHeight-c3d.CellDimsVS*3,
		c3d.VirtualScreenWidth, c3d.CellDimsVS*3,
		npWindow)
	ret.Position = mgl32.Vec2{0, -float32(c3d.VirtualScreenHeight)}
	return ret
}

// printf add log lines to the console.
func (w *consoleWidget) printf(fm string, args ...any) {
	str := fmt.Sprintf(fm, args...)
	lines := strings.Split(wordwrap.WrapString(str, uint(conWidth)), "\n")
	for _, line := range lines {
		w.pushLine(line)
	}
}

// pushLine pushes a line onto the log.
func (w *consoleWidget) pushLine(line string) {
	w.textDirty = true
	w.lines[w.lp] = line
	w.lp++
	if w.lp >= conLines {
		w.lp = 0
	}
}

// stepVisibility steps through visibility modes.
func (w *consoleWidget) stepVisibility() {
	w.vis++
	if w.vis > 2 {
		w.vis = 0
	}
	switch w.vis {
	case 0:
		w.Position[1] = -float32(c3d.VirtualScreenHeight)
	case 1:
		w.Position[1] = -float32(c3d.VirtualScreenHeight / 2)
	case 2:
		w.Position[1] = 0
	}
}

// drawText draws all text into the widget.
func (w *consoleWidget) drawText() {
	// Cleanup
	w.Text.Reset()
	w.textDirty = false
	// Draw log lines
	li := w.lp
	li--
	if li < 0 {
		li = conLines - 1
	}
	p := [2]int{c3d.CellDimsVS, c3d.CellDimsVS * conHeight}
	for {
		if p[1] < c3d.CellDimsVS {
			break
		}
		line := w.lines[li]
		w.Text.Print(p[0], p[1], line)
		li--
		if li < 0 {
			li = conLines - 1
		}
		p[1] -= c3d.LineSpacingVS
	}
	// Adjust prompt offset
	l := w.po
	c := w.cp
	if c-l > (conWidth - 1) {
		l = c - (conWidth - 1)
	}
	if c-l < 0 {
		l = c
	}
	// Draw prompt string
	prompt := w.prompt[l:]
	if len(prompt) >= conWidth {
		prompt = prompt[:conWidth]
	}
	p[1] = c3d.CellDimsVS * (conHeight + 3)
	w.Text.Print(p[0], p[1], prompt)
	// Draw caret
	p[0] = c3d.CellDimsVS * ((c - l) + 1)
	w.Text.Print(p[0], p[1], "_")
}

// update implements the widget interface.
func (w *consoleWidget) update() {
	if w.textDirty {
		w.drawText()
	}
}

// isFocused returns true if the console widget is the input focus.
func (w *consoleWidget) isFocused() bool {
	return w.vis != 0
}

// input handles focused input.
func (w *consoleWidget) input() {
	for _, r := range input.CharsThisFrame {
		pl := w.prompt[:w.cp]
		pr := w.prompt[w.cp:]
		w.prompt = pl + string(r) + pr
		w.cp++
		w.textDirty = true
	}
	if input.WasPressed("ui-left") && w.cp > 0 {
		w.cp--
		w.textDirty = true
	}
	if input.WasPressed("ui-right") && w.cp < len(w.prompt) {
		w.cp++
		w.textDirty = true
	}
	if input.WasPressed("backspace") && w.cp > 0 {
		pl := w.prompt[:w.cp-1]
		pr := w.prompt[w.cp:]
		w.prompt = pl + pr
		w.cp--
		w.textDirty = true
	}
	if input.WasPressed("delete") && w.cp < len(w.prompt) {
		pl := w.prompt[:w.cp]
		pr := w.prompt[w.cp+1:]
		w.prompt = pl + pr
		w.textDirty = true
	}
	if input.WasPressed("confirm") {
		w.handleCommand(w.prompt)
		w.prompt = ""
		w.cp = 0
		w.textDirty = true
	}
	if input.WasPressed("console") {
		w.stepVisibility()
	}
	if input.WasPressed("cancel") {
		w.vis = 2
		w.stepVisibility()
	}
}

// handleCommand handles a command line.
func (w *consoleWidget) handleCommand(l string) {
	fields := strings.Fields(l)
	switch strings.ToLower(fields[0]) {
	case "exit":
		win.SetShouldClose(true)
	default:
		w.printf("error: unknown command %s", fields[0])
	}
}
