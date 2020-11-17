package ms2k

import (
	"math"
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"
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

	noise opensimplex.Noise32
}

// NewWorld creates a new world
func NewWorld() *World {
	ship1 := &Ship{}
	ship2 := &Ship{}

	noise := opensimplex.NewNormalized32(rand.Int63()) // TODO: extract seed

	planets := make([]*Planet, 0)

	return &World{
		Planets:         planets,
		Ships:           []*Ship{ship1, ship2},
		noise:           noise,
		GeneratedChunks: map[int]map[int]struct{}{},
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
			if w.noise.Eval2(float32(i), float32(j)) > 0.87 {
				planet := NewGasPlanet(Position{
					X: cellSize*float64(i) + float64(x*cellSize*chunkSize),
					Y: cellSize*float64(j) + float64(y*cellSize*chunkSize),
				})
				w.Planets = append(w.Planets, planet)
			}
		}
	}
}

func getChunkContaining(p Position) (int, int) {
	return int(math.Floor(p.X / (cellSize * chunkSize))), int(math.Floor(p.Y / (cellSize * chunkSize)))
}
