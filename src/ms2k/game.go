package ms2k

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	World *World
	score int
	won   bool
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

	return nil
}

// Update is used to implement the ebiten.Game interface
func (g *Game) Update() error {
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
				planet.Looted = true
				g.score++
			}
		}
	}

	if g.score == 200 {
		g.won = true
	}

	return nil
}

// Draw is used to implement the ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	if g.won {
		ebitenutil.DebugPrint(screen, "victory")
		return
	}

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
				// screen.DrawImage(rrr, dio)
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

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.score)+"/"+strconv.Itoa(len(g.World.Planets)))
	ebitenutil.DebugPrintAt(screen, g.World.getSelectedShip().Position.String(), 0, 16)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(int(minXToDisplay)), 0, 48)
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
