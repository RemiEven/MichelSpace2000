package ms2k

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
)

const (
	stateInMenu = iota
	stateCreatingGame
	stateInGame
	stateInSettings
	stateWon
	stateLost
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
	assetLibrary *assets.Library

	state int8

	menu             *MainMenu
	gameCreationMenu *GameCreationMenu

	settings *Settings

	World *World
}

// Init initializes a game
func (g *Game) Init() error {
	assetLibrary, err := assets.NewAssetLibrary()
	if err != nil {
		return fmt.Errorf("failed to load asset library: %w", err)
	}
	g.assetLibrary = assetLibrary

	audioContext := NewAudioContext()

	if _, err := NewPlayer(audioContext, g.assetLibrary.Sounds["music"]); err != nil {
		return fmt.Errorf("failed to play music: %w", err)
	}

	g.state = stateInMenu
	g.menu = &MainMenu{}
	g.gameCreationMenu = &GameCreationMenu{}
	g.settings = &Settings{
		keyboardLayout: keyboardLayoutQwerty,
	}

	return nil
}

// Update is used to implement the ebiten.Game interface
func (g *Game) Update() error {
	timeNow := time.Now()

	nextState := g.state
	switch g.state {
	case stateInMenu:
		nextState = g.menu.Update()
		if nextState == stateCreatingGame {
			g.gameCreationMenu.RandomizeSeed()
		}
	case stateCreatingGame:
		nextState = g.gameCreationMenu.Update()
		switch nextState {
		case stateInGame:
			rng, err := rng.NewRNG(string(g.gameCreationMenu.RNG))
			if err != nil {
				log.Fatal(fmt.Errorf("failed to initialize rng: %w", err))
			}
			g.World = NewWorld(rng, timeNow)
		}
	case stateInSettings:
		nextState = g.settings.Update()
	case stateInGame:
		nextState = g.World.Update(timeNow, g.settings)
	case stateLost, stateWon:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			nextState = stateInMenu
			g.menu.selectedIndex = menuStateNewGame
		}
	}
	g.state = nextState

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
	switch g.state {
	case stateInMenu:
		g.menu.Draw(screen, g.assetLibrary)
	case stateInSettings:
		g.settings.Draw(screen, g.assetLibrary)
	case stateCreatingGame:
		g.gameCreationMenu.Draw(screen, g.assetLibrary)
	case stateInGame:
		g.World.Draw(screen, g.assetLibrary)
	case stateLost:
		g.World.Draw(screen, g.assetLibrary)
		fontFace := g.assetLibrary.FontFaces["oxanium"]
		fontFaceHeight := fontFace.Metrics().Height.Ceil()
		fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

		{
			titleLabel := "Game Over"
			boundString := text.BoundString(fontFace, titleLabel)
			ui.DrawBoxAround(screen, g.assetLibrary, (screenWidth-boundString.Dx())/2, fontFaceHeight*11, boundString.Dx(), fontFaceHeight)
			text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*11+fontShift, textColor)
		}
	case stateWon:
		g.World.Draw(screen, g.assetLibrary)
		fontFace := g.assetLibrary.FontFaces["oxanium"]
		fontFaceHeight := fontFace.Metrics().Height.Ceil()
		fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

		{
			titleLabel := "Victory"
			boundString := text.BoundString(fontFace, titleLabel)
			ui.DrawBoxAround(screen, g.assetLibrary, (screenWidth-boundString.Dx())/2, fontFaceHeight*11, boundString.Dx(), fontFaceHeight)
			text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*11+fontShift, textColor)
		}
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
