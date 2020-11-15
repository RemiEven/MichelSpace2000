package main

import (
	"math"
	"math/rand"
)

// Planet holds all information about a planet
type Planet struct {
	Position Position
	Looted   bool
	Hue      float64
}

// NewGasPlanet creates a new planet of type Gas with a random Hue placed on the given position
func NewGasPlanet(position Position) *Planet {
	return &Planet{
		Position: position,
		Hue:      rand.Float64() * 2 * math.Pi,
	}
}
