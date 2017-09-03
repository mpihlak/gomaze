// Generate a maze, compute the shortest path through it and render to stdout.
// Optionally pass [width][height] as argv to control the dimensions of the maze.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/mpihlak/maze"
)

func main() {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	width, height := 40, 20

	if len(os.Args) > 1 {
		if w, err := strconv.Atoi(os.Args[1]); err == nil {
			width = w
		}
	}
	if len(os.Args) > 2 {
		if h, err := strconv.Atoi(os.Args[2]); err == nil {
			height = h
		}
	}

	render := maze.NewStreamRenderer()
	defer render.Done()

	level := maze.GenerateRandomMaze(width, height-1)
	actor := maze.Actor{Character: '&', CurrPos: level.Exits[0], EndPos: level.Exits[1]}

	level.AddActor(&actor)
	maze.CalculateShortestPath(level, &actor, level.Exits[1])
	maze.Render(level, fmt.Sprintf("Seed=%v Shortest path length=%v.", seed, len(actor.Path)), render)
}
