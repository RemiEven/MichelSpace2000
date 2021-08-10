package ms2k

import (
	"math"
	"strconv"
)

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
			if value := w.rng.GetValueAtPosition(float32(i+x*chunkSize), float32(j+y*chunkSize)); value >= 0.92 {
				planet := &Planet{
					Name: toPlanetName(value),
					Position: Position{
						X: cellSize*float64(i) + float64(x*cellSize*chunkSize),
						Y: cellSize*float64(j) + float64(y*cellSize*chunkSize),
					},
					Hue: float64(w.rng.GetValueAtPosition(-float32(i+x*chunkSize)/20, -float32(j+y*chunkSize)/20) * 2 * math.Pi),
				}

				if value >= 0.96 {
					planet.AddMoon(float64((value - 0.96) / (1.0 - 0.96) * 4.0 * math.Pi))
				}

				w.Planets = append(w.Planets, planet)
			} else if value < 0.02 {
				wormHole := &WormHole{
					Position: Position{
						X: cellSize*float64(i) + float64(x*cellSize*chunkSize),
						Y: cellSize*float64(j) + float64(y*cellSize*chunkSize),
					},
				}

				w.WormHoles = append(w.WormHoles, wormHole)
			}
		}
	}
}

func getChunkContaining(p Position) (int, int) {
	return int(math.Floor(p.X / (cellSize * chunkSize))), int(math.Floor(p.Y / (cellSize * chunkSize)))
}

func toPlanetName(number float32) string {
	n := int(number * 100_000_000)
	starCatalogue := [10]string{"GJ", "Kepler", "Corot", "HAT", "HD", "SAO", "FK", "YBS", "HIP", "LP"}[(n/1_000_000)%10]
	starNumber := (n / 10) % 100_000
	planetLetter := 'b' + (n % 10)
	return starCatalogue + " " + strconv.Itoa(starNumber) + " " + string([]rune{rune(planetLetter)})
}
