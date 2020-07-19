package main

import "math"

type Position struct {
	X, Y float64
}

func (p *Position) DistanceTo(other *Position) float64 {
	return math.Sqrt(math.Pow(p.X-other.X, 2) + math.Pow(p.Y-other.Y, 2))
}

type Direction int

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

type Ship struct {
	Position  Position
	Direction Direction
}
