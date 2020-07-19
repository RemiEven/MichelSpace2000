package main

import (
	"image"
	_ "image/png"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 400
)

type Game struct {
	images map[string]*ebiten.Image
	world  *World
	score  int
	won    bool
}

func (g *Game) Update(_ *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.world.selectPreviousShip()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.world.selectNextShip()
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.world.getSelectedShip().Position.Y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.world.getSelectedShip().Position.Y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.world.getSelectedShip().Position.X--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.world.getSelectedShip().Position.X++
	}

	for _, ship := range g.world.Ships {
		for _, planet := range g.world.Planets {
			if !planet.Looted && ship.Position.DistanceTo(&planet.Position) < 50 {
				planet.Looted = true
				g.score++
			}
		}
	}

	if g.score == len(g.world.Planets) {
		g.won = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.won {
		ebitenutil.DebugPrint(screen, "victory")
		return
	}

	viewPortCenter := g.world.getSelectedShip().Position

	for _, planet := range g.world.Planets {
		dio := &ebiten.DrawImageOptions{}
		scale := 0.25
		dio.GeoM.Scale(scale, scale)
		dio.GeoM.Translate(-viewPortCenter.X, -viewPortCenter.Y)
		dio.GeoM.Translate(screenWidth/2, screenHeight/2)
		dio.GeoM.Translate(planet.Position.X, planet.Position.Y)
		imageWidth, imageHeight := g.images["planet"].Size()
		dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)
		if planet.Looted {
			dio.ColorM.ChangeHSV(0, 0, 1)
		}
		screen.DrawImage(g.images["planet"], dio)
	}
	for _, ship := range g.world.Ships {
		dio := &ebiten.DrawImageOptions{}
		scale := 1.0
		dio.GeoM.Scale(scale, scale)
		dio.GeoM.Translate(-viewPortCenter.X, -viewPortCenter.Y)
		dio.GeoM.Translate(screenWidth/2, screenHeight/2)
		dio.GeoM.Translate(ship.Position.X, ship.Position.Y)
		imageWidth, imageHeight := g.images["ship"].Size()
		dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)
		screen.DrawImage(g.images["ship"], dio)
	}

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.score))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	shipsImg, _, err := ebitenutil.NewImageFromFile("./modular_ships.png", ebiten.FilterDefault)
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
	planetImg, _, err := ebitenutil.NewImageFromFile("./Green Gas Planet.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("BLBLBLBLBLBLBL")
	if err = ebiten.RunGame(&Game{
		images: map[string]*ebiten.Image{
			"ship":   shipImg,
			"planet": planetImg,
		},
		world: NewWorld(),
	}); err != nil {
		log.Fatal(err)
	}
}
