package main

import (
	"log"

	"github.com/RemiEven/michelSpace2000/src/ms2k"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("MichelSpace2000")

	game := &ms2k.Game{
		World: ms2k.NewWorld(),
	}

	if err := game.Init(); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
