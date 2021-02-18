package ms2k

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
)

const (
	screenWidth  = 640 * 2
	screenHeight = 400 * 2

	viewportBorderMargin = 32 // should be equal or bigger than half the side length of the biggest sprite to avoid clipping
)

var (
	zoomFactor = 1.0
)

// Game contains all loaded game assets with current game data
type Game struct {
	assetLibrary *AssetLibrary

	state int8

	menu *MainMenu

	World *World
	score int

	lose *Operation
}

// Init initializes a game
func (g *Game) Init() error {
	assetLibrary, err := NewAssetLibrary()
	if err != nil {
		return fmt.Errorf("failed to load asset library: %w", err)
	}
	g.assetLibrary = assetLibrary

	audioContext := NewAudioContext()

	if _, err := NewPlayer(audioContext, g.assetLibrary.sounds["music"]); err != nil {
		return fmt.Errorf("failed to play music: %w", err)
	}

	g.state = stateInMenu
	g.menu = &MainMenu{}

	return nil
}

// Update is used to implement the ebiten.Game interface
func (g *Game) Update() error {
	timeNow := time.Now()

	switch g.state {
	case stateInMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch g.menu.state() {
			case menuStateNewGame:
				rng, err := rng.NewRNG("") // TODO: show a text input to have a seed
				if err != nil {
					log.Fatal(fmt.Errorf("failed to initialize rng: %w", err))
				}
				g.World = NewWorld(rng)

				g.state = stateInGame
				g.score = 0
				g.lose = &Operation{
					lastUpdate: timeNow,
					speed:      1. / 30.,
				}
			case menuStateExit:
				os.Exit(0)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.menu.selectedIndex++
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.menu.selectedIndex--
		}
	case stateInGame:
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.World.selectPreviousShip()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			g.World.selectNextShip()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			zoomFactor = zoomFactor * 2
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			zoomFactor = zoomFactor / 2
		}

		var (
			up    = ebiten.IsKeyPressed(ebiten.KeyUp)
			down  = ebiten.IsKeyPressed(ebiten.KeyDown)
			left  = ebiten.IsKeyPressed(ebiten.KeyLeft)
			right = ebiten.IsKeyPressed(ebiten.KeyRight)
		)

		var (
			goesSouth = down && !up
			goesNorth = up && !down
			goesWest  = left && !right
			goesEast  = right && !left
		)

		selectedShip := g.World.getSelectedShip()
		speed := 3.0
		if goesNorth {
			selectedShip.Position.Y -= speed
		}
		if goesSouth {
			selectedShip.Position.Y += speed
		}
		if goesWest {
			selectedShip.Position.X -= speed
		}
		if goesEast {
			selectedShip.Position.X += speed
		}

		switch {
		case goesNorth && goesWest:
			selectedShip.Direction = Northwest
		case goesWest && goesSouth:
			selectedShip.Direction = Southwest
		case goesSouth && goesEast:
			selectedShip.Direction = Southeast
		case goesEast && goesNorth:
			selectedShip.Direction = Northeast
		case goesNorth:
			selectedShip.Direction = North
		case goesWest:
			selectedShip.Direction = West
		case goesSouth:
			selectedShip.Direction = South
		case goesEast:
			selectedShip.Direction = East
		}

		g.World.ensureChunksAroundAreGenerated(selectedShip.Position)

		for _, ship := range g.World.Ships {
			for _, planet := range g.World.Planets {
				if !planet.Looted && ship.Position.DistanceTo(&planet.Position) < 50 {
					if _, ok := ship.PlanetScans[planet]; !ok {
						ship.PlanetScans[planet] = &Operation{
							lastUpdate: timeNow,
							speed:      50,
						}
					}
				}
			}
			for planet, scan := range ship.PlanetScans {
				if !planet.Looted && ship.Position.DistanceTo(&planet.Position) < 50 {
					scan.Update(timeNow)
					if scan.IsCompleted() {
						delete(ship.PlanetScans, planet)
						g.score++
						planet.Looted = true
					}
				} else {
					delete(ship.PlanetScans, planet)
				}
			}
		}

		if g.score >= 10 {
			g.state = stateWon
			g.World = nil
			return nil
		}
		g.lose.Update(timeNow)
		if g.lose.IsCompleted() {
			g.World = nil
			g.state = stateLost
		}
	case stateLost, stateWon:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = stateInMenu
			g.menu.selectedIndex = menuStateNewGame
		}
	}

	return nil
}

var (
	white     = color.White
	lightBlue = color.RGBA{R: 0xaa, G: 0xaa, B: 0xff, A: 0xff}

	textColor         = white
	selectedTextColor = lightBlue
)

// Draw is used to implement the ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	fontFace := g.assetLibrary.fontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()

	switch g.state {
	case stateInMenu:
		{
			titleLabel := "MichelSpace2000"
			boundString := text.BoundString(fontFace, titleLabel)
			text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*2, textColor)
		}

		color := func(menuOption int8) color.Color {
			if menuOption == g.menu.state() {
				return selectedTextColor
			}
			return textColor
		}

		{
			newGameLabel := "New game"
			boundString := text.BoundString(fontFace, newGameLabel)
			text.Draw(screen, newGameLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*4, color(menuStateNewGame))
		}
		{
			exitLabel := "Exit"
			boundString := text.BoundString(fontFace, exitLabel)
			text.Draw(screen, exitLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*5, color(menuStateExit))
		}
	case stateInGame:
		viewPortCenter := g.World.getSelectedShip().Position

		scale := 1.0
		parallaxFactor := math.Pow(3.0, zoomFactor)
		imageWidth, imageHeight := g.assetLibrary.images["bg"].Size()
		topLeftBackgroundTileX := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*viewPortCenter.X - screenWidth/2 /*/zoomFactor*/) / (float64(imageWidth) * scale)))
		topLeftBackgroundTileY := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*viewPortCenter.Y - screenHeight/2 /*/zoomFactor*/) / (float64(imageHeight) * scale)))
		bottomRightBackgroundTileX := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*viewPortCenter.X + screenWidth/2 /*/zoomFactor*/) / (float64(imageWidth) * scale)))
		bottomRightBackgroundTileY := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*viewPortCenter.Y + screenHeight/2 /*/zoomFactor*/) / (float64(imageHeight) * scale)))
		x := topLeftBackgroundTileX
		for x <= bottomRightBackgroundTileX {
			y := topLeftBackgroundTileY
			for y <= bottomRightBackgroundTileY {
				dio := &ebiten.DrawImageOptions{}
				dio.GeoM.Translate(float64(x)*scale*float64(imageWidth)+screenWidth/2.0 /*/zoomFactor*/ -(parallaxFactor-1.0)*viewPortCenter.X/parallaxFactor, float64(y)*scale*float64(imageHeight)+screenHeight/2.0 /*/zoomFactor*/ -(parallaxFactor-1.0)*viewPortCenter.Y/parallaxFactor)
				screen.DrawImage(g.assetLibrary.images["bg"], dio)
				y++
			}
			x++
		}

		minXToDisplay := viewPortCenter.X - (screenWidth/2/zoomFactor + viewportBorderMargin)
		maxXToDisplay := viewPortCenter.X + (screenWidth/2/zoomFactor + viewportBorderMargin)
		minYToDisplay := viewPortCenter.Y - (screenHeight/2/zoomFactor + viewportBorderMargin)
		maxYToDisplay := viewPortCenter.Y + (screenHeight/2/zoomFactor + viewportBorderMargin)

		{
			imageWidth, imageHeight := g.assetLibrary.images["planet"].Size()
			for _, planet := range g.World.Planets {
				if isInBox(planet.Position.X, planet.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
					dio := &ebiten.DrawImageOptions{}
					scale := 0.25 * zoomFactor
					dio.GeoM.Scale(scale, scale)
					dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

					translateToDrawPosition(planet.Position, viewPortCenter, &dio.GeoM, zoomFactor)

					dio.ColorM.ChangeHSV(planet.Hue, 1, 1)
					if planet.Looted {
						dio.ColorM.ChangeHSV(0, 0, 1)
					}
					screen.DrawImage(g.assetLibrary.images["planet"], dio)
				}
			}
		}

		{
			imageWidth, imageHeight := g.assetLibrary.images["ship"].Size()
			for _, ship := range g.World.Ships {
				if isInBox(ship.Position.X, ship.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
					dio := &ebiten.DrawImageOptions{}
					scale := 1.0 * zoomFactor
					dio.GeoM.Scale(scale, scale)
					dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)
					dio.GeoM.Rotate(-2.0 * math.Pi / 8.0 * float64(ship.Direction))

					translateToDrawPosition(ship.Position, viewPortCenter, &dio.GeoM, zoomFactor)

					screen.DrawImage(g.assetLibrary.images["ship"], dio)
				}
			}
		}

		textBgColor := color.Black
		ebitenutil.DrawRect(screen, 0, 0, 250, 120, textBgColor)

		text.Draw(screen, strconv.Itoa(g.score)+"/"+strconv.Itoa(10), fontFace, 0, 26, textColor)
		text.Draw(screen, g.World.getSelectedShip().Position.String(), fontFace, 0, 54, textColor)
		text.Draw(screen, strconv.Itoa(int(minXToDisplay)), fontFace, 0, 110, textColor)
	case stateLost:
		text.Draw(screen, "Game over", fontFace, 0, 26, textColor)
	case stateWon:
		text.Draw(screen, "Victory", fontFace, 0, 26, textColor)
	}

}

func translateToDrawPosition(gamePosition, viewPortCenter Position, geoM *ebiten.GeoM, zoomFactor float64) {
	geoM.Translate(-viewPortCenter.X*zoomFactor, -viewPortCenter.Y*zoomFactor)
	geoM.Translate(screenWidth/2, screenHeight/2)
	geoM.Translate(gamePosition.X*zoomFactor, gamePosition.Y*zoomFactor)
}

// Layout is used to implement the ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func isInBox(x, y, minX, maxX, minY, maxY float64) bool {
	return minX <= x && x <= maxX && minY <= y && y <= maxY
}
