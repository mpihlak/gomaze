// Package maze, controlling the moving objects
package maze

import (
	"fmt"
	"time"
)

// Controller manages the actors on the level, responds to keyboard events and
// renders the maze.
type Controller struct {
	frame  int
	level  *Level
	render Renderer
}

func NewController(level *Level, render Renderer) *Controller {
	c := Controller{level: level, render: render}
	return &c
}

func (c *Controller) Start() {
	for _, actor := range c.level.Actors {
		actor.PathNav.Initialize(c.level, actor)
	}
}

// RunLoop is called in a loop to update the state of the moving objects,
// render the maze and collect keyboard events.
func (c *Controller) RunLoop() bool {
	Render(*c.level, fmt.Sprintf("render #%d", c.frame), c.render)

	isDone := false

	for _, actor := range c.level.Actors {
		actor.PathNav.NextPosition()
		if actor.HasFinished() {
			isDone = true
		}
	}

	select {
	case k := <-c.render.GetKeyboardEvent():
		if k == KBEventCancel {
			isDone = true
		}
	default:
	}

	time.Sleep(100 * time.Millisecond)
	c.frame++

	return !isDone
}

func (c *Controller) Done() {
	Render(*c.level, fmt.Sprintf("Woohoo! Done after %d iterations. Press any key to exit...", c.frame), c.render)
	<-c.render.GetKeyboardEvent()
}
