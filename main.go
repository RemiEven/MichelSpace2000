package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/RemiEven/michelSpace2000/src/ms2k"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	shipsImg, _, err := ebitenutil.NewImageFromFile("./assets/img/modular_ships.png")
	if err != nil {
		log.Fatal(err)
	}
	shipImg := shipsImg.SubImage(image.Rect(
		80,
		320,
		112,
		352,
	)).(*ebiten.Image)
	if err != nil {
		log.Fatal(err)
	}
	planetImg, _, err := ebitenutil.NewImageFromFile("./assets/img/Green Gas Planet.png")
	if err != nil {
		log.Fatal(err)
	}
	bgImg, _, err := ebitenutil.NewImageFromFile("./assets/img/back.png")
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("BLBLBLBLBLBLBL")

	audioContext := ms2k.NewAudioContext()

	music, err := ioutil.ReadFile("./assets/audio/Hardmoon_-_Deep_space.mp3")
	if err != nil {
		log.Fatal(err)
	}
	player, err := ms2k.NewPlayer(audioContext, music)
	defer player.Close()

	if err = ebiten.RunGame(&ms2k.Game{
		Images: map[string]*ebiten.Image{
			"ship":   shipImg,
			"planet": planetImg,
			"bg":     bgImg,
		},
		World: ms2k.NewWorld(),
	}); err != nil {
		log.Fatal(err)
	}
}
