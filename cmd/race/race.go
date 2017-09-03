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

	movez := make(maze.WalkerChannel)

	level.AddActor(&maze.Actor{Character: '@', CurrPos: level.Exits[0], EndPos: level.Exits[1]})
	level.AddActor(&maze.Actor{Character: '&', CurrPos: level.Exits[1], EndPos: level.Exits[0]})

	for _, w := range level.Actors {
		go maze.WalkThrough(level, w, &maze.ShortestLineWalker{}, movez)
	}

	done := false
	for frame := 1; !done; frame++ {
		select {
		case a := <-movez:
			if a.HasFinished() {
				done = true
			}
		case k := <-render.GetKeyboardEvent():
			if k == maze.KBEventCancel {
				done = true
			}
		}
		maze.Render(level, fmt.Sprintf("render #%d", frame), render)
	}

	maze.Render(level, "Woohoo! Press any key to exit...", render)
	<-render.GetKeyboardEvent()
}
