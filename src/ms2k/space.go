package ms2k

import (
	"math"

	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
)

const (
	cellSize  = 50
	chunkSize = 32
)

// World contains data such as all the Planets & Ships of the game
type World struct {
	Planets           []*Planet
	GeneratedChunks   map[int]map[int]struct{}
	Ships             []*Ship
	selectedShipIndex int

	rng *rng.RNG
}

// NewWorld creates a new world
func NewWorld(rng *rng.RNG) *World {
	ship1 := &Ship{}
	ship2 := &Ship{}

	planets := make([]*Planet, 0)

	return &World{
		Planets:         planets,
		Ships:           []*Ship{ship1, ship2},
		GeneratedChunks: map[int]map[int]struct{}{},

		rng: rng,
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

func (w *World) ensureChunksAroundAreGenerated(p Position) {
	x0, y0 := getChunkContaining(p)
	for x := x0 - 1; x <= x0+1; x++ {
		for y := y0 - 1; y <= y0+1; y++ {
			if w.GeneratedChunks[x] == nil {
				w.GeneratedChunks[x] = map[int]struct{}{}
			}
			if _, ok := w.GeneratedChunks[x][y]; !ok {
				w.generateChunk(x, y)
				w.GeneratedChunks[x][y] = struct{}{}
			}
		}
	}
}

func (w *World) generateChunk(x, y int) {
	for i := 0; i < chunkSize; i++ {
		for j := 0; j < chunkSize; j++ {
			if w.rng.GetValueAtPosition(float32(i+x*chunkSize), float32(j+y*chunkSize)) >= 0.9 {
				planet := &Planet{
					Position: Position{
						X: cellSize*float64(i) + float64(x*cellSize*chunkSize),
						Y: cellSize*float64(j) + float64(y*cellSize*chunkSize),
					},
					Hue: float64(w.rng.GetValueAtPosition(-float32(i+x*chunkSize)/20, -float32(j+y*chunkSize)/20) * 2 * math.Pi),
				}
				w.Planets = append(w.Planets, planet)
			}
		}
	}
}

func getChunkContaining(p Position) (int, int) {
	return int(math.Floor(p.X / (cellSize * chunkSize))), int(math.Floor(p.Y / (cellSize * chunkSize)))
}
