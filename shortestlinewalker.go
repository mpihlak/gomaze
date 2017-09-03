// Package maze ... simple navigator, always try to move towards the direction
// that has shortest straight line distance to finish.
//
// This makes a bunch of decisions that humans wouldn't make, but gets through
// the maze just fine. Possible enhancements:
// * Add a concept of visibility for more human-like traversal.
// * Don't look into visible cul de sacs.
// * Add awareness of the level layout
package maze

import "math"

// Keeps track of which directions have already been tried
type visitState struct {
	tried [4]bool
}

// ShortestLineWalker tries to navigate the maze by always aiming at the shortest
// line distance from current location.
type ShortestLineWalker struct {
	visitMap map[Position]visitState // map of all the directions taken
	actor    *Actor
	level    *Level
}

// Initialize sets up the walker state. It is not safe to start walking without initializing first.
func (walker *ShortestLineWalker) Initialize(level *Level, actor *Actor) {
	walker.visitMap = make(map[Position]visitState)
	walker.actor = actor
	walker.level = level
	actor.Path = make([]Position, 0)
}

// NextPosition picks a direction that has the shortest straight line distance from our
// current position. Avoid directions that we have already tried from here and if no options
// remain take a step back.
//
// For each cell on the map we keep track of the directions tries from there. This way
// we always know what has been tried before.
//
func (walker *ShortestLineWalker) NextPosition() {
	shortestLine := math.MaxInt32
	bestDirIndex := -1
	start := walker.actor.CurrPos
	finish := walker.actor.EndPos
	visitedDirections := walker.visitMap[start]

	// Try all valid directions that we have not already visited.
	// Pick the one that has shortest distance to finish.
	for index, dir := range ValidDirections {
		if visitedDirections.tried[index] {
			continue
		}

		// Calculate the euclidean distance from new position to finish
		xx := (start.col + dir.xd - finish.col)
		yy := (start.row + dir.yd - finish.row)
		dist := int(math.Sqrt(float64(xx*xx + yy*yy)))
		newPos := Position{row: start.row + dir.yd, col: start.col + dir.xd}

		// We might be stepping on a cell that we've already visited (eg. back
		// the way we came). Makes sense to consider unexplored cells first, so
		// add some "weight" to the distance.
		if walker.visitMap[newPos] != (visitState{}) {
			dist = dist * 1000
		}

		if dist < shortestLine && walker.level.CanMove(newPos) {
			shortestLine = dist
			bestDirIndex = index
		}
	}

	if bestDirIndex < 0 {
		// We've failed :(
    panic("I'm stuck here ...")
	} else {
		walker.actor.CurrPos = AddDirection(start, ValidDirections[bestDirIndex])
		walker.actor.Path = append(walker.actor.Path, walker.actor.CurrPos)

		// Record that we've tried this direction from here
		visitedDirections.tried[bestDirIndex] = true
		walker.visitMap[start] = visitedDirections
	}
}
