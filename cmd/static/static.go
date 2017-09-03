package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/mpihlak/maze"
)

func main() {
	var asciiArtLevel = "" +
		"#@#################\n" +
		"# #####           #\n" +
		"#     ## ######   #\n" +
		"##### #         # #\n" +
		"#     #  # #  # # #\n" +
		"# #####  # #### # #\n" +
		"#        #      # #\n" +
		"#################=#\n"

	scanner := bufio.NewScanner(strings.NewReader(asciiArtLevel))
	level := maze.ReadLevel(scanner)
	movez := make(maze.WalkerChannel)

	render := maze.NewStreamRenderer()
	defer render.Done()

	for _, actor := range level.Actors {
		actor.EndPos = level.Exits[0]
		go maze.WalkThrough(level, actor, &maze.ShortestPathWalker{}, movez)
	}

	done := false
	for counter := 0; !done; counter++ {
		select {
		case a := <-movez:
			if a.HasFinished() {
				done = true
			}
		}

		maze.Render(level, fmt.Sprintf("Iteration #%d", counter), render)
		counter++
	}
}
