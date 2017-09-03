// Package maze provides level data structures, initialization and loading
package maze

import (
	"bufio"
)

// Position coordinates on the level grid
type Position struct {
	row, col int
}

// Direction is something we can move towards
type Direction struct {
	xd int // Column delta
	yd int // Row delta
}

// ValidDirections specifies all the valid directions we can move towards
var ValidDirections = [4]Direction{
	{0, 1},  // down
	{1, 0},  // right
	{-1, 0}, // left
	{0, -1}, // up
}

// Tile types
const (
	EmptyTile = iota
	WallTile
)

// Tile is a map element on the level
type Tile struct {
	tileType  int
	Character rune
}

// Actor is something that can move around on the level
type Actor struct {
	Character rune       // Display character
	CurrPos   Position   // Current location
	EndPos    Position   // Destination, if calculated
	Path      []Position // Path, if calculated.
  PathNav   Walker
}

// Level describes the map and everything on it
type Level struct {
	width  int
	height int
	tiles  [][]Tile   // Level map[row][col]
	Actors []*Actor   // Various moving actors on the level
	Exits  []Position // Exits on the level
}

// ReadLevel reads a level from ASCII art
func ReadLevel(scanner *bufio.Scanner) Level {
	level := Level{}
	for row := 0; scanner.Scan(); row++ {
		var tileRow []Tile
		for col, c := range scanner.Text() {
			tileType := EmptyTile
			pos := Position{row: row, col: col}

			switch c {
			case '@', '?', '!', '&':
				level.Actors = append(level.Actors, &Actor{Character: c, CurrPos: pos})
				c = ' '
			case '=':
				level.Exits = append(level.Exits, pos)
			case '#':
				tileType = WallTile
			}
			tileRow = append(tileRow, Tile{tileType: tileType, Character: c})
		}
		level.tiles = append(level.tiles, tileRow)
		if row == 0 {
			level.width = len(tileRow)
		}
		level.height++
	}
	return level
}

// WithinBounds checks if the position is on the level
func (level Level) WithinBounds(pos Position) bool {
	return pos.col >= 0 && pos.row >= 0 && pos.col < level.width && pos.row < level.height
}

// WithinFrame checks if the position is within the level frame
func (level Level) WithinFrame(pos Position) bool {
	return pos.row > 0 && pos.col > 0 && pos.row < level.height-1 && pos.col < level.width-1
}

// IsWalkable checks if we can walk on this position
func (level Level) IsWalkable(pos Position) bool {
	if t := level.tiles[pos.row][pos.col]; t.tileType != EmptyTile {
		return false
	}
	return true
}

// CanMove tells if the position on the level is vacant or not
func (level Level) CanMove(pos Position) bool {

	if !level.WithinBounds(pos) {
		return false
	}

	if !level.IsWalkable(pos) {
		return false
	}

	// How about other actors?
	for _, actor := range level.Actors {
		if actor.CurrPos == pos {
			return false
		}
	}

	return true
}

// HasFinished returns true if the actor has reached its destination
func (a Actor) HasFinished() bool {
	return a.CurrPos == a.EndPos
}

// AddDirection adds the direction to position and returns the new position
func AddDirection(pos Position, d Direction) Position {
	return Position{row: pos.row + d.yd, col: pos.col + d.xd}
}

// AddActor adds a new Actor to the level
func (level *Level) AddActor(a *Actor) {
	level.Actors = append(level.Actors, a)
}

// AddWalker adds a new walker to the level
func (level *Level) AddWalker(c rune, startPos, endPos Position, navigator Walker) {
  actor := &Actor{ Character: c, CurrPos: startPos, EndPos: endPos, PathNav: navigator }
  navigator.Initialize(level, actor)
  level.Actors = append(level.Actors, actor)
}
