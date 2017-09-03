package maze

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

var (
	termboxInitialized = false
)

const (
	// KBEventUnknown -- unknown keyboard event
	KBEventUnknown = iota
	// KBEventCancel -- user has canceled the rendering (ESC or ^C)
	KBEventCancel
	// KBEventPause -- user has paused the rendering (^S or break)
	KBEventPause
)

// KeyboardEventChannel is used for passing keyboard events from the renderer to it's client.
type KeyboardEventChannel chan int

// Renderer is an interface that has knows how to display a maze on a terminal
// The interface is intended to be compatible with a dummy terminal so that the
// maze can be rendered into a file if needed.
type Renderer interface {
	Reset()
	PutChar(c rune)
	GetKeyboardEvent() KeyboardEventChannel
	NextLine()
	Flush()
	Done()
	Size() (int, int)
}

// Render draws the level and the path through it
func Render(level Level, banner string, r Renderer) {
	// Map out actors and their paths for quick lookup
	actorMap := make(map[Position]rune)
	for _, actor := range level.Actors {
		for _, pos := range actor.Path {
			// Avoid overriding actors with breadcrumbs, hence the lookup
			if _, ok := actorMap[pos]; !ok {
				actorMap[pos] = '.'
			}
		}
		actorMap[actor.CurrPos] = actor.Character
	}

	r.Reset()
	for _, c := range banner {
		r.PutChar(rune(c))
	}
	r.NextLine()

	// Display the level, tiles, actors and paths
	for row, tileRow := range level.tiles {
		for col, tile := range tileRow {
			pos := Position{row: row, col: col}
			c := tile.Character
			if ac, ok := actorMap[pos]; ok {
				c = ac
			}
			r.PutChar(c)
		}
		r.NextLine()
	}
	r.NextLine()
	r.Flush()
}

// TermboxRenderer uses the termbox library for rendering the maze.
type TermboxRenderer struct {
	row, col int
	kbEvents KeyboardEventChannel
}

// NewTermboxRenderer initialized the termbox library and creates a new instance of a renderer.
func NewTermboxRenderer() *TermboxRenderer {
	t := TermboxRenderer{}
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	t.kbEvents = make(KeyboardEventChannel)
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					fallthrough
				case termbox.KeyCtrlC:
					t.kbEvents <- KBEventCancel
				default:
					t.kbEvents <- KBEventUnknown
				}
			}
		}
	}()

	return &t
}

// GetKeyboardEvent returns a channel that can be polled for keyboard events from this renderer.
func (t *TermboxRenderer) GetKeyboardEvent() KeyboardEventChannel {
	return t.kbEvents
}

// Done is called to shut down the renderer.
func (t *TermboxRenderer) Done() {
	termbox.Close()
}

// Size returns terminal width and height.
func (t *TermboxRenderer) Size() (int, int) {
	return termbox.Size()
}

// NextLine advances the current row and resets the column to the start of the row.
func (t *TermboxRenderer) NextLine() {
	t.row++
	t.col = 0
}

// PutChar puts the character into the current position indicated by row and column and
// advances the column.
func (t *TermboxRenderer) PutChar(c rune) {
	termbox.SetCell(t.col, t.row, c, termbox.ColorDefault, termbox.ColorDefault)
	t.col++
}

// Reset the terminal so that we start again from the top left corner.
func (t *TermboxRenderer) Reset() {
	t.col = 0
	t.row = 0
}

// Flush renders the current state of the buffer.
func (t *TermboxRenderer) Flush() {
	termbox.Flush()
}

// StreamRenderer renders the maze on an output stream (file, stdout, etc.)
type StreamRenderer struct {
}

// NewStreamRenderer initializes and returns a new StreamRenderer
func NewStreamRenderer() *StreamRenderer {
	return &StreamRenderer{}
}

// GetKeyboardEvent returns a channel that can be polled for keyboard events from this renderer.
func (t *StreamRenderer) GetKeyboardEvent() KeyboardEventChannel {
	return nil
}

// Done is called to shut down the renderer.
func (t *StreamRenderer) Done() {
}

// Size returns terminal width and height.
func (t *StreamRenderer) Size() (int, int) {
	panic("StreamRenderer has no Size()")
}

// NextLine advances the current row and resets the column to the start of the row.
func (t *StreamRenderer) NextLine() {
	fmt.Print("\n")
}

// PutChar puts the character into the current position indicated by row and column and
// advances the column.
func (t *StreamRenderer) PutChar(c rune) {
	fmt.Print(string(c))
}

// Reset the terminal so that we start again from the top left corner.
func (t *StreamRenderer) Reset() {}

// Flush renders the current state of the buffer.
func (t *StreamRenderer) Flush() {}
