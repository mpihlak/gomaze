// Package maze ... walker interfaces.
package maze

import (
	"time"
)

// WalkerChannel is used by actors for posting movement updates
type WalkerChannel chan *Actor

// Walker specifies the interface that can be used to walk an Actor through the maze
type Walker interface {
	Initialize(level *Level, actor *Actor)
	NextPosition()
}

// WalkThrough walks the actor through the maze, pushing new positions to movez channel
func WalkThrough(level Level, actor *Actor, walker Walker, movez WalkerChannel) {
	walker.Initialize(&level, actor)

	for {
		movez <- actor
		time.Sleep(100 * time.Millisecond)
		walker.NextPosition()
	}
}
