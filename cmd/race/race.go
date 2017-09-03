// Generate a random maze and race some walkers through it
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mpihlak/maze"
)

func main() {
	var seed int64

	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &seed)
	} else {
		seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(seed)

	render := maze.NewTermboxRenderer()
	defer render.Done()
  width, height := render.Size()
	level := maze.GenerateRandomMaze(width, height-1)

  a1 := maze.NewActor('@', level.Exits[0], level.Exits[1], &maze.ShortestLineWalker{})
	level.AddActor(a1)

  a2 := maze.NewActor('&', level.Exits[1], level.Exits[0], &maze.ShortestPathWalker{})
	level.AddActor(a2)

  controller := maze.NewController(&level, render)
  controller.Start()

  for controller.RunLoop() {
    // RunLoop takes care of rendering and keyboard events.
    // Do whatever you want here.
  }

  controller.Done()
}
