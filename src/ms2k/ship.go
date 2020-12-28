package ms2k

import (
	"fmt"
	"math"
)

// Position contains X and Y coordinates of a position in space
type Position struct {
	X, Y float64
}

// DistanceTo measures the distance between the position and another
func (p *Position) DistanceTo(other *Position) float64 {
	return math.Hypot(p.X-other.X, p.Y-other.Y)
}

// String returns a text representation of the position
func (p *Position) String() string {
	return fmt.Sprintf("X: %8.0f parsecs\nY: %8.0f parsecs", p.X/10, p.Y/10)
}

// Direction is an enum used to know towards where an oriented object is pointing
type Direction int

// Enum of all supported directions
const (
	North Direction = iota
	Northwest
	West
	Southwest
	South
	Southeast
	East
	Northeast
)

// Ship holds the position and direction of a spaceship
type Ship struct {
	Position  Position
	Direction Direction

	PlanetScans map[*Planet]*Operation
}
