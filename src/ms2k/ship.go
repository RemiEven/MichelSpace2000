package ms2k

// Ship holds the position and direction of a spaceship
type Ship struct {
	Position  Position
	Direction Direction

	PlanetScans map[*Planet]*Operation

	Radar *Radar
}
