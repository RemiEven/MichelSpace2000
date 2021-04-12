package ms2k

import (
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
)

const (
	cellSize  = 50
	chunkSize = 32
)

// World contains data such as all the Planets & Ships of the game
type World struct {
	Planets           []*Planet
	WormHoles         []*WormHole
	GeneratedChunks   map[int]map[int]struct{}
	Ships             []*Ship
	selectedShipIndex int

	rng *rng.RNG

	score int

	lose *Operation
}

// NewWorld creates a new world
func NewWorld(rng *rng.RNG, timeNow time.Time) *World {
	ship1 := &Ship{
		PlanetScans: map[*Planet]*Operation{},
	}
	ship2 := &Ship{
		PlanetScans: map[*Planet]*Operation{},
	}

	planets := make([]*Planet, 0)

	return &World{
		Planets:         planets,
		Ships:           []*Ship{ship1, ship2},
		GeneratedChunks: map[int]map[int]struct{}{},

		lose: &Operation{
			lastUpdate: timeNow,
			speed:      4. / 5.,
		},
		rng: rng,
	}
}

// Update updates the world
func (w *World) Update(timeNow time.Time, settings *Settings) int8 {
	if inpututil.IsKeyJustPressed(keyMapping.PreviousShip) {
		w.selectPreviousShip()
	}
	if inpututil.IsKeyJustPressed(keyMapping.NextShip) {
		w.selectNextShip()
	}

	if inpututil.IsKeyJustPressed(keyMapping.ZoomIn) {
		zoomFactor = zoomFactor * 2
	}
	if inpututil.IsKeyJustPressed(keyMapping.ZoomOut) {
		zoomFactor = zoomFactor / 2
	}

	var (
		up    = ebiten.IsKeyPressed(keyMapping.Up)
		down  = ebiten.IsKeyPressed(keyMapping.Down)
		left  = ebiten.IsKeyPressed(keyMapping.Left)
		right = ebiten.IsKeyPressed(keyMapping.Right)
	)

	var (
		goesSouth = down && !up
		goesNorth = up && !down
		goesWest  = left && !right
		goesEast  = right && !left
	)

	selectedShip := w.getSelectedShip()
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

	w.ensureChunksAroundAreGenerated(selectedShip.Position)

	for _, ship := range w.Ships {
		for _, planet := range w.Planets {
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
					w.score++
					planet.Looted = true
				}
			} else {
				delete(ship.PlanetScans, planet)
			}
		}
	}

	if w.score >= 10 {
		return stateWon
	}
	w.lose.Update(timeNow)
	if w.lose.IsCompleted() {
		return stateLost
	}

	return stateInGame
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

// Draw draws the world
func (w *World) Draw(screen *ebiten.Image, assetLibrary *AssetLibrary) {
	fontFace := assetLibrary.fontFaces["oxanium"]

	viewPortCenter := w.getSelectedShip().Position

	scale := 1.0
	parallaxFactor := math.Pow(3.0, zoomFactor)
	imageWidth, imageHeight := assetLibrary.images["bg"].Size()
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
			screen.DrawImage(assetLibrary.images["bg"], dio)
			y++
		}
		x++
	}

	minXToDisplay := viewPortCenter.X - (screenWidth/2/zoomFactor + viewportBorderMargin)
	maxXToDisplay := viewPortCenter.X + (screenWidth/2/zoomFactor + viewportBorderMargin)
	minYToDisplay := viewPortCenter.Y - (screenHeight/2/zoomFactor + viewportBorderMargin)
	maxYToDisplay := viewPortCenter.Y + (screenHeight/2/zoomFactor + viewportBorderMargin)

	{
		imageWidth, imageHeight := assetLibrary.images["wormHole"].Size()
		for _, wormHole := range w.WormHoles {
			if isInBox(wormHole.Position.X, wormHole.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
				dio := &ebiten.DrawImageOptions{}
				scale := 2 * zoomFactor
				dio.GeoM.Scale(scale, scale)
				dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

				translateToDrawPosition(wormHole.Position, viewPortCenter, &dio.GeoM, zoomFactor)

				screen.DrawImage(assetLibrary.images["wormHole"], dio)
			}
		}
	}

	{
		imageWidth, imageHeight := assetLibrary.images["planet"].Size()
		for _, planet := range w.Planets {
			if isInBox(planet.Position.X, planet.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
				dio := &ebiten.DrawImageOptions{}
				scale := 0.25 * zoomFactor
				dio.GeoM.Scale(scale, scale)
				dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

				translateToDrawPosition(planet.Position, viewPortCenter, &dio.GeoM, zoomFactor)

				dio.ColorM.ChangeHSV(planet.Hue, 1, 1)
				screen.DrawImage(assetLibrary.images["planet"], dio)
				moonImageWidth, moonImageHeight := assetLibrary.images["moon"].Size()
				for _, moon := range planet.Moons {
					dio := &ebiten.DrawImageOptions{}
					scale := zoomFactor
					dio.GeoM.Scale(scale, scale)
					dio.GeoM.Translate(-float64(moonImageWidth)/2.0*scale, -float64(moonImageHeight)/2.0*scale)
					translateToDrawPosition(moon.Position, viewPortCenter, &dio.GeoM, zoomFactor)
					screen.DrawImage(assetLibrary.images["moon"], dio)
				}
				if planet.Looted {
					satelliteImageWidth, satelliteImageHeight := assetLibrary.images["satellite"].Size()
					dio := &ebiten.DrawImageOptions{}
					scale := zoomFactor
					dio.GeoM.Scale(scale, scale)
					dio.GeoM.Translate(-float64(satelliteImageWidth)/2.0*scale, -float64(satelliteImageHeight)/2.0*scale)

					distance := 38
					position := Position{
						X: planet.Position.X + math.Sqrt2*float64(distance/2),
						Y: planet.Position.Y - math.Sqrt2*float64(distance/2),
					}
					translateToDrawPosition(position, viewPortCenter, &dio.GeoM, zoomFactor)
					screen.DrawImage(assetLibrary.images["satellite"], dio)
				}
			}
		}
	}

	{
		imageWidth, imageHeight := assetLibrary.images["earth"].Size()
		if isInBox(0, 0, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
			dio := &ebiten.DrawImageOptions{}
			scale := 2.0 * zoomFactor
			dio.GeoM.Scale(scale, scale)
			dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

			translateToDrawPosition(Position{}, viewPortCenter, &dio.GeoM, zoomFactor)

			screen.DrawImage(assetLibrary.images["earth"], dio)
		}
	}

	{
		imageWidth, imageHeight := assetLibrary.images["ship"].Size()
		for _, ship := range w.Ships {
			if isInBox(ship.Position.X, ship.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
				dio := &ebiten.DrawImageOptions{}
				scale := 1.0 * zoomFactor
				dio.GeoM.Scale(scale, scale)
				dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)
				dio.GeoM.Rotate(-2.0 * math.Pi / 8.0 * float64(ship.Direction))

				translateToDrawPosition(ship.Position, viewPortCenter, &dio.GeoM, zoomFactor)

				screen.DrawImage(assetLibrary.images["ship"], dio)
			}
		}
	}

	textBgColor := color.Black
	ebitenutil.DrawRect(screen, 0, 0, 250, 120, textBgColor)

	text.Draw(screen, strconv.Itoa(w.score)+"/"+strconv.Itoa(10), fontFace, 4, 26, textColor)
	text.Draw(screen, w.getSelectedShip().Position.String(), fontFace, 4, 54, textColor)
	text.Draw(screen, loseOperationToDoomsdayClockTime(w.lose), fontFace, 4, 110, textColor)
}

func loseOperationToDoomsdayClockTime(operation *Operation) string {
	switch {
	case operation.completedPercentage < 20:
		return "23:55"
	case operation.completedPercentage < 40:
		return "23:56"
	case operation.completedPercentage < 60:
		return "23:57"
	case operation.completedPercentage < 80:
		return "23:58"
	default:
		return "23:59"
	}
}
