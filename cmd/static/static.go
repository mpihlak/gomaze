package main

import (
	"bufio"
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

	render := maze.NewTermboxRenderer()
	defer render.Done()

  // Exits are not marked on the level, so we just make one up
	for _, actor := range level.Actors {
		actor.EndPos = level.Exits[0]
    actor.PathNav = &maze.ShortestPathWalker{}
	}

  controller := maze.NewController(&level, render)
  controller.Start()
  for controller.RunLoop() { }
  controller.Done()

}
