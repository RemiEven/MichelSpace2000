package ms2k

import (
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
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

	displayedPlanetName string
}

// NewWorld creates a new world
func NewWorld(rng *rng.RNG, timeNow time.Time) *World {
	ship1 := &Ship{
		PlanetScans: map[*Planet]*Operation{},
	}
	ship2 := &Ship{
		PlanetScans: map[*Planet]*Operation{},
	}

	planets := make([]*Planet, 1)
	planets[0] = &Planet{
		Name:   "Earth",
		Looted: true,
	}

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
		var closestPlanet *Planet
		distanceToClosestPlanet := math.MaxFloat64
		for _, planet := range w.Planets {
			distanceToShip := ship.Position.DistanceTo(&planet.Position)
			if distanceToShip < distanceToClosestPlanet {
				distanceToClosestPlanet = distanceToShip
				closestPlanet = planet
			}
			if !planet.Looted && distanceToShip < 50 {
				if _, ok := ship.PlanetScans[planet]; !ok {
					ship.PlanetScans[planet] = &Operation{
						lastUpdate: timeNow,
						speed:      50,
					}
				}
			}
		}

		if ship == selectedShip {
			if distanceToClosestPlanet < 50 {
				w.displayedPlanetName = closestPlanet.Name
			} else {
				w.displayedPlanetName = ""
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
func (w *World) Draw(screen *ebiten.Image, assetLibrary *assets.Library) {
	fontFace := assetLibrary.FontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()
	fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

	viewPortCenter := w.getSelectedShip().Position

	drawSpaceBackground(screen, assetLibrary, viewPortCenter)

	screenBounds := screen.Bounds()
	screenWidth, screenHeight := float64(screenBounds.Dx()), float64(screenBounds.Dy())

	minXToDisplay := viewPortCenter.X - (screenWidth/2/zoomFactor + viewportBorderMargin)
	maxXToDisplay := viewPortCenter.X + (screenWidth/2/zoomFactor + viewportBorderMargin)
	minYToDisplay := viewPortCenter.Y - (screenHeight/2/zoomFactor + viewportBorderMargin)
	maxYToDisplay := viewPortCenter.Y + (screenHeight/2/zoomFactor + viewportBorderMargin)

	{
		imageWidth, imageHeight := assetLibrary.Images["wormHole"].Size()
		for _, wormHole := range w.WormHoles {
			if isInBox(wormHole.Position.X, wormHole.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
				dio := &ebiten.DrawImageOptions{}
				scale := 2 * zoomFactor
				dio.GeoM.Scale(scale, scale)
				dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

				translateToDrawPosition(&screenBounds, wormHole.Position, viewPortCenter, &dio.GeoM, zoomFactor)

				screen.DrawImage(assetLibrary.Images["wormHole"], dio)
			}
		}
	}

	{
		imageWidth, imageHeight := assetLibrary.Images["planet"].Size()
		for _, planet := range w.Planets {
			if planet == w.Planets[0] {
				continue
			}
			if isInBox(planet.Position.X, planet.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
				dio := &ebiten.DrawImageOptions{}
				scale := 0.25 * zoomFactor
				dio.GeoM.Scale(scale, scale)
				dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

				translateToDrawPosition(&screenBounds, planet.Position, viewPortCenter, &dio.GeoM, zoomFactor)

				dio.ColorM.ChangeHSV(planet.Hue, 1, 1)
				screen.DrawImage(assetLibrary.Images["planet"], dio)
				moonImageWidth, moonImageHeight := assetLibrary.Images["moon"].Size()
				for _, moon := range planet.Moons {
					dio := &ebiten.DrawImageOptions{}
					scale := zoomFactor
					dio.GeoM.Scale(scale, scale)
					dio.GeoM.Translate(-float64(moonImageWidth)/2.0*scale, -float64(moonImageHeight)/2.0*scale)
					translateToDrawPosition(&screenBounds, moon.Position, viewPortCenter, &dio.GeoM, zoomFactor)
					screen.DrawImage(assetLibrary.Images["moon"], dio)
				}
				if planet.Looted {
					satelliteImageWidth, satelliteImageHeight := assetLibrary.Images["satellite"].Size()
					dio := &ebiten.DrawImageOptions{}
					scale := zoomFactor
					dio.GeoM.Scale(scale, scale)
					dio.GeoM.Translate(-float64(satelliteImageWidth)/2.0*scale, -float64(satelliteImageHeight)/2.0*scale)

					distance := 38
					position := Position{
						X: planet.Position.X + math.Sqrt2*float64(distance/2),
						Y: planet.Position.Y - math.Sqrt2*float64(distance/2),
					}
					translateToDrawPosition(&screenBounds, position, viewPortCenter, &dio.GeoM, zoomFactor)
					screen.DrawImage(assetLibrary.Images["satellite"], dio)
				}
			}
		}
	}

	{
		imageWidth, imageHeight := assetLibrary.Images["earth"].Size()
		if isInBox(0, 0, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
			dio := &ebiten.DrawImageOptions{}
			scale := 2.0 * zoomFactor
			dio.GeoM.Scale(scale, scale)
			dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)

			translateToDrawPosition(&screenBounds, Position{}, viewPortCenter, &dio.GeoM, zoomFactor)

			screen.DrawImage(assetLibrary.Images["earth"], dio)
		}
	}

	{
		imageWidth, imageHeight := assetLibrary.Images["ship"].Size()
		for _, ship := range w.Ships {
			if isInBox(ship.Position.X, ship.Position.Y, minXToDisplay, maxXToDisplay, minYToDisplay, maxYToDisplay) {
				dio := &ebiten.DrawImageOptions{}
				scale := 1.0 * zoomFactor
				dio.GeoM.Scale(scale, scale)
				dio.GeoM.Translate(-float64(imageWidth)/2.0*scale, -float64(imageHeight)/2.0*scale)
				dio.GeoM.Rotate(-2.0 * math.Pi / 8.0 * float64(ship.Direction))

				translateToDrawPosition(&screenBounds, ship.Position, viewPortCenter, &dio.GeoM, zoomFactor)

				screen.DrawImage(assetLibrary.Images["ship"], dio)
			}
		}
	}

	ui.DrawBoxAround(screen, assetLibrary, 0, 0, 250, 120, ui.Bottom|ui.Right)
	text.Draw(screen, strconv.Itoa(w.score)+"/"+strconv.Itoa(10), fontFace, 4, 26, textColor)
	text.Draw(screen, w.getSelectedShip().Position.String(), fontFace, 4, 54, textColor)
	text.Draw(screen, loseOperationToDoomsdayClockTime(w.lose), fontFace, 4, 110, textColor)

	if w.displayedPlanetName != "" {
		largestPossibleBoundString := text.BoundString(fontFace, "Kepler 99999 jh")
		ui.DrawBoxAround(screen, assetLibrary, (int(screenWidth)-largestPossibleBoundString.Dx())/2, int(screenHeight)-fontFaceHeight, largestPossibleBoundString.Dx(), fontFaceHeight, ui.Left|ui.Top|ui.Right)
		boundString := text.BoundString(fontFace, w.displayedPlanetName)
		text.Draw(screen, w.displayedPlanetName, fontFace, (int(screenWidth)-boundString.Dx())/2, int(screenHeight)-fontFaceHeight+fontShift, textColor)
	}
}

func loseOperationToDoomsdayClockTime(operation *Operation) string {
	numberOfSeconds := 5 * 60
	secondsPerPercent := numberOfSeconds / 100

	minutes := 55 + int(float64(secondsPerPercent)*operation.completedPercentage)/60
	seconds := int(float64(secondsPerPercent)*operation.completedPercentage) % 60

	return "23:" + strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
}
