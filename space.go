package main

type World struct {
	Planets           []*Planet
	Ships             []*Ship
	selectedShipIndex int
}

func NewWorld() *World {
	ship1 := &Ship{}
	ship2 := &Ship{}

	return &World{
		Planets: []*Planet{
			{
				Position: Position{
					X: 150,
					Y: 200,
				},
			},
			{
				Position: Position{
					X: 30,
					Y: 100,
				},
			},
			{
				Position: Position{
					X: 40,
					Y: 180,
				},
			},
		},
		Ships: []*Ship{ship1, ship2},
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
