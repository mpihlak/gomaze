// Package maze ... walk the maze using a shortest path.
package maze

import (
	"container/heap"
	"math"
)

// PathNode represents a node on the graph. It's created from an empty space on the map.
type PathNode struct {
	parent   *PathNode
	distance int
	index    int
	pos      Position
}

// PathNodePQ is a priority queue of *PathNode items
type PathNodePQ []*PathNode

func (pq PathNodePQ) Len() int           { return int(len(pq)) }
func (pq PathNodePQ) Less(i, j int) bool { return pq[i].distance < pq[j].distance }

func (pq PathNodePQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// Push a new element into the PQ
func (pq *PathNodePQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PathNode)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes an element from the PQ
func (pq *PathNodePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// CalculateShortestPath generates the shortest path from the current location of the actor
// to it's endPos using Dijkstra's algorithm. The level is not mutated in the process, the
// calculated path is stored in actor.
//
// The graph is represented here as a matrix[rows][columns] of nodes where an
// empty position in the matrix is a node (eg. something that can be walked on)
// and it's connected to it's neighboring nodes with edges of weight 1.
func CalculateShortestPath(level Level, actor *Actor, endPos Position) {

	// nodeMap is a copy of the map that we can scribble on. It will be
	// updated with the path information and node distances.
	nodeMap := make([][]PathNode, level.height)

	// nodes holds all the nodes (pointers to nodeMap) that we have not yet
	// considered ordered by distance from the starting position.
	var nodes PathNodePQ

	// First build the map of nodes, so that every empty tile on the level will
	// be a node. Initialize the node's distance to infinity, if it's not the
	// starting node.  Build a priority queue of all the nodes by distance to the
	// starting node.
	for row, rowTiles := range level.tiles {
		nodeMap[row] = make([]PathNode, level.width)
		for col := range rowTiles {
			if level.tiles[row][col].tileType != EmptyTile {
				continue
			}
			n := PathNode{pos: Position{row: row, col: col}}
			if n.pos == actor.CurrPos {
				n.distance = 0
			} else {
				n.distance = math.MaxInt32
			}
			// Populate the map and push the node to the priority queue
			n.index = len(nodes)
			nodeMap[row][col] = n
			nodes = append(nodes, &nodeMap[row][col])
		}
	}

	// OK, actually arrange the PQ.
	heap.Init(&nodes)

	// Function to reduce the distance of sourceNode's neighbor
	//
	// This will be called for all the neighbours of the source node to update their
	// distances from the source. The distance from source is of course always 1, as we
	// are only looking at immediate neighbours.
	//
	// After fixing the distance, we register that the shortest path to destination came from the
	// source node and fix the order of the PQ.
	reduce := func(sourceNode *PathNode, row, col int) {
		if !level.CanMove(Position{row: row, col: col}) {
			return
		}
		destNode := &nodeMap[row][col]
		// Skip nodes that have already been removed from the queue and those that have
		// no entry in nodeMap (eg. non-empty cell on map)
		if destNode.index >= 0 && *destNode != (PathNode{}) {
			newDistance := sourceNode.distance + 1
			if destNode.distance > newDistance {
				destNode.distance = newDistance
				destNode.parent = sourceNode

				// Adjust the position in the PQ
				heap.Fix(&nodes, destNode.index)
			}
		}
	}

	// PathNode at finish position
	var finishNode *PathNode

	// The nodes PQ will have the nodes ordered by distance from the start position. We go through
	// the queue popping the closest node off the top (at first pass this is the starting node).
	// Then for each of this node's neighbours we update it's distances to 1 and update the PQ ordering.
	//
	// Proceed by popping the nearest node and reducing it's neighbours until the queue runs out or
	// we find a node that's at the finish position.
	//
	// The loop invariant here is that the nodes PQ always has the node closes
	for nodes.Len() > 0 {
		n := heap.Pop(&nodes).(*PathNode)

		// Quit if we're already at finish position
		if n.pos == endPos {
			finishNode = n
			break
		}

		// Reduce the distances for the nodes neighbors
		reduce(n, n.pos.row, n.pos.col+1)
		reduce(n, n.pos.row+1, n.pos.col)
		reduce(n, n.pos.row, n.pos.col-1)
		reduce(n, n.pos.row-1, n.pos.col)
	}

	// Map the path by tracing back from finish to start.
	actor.EndPos = endPos
	actor.Path = make([]Position, 0)
	for p := finishNode; p != nil; p = p.parent {
		actor.Path = append(actor.Path, p.pos)
	}
}

// ShortestPathWalker will navigate the maze using the shortest path
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
