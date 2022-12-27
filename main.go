package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/RemiEven/michelSpace2000/src/ms2k"
)

func main() {

	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("MichelSpace2000")

	game := &ms2k.Game{}

	if err := game.Init(); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
