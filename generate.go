// Package maze deals with maze generation.
package maze

import (
	"math/rand"
)

const (
	// WallBlock is used for drawing all maze walls
	WallBlock = 0x2588
)

// GenerateRandomMaze generates a maze that does not contain disconnected rooms.
// Eg. if we place an actor to an empty spot on the level, it should be able to
// navigate to every other empty spot.
// There will be 2 exits one in the top left and another in the bottom right.
func GenerateRandomMaze(width, height int) Level {
	level := MakeEmptyLevel(width, height)

	// Fill the inside of the level with tiles, we're gonna plough into it to make a maze.
	for row, tileRow := range level.tiles {
		for col := range tileRow {
			if row > 0 && col > 0 && row < level.height-1 && col < level.width-1 {
				level.tiles[row][col] = Tile{WallTile, WallBlock}
			}
		}
	}

	// Generate a path from top left to bottom right.
	startPos := Position{row: 0, col: 1}
	exitPos := Position{row: level.height - 1, col: level.width - 2}

	// turnStack stores all the turning points, so that we can pop a previous location
	// when we run into dead end.
	turnStack := make([]Position, 0)
	turnStack = append(turnStack, startPos)
	pos := startPos

	// Keep generating while we have not considered all options.
	// not done until we've popped the last element off the turnStack
	for len(turnStack) > 0 {
		steps := 0

		// Randomize the directions to be tried
		shuffle := ValidDirections
		for i := 0; i < len(shuffle); i++ {
			rndPos := rand.Intn(len(shuffle))
			shuffle[i], shuffle[rndPos] = shuffle[rndPos], shuffle[i]
		}

		// Try moving in random directions, until we can make at least 1 step
		for _, dir := range shuffle {
			maxSteps := rand.Intn(4) + 2
			for steps < maxSteps {
				newPos := AddDirection(pos, dir)

				// Stop walking if we've been here before
				if !level.WithinFrame(newPos) || level.IsWalkable(newPos) {
					break
				}
				if !hasEnoughWalls(level, newPos, pos) {
					break
				}
				level.tiles[newPos.row][newPos.col] = Tile{EmptyTile, ' '}
				pos = newPos
				steps++
			}

			if steps > 0 {
				turnStack = append(turnStack, pos)
				break
			}
		}

		// If we didn't manage to go anywhere, pop the last position and restart from there
		if steps == 0 {
			pos = turnStack[len(turnStack)-1]
			turnStack = turnStack[:len(turnStack)-1]
		}
	}

	level.CreateHorizontalExit(startPos)
	level.CreateHorizontalExit(exitPos)
	return level
}

// hasEnoughWalls validates that there is enough surrounding space around the position.
// In the surrounding area there should be at least 1 wall tile in every direction (except
// the direction we just came from)
func hasEnoughWalls(level Level, pos Position, origin Position) bool {
	for _, dir := range ValidDirections {
		newPos := AddDirection(pos, dir)
		if newPos == origin {
			// Skip that. It's empty, we just came from there.
			continue
		}
		if !level.WithinFrame(newPos) {
			// Okey, we're surrounded by the frame
			continue
		}
		if level.IsWalkable(newPos) {
			// Oh, we've broken through. No good.
			return false
		}
	}
	return true
}

// MakeEmptyLevel generates a level frame, borders, corners. etc.
func MakeEmptyLevel(width, height int) Level {
	level := Level{width: width, height: height}
	level.tiles = make([][]Tile, height)
	for row := range level.tiles {
		level.tiles[row] = make([]Tile, width)
		for col := range level.tiles[row] {
			var t Tile
			if row == 0 || row == height-1 {
				t = Tile{WallTile, WallBlock}
			} else if col == 0 || col == width-1 {
				t = Tile{WallTile, WallBlock}
			} else {
				t = Tile{EmptyTile, ' '}
			}
			level.tiles[row][col] = t
		}
	}
	return level
}

// CreateHorizontalExit creates an exit in the horizontal frame of the level.
func (level *Level) CreateHorizontalExit(pos Position) {
	level.tiles[pos.row][pos.col] = Tile{EmptyTile, ' '}
	level.Exits = append(level.Exits, pos)

	// Plough through the level to remove any obstacles that might be blocking the exit
	dir := Direction{xd: 0, yd: +1}

	switch {
	case pos.row == 0:
		dir = Direction{xd: 0, yd: 1}
	case pos.col == 0:
		dir = Direction{xd: 1, yd: 0}
	case pos.row == level.height-1:
		dir = Direction{xd: 0, yd: -1}
	case pos.col == level.width-1:
		dir = Direction{xd: -1, yd: 0}
	default:
		panic("Can't plough the exit!")
	}

	// Plough through the level until we hit an opening
	for pos = AddDirection(pos, dir); level.WithinFrame(pos) && !level.CanMove(pos); pos = AddDirection(pos, dir) {
		level.tiles[pos.row][pos.col] = Tile{EmptyTile, ' '}
	}
}
