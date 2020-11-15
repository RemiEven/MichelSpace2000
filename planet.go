package main

import (
	"math"
	"math/rand"
)

type Planet struct {
	Position Position
	Looted   bool
	Hue      float64
}

func NewGasPlanet() *Planet {
	return &Planet{
		// Position: Position{
		// 	X: rand.Float64()*worldSize - worldSize/2,
		// 	Y: rand.Float64()*worldSize - worldSize/2,
		// },
		Hue: rand.Float64() * 2 * math.Pi,
	}
}
