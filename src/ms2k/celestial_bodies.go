package ms2k

import "math"

// Planet holds all information about a planet
type Planet struct {
	Name     string
	Position Position
	Looted   bool
	Hue      float64
	Moons    []*Moon
}

func (planet *Planet) AddMoon(angle float64) {
	distance := 45.0
	planet.Moons = append(planet.Moons, &Moon{
		Position: Position{
			X: planet.Position.X + distance*math.Cos(angle),
			Y: planet.Position.Y + distance*math.Sin(angle),
		},
	})
}

// Moon holds all information about a moon
type Moon struct {
	Position Position
}

// WormHole holds all information about a worm hole
type WormHole struct {
	Position Position
}
