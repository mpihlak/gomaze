// Package maze, moving objects
package maze

// Actor is something that can move around on the level
type Actor struct {
	Character rune       // Display character
	CurrPos   Position   // Current location
	EndPos    Position   // Destination, if calculated
	Path      []Position // Path, if calculated.
	PathNav   Walker
}

// Walker specifies the interface that can be used to walk an Actor through the maze
type Walker interface {
	Initialize(level *Level, actor *Actor)
	NextPosition()
}

// NewActor creates and initializes and Actor
func NewActor(c rune, startPos, endPos Position, navigator Walker) *Actor {
	actor := Actor{Character: c, CurrPos: startPos, EndPos: endPos, PathNav: navigator}
	return &actor
}

// HasFinished returns true if the actor has reached its destination
func (a Actor) HasFinished() bool {
	return a.CurrPos == a.EndPos
}
