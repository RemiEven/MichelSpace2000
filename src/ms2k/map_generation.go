package ms2k

import "math"

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
			if value := w.rng.GetValueAtPosition(float32(i+x*chunkSize), float32(j+y*chunkSize)); value >= 0.9 {
				planet := &Planet{
					Position: Position{
						X: cellSize*float64(i) + float64(x*cellSize*chunkSize),
						Y: cellSize*float64(j) + float64(y*cellSize*chunkSize),
					},
					Hue: float64(w.rng.GetValueAtPosition(-float32(i+x*chunkSize)/20, -float32(j+y*chunkSize)/20) * 2 * math.Pi),
				}

				if value >= 0.95 {
					planet.AddMoon(float64((value - 0.95) * 100))
				}

				w.Planets = append(w.Planets, planet)
			}
		}
	}
}

func getChunkContaining(p Position) (int, int) {
	return int(math.Floor(p.X / (cellSize * chunkSize))), int(math.Floor(p.Y / (cellSize * chunkSize)))
}
