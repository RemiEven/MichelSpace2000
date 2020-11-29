package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/RemiEven/michelSpace2000/src/ms2k"
	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
)

func main() {
	var seed string
	flag.StringVar(&seed, "seed", "", "used to seed RNG")

	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("MichelSpace2000")

	rng, err := rng.NewRNG(seed)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize rng: %w", err))
	}

	game := &ms2k.Game{
		World: ms2k.NewWorld(rng),
	}

	if err := game.Init(); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
