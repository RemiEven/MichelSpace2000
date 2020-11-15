package main

import (
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"
)

// World contains data such as all the Planets & Ships of the game
type World struct {
	Planets           []*Planet
	Ships             []*Ship
	selectedShipIndex int

	noise opensimplex.Noise32
}

// NewWorld creates a new world
func NewWorld() *World {
	ship1 := &Ship{}
	ship2 := &Ship{}

	noise := opensimplex.NewNormalized32(rand.Int63()) // TODO: extract seed

	planets := make([]*Planet, 0)
	for i := -10; i <= 10; i++ {
		for j := -10; j <= 10; j++ {
			if noise.Eval2(float32(i), float32(j)) > 0.8 {
				planet := NewGasPlanet(Position{
					X: 50.0 * float64(i),
					Y: 50.0 * float64(j),
				})
				planets = append(planets, planet)
			}
		}
	}

	return &World{
		Planets: planets,
		Ships:   []*Ship{ship1, ship2},
		noise:   noise,
	}
}

func (w *World) getSelectedShip() *Ship {
	return w.Ships[w.selectedShipIndex]
}

func (w *World) selectNextShip() {
	w.selectedShipIndex = (w.selectedShipIndex + 1) % len(w.Ships)
}

func (w *World) selectPreviousShip() {
	w.selectedShipIndex = (w.selectedShipIndex + len(w.Ships) - 1) % len(w.Ships)
}
