// Package maze ... walk the maze using a shortest path.
package maze

// PathNode represents a node on the graph. It's created from an empty space on the map.
type PathNode struct {
	parent   *PathNode
	distance int
	pos      Position
}

// CalculateShortestPath generates the shortest path from the current location of the actor to it's
// endPos using BFS. The level is not mutated in the process, the calculated path is stored in
// actor.
//
// The graph is represented here as a matrix[rows][columns] of nodes where an empty position in the
// matrix is a node (eg. something that can be walked on) and it's connected to it's neighboring
// nodes with edges of weight 1.
func CalculateShortestPath(level Level, actor *Actor, endPos Position) {
	// First build a map of all the empty tiles marked as unvisited. For convenience
	// We mark walls as visited, so that we don't consider them as nodes to visit.
	visitedTiles := make([][]bool, level.height)
	for row, rowTiles := range level.tiles {
		visitedTiles[row] = make([]bool, level.width)
		for col := range rowTiles {
			if level.tiles[row][col].tileType != EmptyTile {
				visitedTiles[row][col] = true
			} else {
				visitedTiles[row][col] = false
			}
		}
	}

	// PathNode at finish position
	var finishNode *PathNode

	// Queue of nodes that we're going to look at
	var nodes = []PathNode{
		PathNode{pos: Position{row: actor.CurrPos.row, col: actor.CurrPos.col}},
	}

	// Try stepping onto this node, if OK add it to the end of the queue
	tryStep := func(sourceNode *PathNode, row, col int) {
		newPos := Position{row: row, col: col}
		if !level.CanMove(newPos) || visitedTiles[row][col] {
			return
		}

		destNode := PathNode{
			distance: sourceNode.distance + 1,
			pos:      newPos,
		}

		destNode.distance = sourceNode.distance + 1
		destNode.parent = sourceNode
		nodes = append(nodes, destNode)
		visitedTiles[row][col] = true
	}

	for len(nodes) > 0 {
		n := nodes[0]
		nodes = nodes[1:]

		// Quit if we're already at finish position
		if n.pos == endPos {
			finishNode = &n
			break
		}

		// Reduce the distances for the nodes neighbors
		tryStep(&n, n.pos.row, n.pos.col+1)
		tryStep(&n, n.pos.row+1, n.pos.col)
		tryStep(&n, n.pos.row, n.pos.col-1)
		tryStep(&n, n.pos.row-1, n.pos.col)
	}

	// Map the path by tracing back from finish to start.
	actor.EndPos = endPos
	actor.Path = make([]Position, 0)
	for p := finishNode; p != nil; p = p.parent {
		actor.Path = append(actor.Path, p.pos)
	}
}

// ShortestPathWalker will navigate the maze using breadth first search
type ShortestPathWalker struct {
	actor     *Actor
	pathIndex int
}

// Initialize sets up the walker state. It is not safe to start walking without initializing first.
func (walker *ShortestPathWalker) Initialize(level *Level, actor *Actor) {
	walker.actor = actor
	CalculateShortestPath(*level, actor, actor.EndPos)
	walker.pathIndex = len(actor.Path) - 1
}

// HasFinished returns true if the walker has reached it's destination
func (walker *ShortestPathWalker) HasFinished() bool {
	return walker.actor.CurrPos == walker.actor.EndPos
}

// GetActor returns the internal actor
func (walker *ShortestPathWalker) GetActor() *Actor {
	return walker.actor
}

// NextPosition advances the actor to the next step on the path. Shortest path, in this case.
// Note that the actor.Path is calculated only once and if something changes on the map, we
// are in trouble.
//
// Also actor.Path has been calculated in the reverse order, so we must go through it backwards.
//
func (walker *ShortestPathWalker) NextPosition() {
	if walker.pathIndex >= 0 {
		walker.actor.CurrPos = walker.actor.Path[walker.pathIndex]
		walker.pathIndex--
	}
}
