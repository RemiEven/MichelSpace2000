package ms2k

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/audio"
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
	stateInCredits
	stateLoadingAssets
	stateLoadingAssetsError
)

const (
	viewportBorderMargin = 32 // should be equal or bigger than half the side length of the biggest sprite to avoid clipping
)

var (
	zoomFactor = 1.0
)

// Game contains all loaded game assets with current game data
type Game struct {
	assetLibraryReadyChan <-chan *assets.Library
	assetLibraryErrChan   <-chan error
	assetLibrary          *assets.Library
	loadingAssetErr       error

	state int8

	menu             *MainMenu
	gameCreationMenu *GameCreationMenu

	settings *Settings

	World *World

	creditScreen *CreditScreen
}

// Init initializes a game
func (g *Game) Init() error {
	g.assetLibraryReadyChan, g.assetLibraryErrChan = assets.NewAssetLibrary()
	g.state = stateLoadingAssets
	return nil
}

// Update is used to implement the ebiten.Game interface
func (g *Game) Update() error {
	timeNow := time.Now()

	nextState := g.state
	switch g.state {
	case stateLoadingAssets:
		select {
		case err := <-g.assetLibraryErrChan:
			g.loadingAssetErr = fmt.Errorf("failed to load asset library: %w", err)
			nextState = stateLoadingAssetsError
		case al := <-g.assetLibraryReadyChan:
			g.assetLibrary = al

			audio.Init(g.assetLibrary)
			audio.PlaySound("music")

			g.menu = NewMainMenu(g.assetLibrary, allowExit)
			g.gameCreationMenu = NewGameCreationMenu(g.assetLibrary)
			g.settings = NewSettings(g.assetLibrary)
			g.creditScreen = NewCreditScreen(g.assetLibrary)
			nextState = stateInMenu
		default:
			// do nothing and check again next frame
		}
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
			g.World = NewWorld(rng, timeNow, g.assetLibrary)
		}
	case stateInSettings:
		nextState = g.settings.Update()
	case stateInGame:
		nextState = g.World.Update(timeNow, g.settings)
	case stateLost, stateWon:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			nextState = stateInMenu
			g.menu.Reset()
		}
	case stateInCredits:
		nextState = g.creditScreen.Update()
	}
	g.state = nextState

	return nil
}

// Draw is used to implement the ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	screenBounds := screen.Bounds()

	switch g.state {
	case stateLoadingAssets:
		ebitenutil.DebugPrint(screen, "Loading assets, please wait...")
	case stateLoadingAssetsError:
		ebitenutil.DebugPrint(screen, g.loadingAssetErr.Error())
	case stateInMenu:
		g.menu.Draw(screen)
	case stateInSettings:
		g.settings.Draw(screen)
	case stateCreatingGame:
		g.gameCreationMenu.Draw(screen)
	case stateInGame:
		g.World.Draw(screen)
	case stateLost:
		g.World.Draw(screen)
		fontFace, _ := g.assetLibrary.FontFaces.Load("oxanium")
		fontFaceHeight := fontFace.Metrics().Height.Ceil()
		fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

		{
			titleLabel := "Game Over"
			boundString := text.BoundString(fontFace, titleLabel)
			ui.DrawBoxAround(screen, g.assetLibrary, (screenBounds.Dx()-boundString.Dx())/2, fontFaceHeight*11, boundString.Dx(), fontFaceHeight, ui.AllBorders)
			text.Draw(screen, titleLabel, fontFace, (screenBounds.Dx()-boundString.Dx())/2, fontFaceHeight*11+fontShift, ui.TextColor)
		}
	case stateWon:
		g.World.Draw(screen)
		fontFace, _ := g.assetLibrary.FontFaces.Load("oxanium")
		fontFaceHeight := fontFace.Metrics().Height.Ceil()
		fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

		{
			titleLabel := "Victory"
			boundString := text.BoundString(fontFace, titleLabel)
			ui.DrawBoxAround(screen, g.assetLibrary, (screenBounds.Dx()-boundString.Dx())/2, fontFaceHeight*11, boundString.Dx(), fontFaceHeight, ui.AllBorders)
			text.Draw(screen, titleLabel, fontFace, (screenBounds.Dx()-boundString.Dx())/2, fontFaceHeight*11+fontShift, ui.TextColor)
		}
	case stateInCredits:
		g.creditScreen.Draw(screen)
	}
}

func translateToDrawPosition(screenBounds *image.Rectangle, gamePosition, viewPortCenter Position, geoM *ebiten.GeoM, zoomFactor float64) {
	screenWidth, screenHeight := float64(screenBounds.Dx()), float64(screenBounds.Dy())
	geoM.Translate(-viewPortCenter.X*zoomFactor, -viewPortCenter.Y*zoomFactor)
	geoM.Translate(screenWidth/2, screenHeight/2)
	geoM.Translate(gamePosition.X*zoomFactor, gamePosition.Y*zoomFactor)
}

// Layout is used to implement the ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func isInBox(x, y, minX, maxX, minY, maxY float64) bool {
	return minX <= x && x <= maxX && minY <= y && y <= maxY
}
