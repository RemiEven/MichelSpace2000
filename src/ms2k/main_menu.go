package ms2k

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
)

const (
	menuStateNewGame = iota
	menuStateSettings
	menuStateExit
)

var (
	menuStates    = []int8{menuStateNewGame, menuStateSettings, menuStateExit}
	lenMenuStates = len(menuStates)
)

// MainMenu is the main menu of the game
type MainMenu struct {
	selectedIndex int
}

func (menu *MainMenu) state() int8 {
	index := menu.selectedIndex % lenMenuStates
	if index < 0 {
		index += lenMenuStates
	}
	return menuStates[index]
}

// Update updates the MainMenu
func (menu *MainMenu) Update() int8 {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch menu.state() {
		case menuStateNewGame:
			return stateCreatingGame
		case menuStateSettings:
			return stateInSettings
		case menuStateExit:
			os.Exit(0)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		menu.selectedIndex++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		menu.selectedIndex--
	}
	return stateInMenu
}

// Draw draws the MainMenu
func (menu *MainMenu) Draw(screen *ebiten.Image, assetLibrary *assets.Library) {
	drawSpaceBackground(screen, assetLibrary, Position{})

	fontFace := assetLibrary.FontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()
	fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

	longestLabel := "MichelSpace2000"
	largestBoundString := text.BoundString(fontFace, longestLabel)

	{
		titleLabel := "MichelSpace2000"
		boundString := text.BoundString(fontFace, titleLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*5, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*5+fontShift, textColor)
	}

	color := func(menuOption int8) color.Color {
		if menuOption == menu.state() {
			return selectedTextColor
		}
		return textColor
	}

	{
		newGameLabel := "New game"
		boundString := text.BoundString(fontFace, newGameLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*9, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, newGameLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*9+fontShift, color(menuStateNewGame))
	}
	{
		settingsLabel := "Controls"
		boundString := text.BoundString(fontFace, settingsLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*11, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, settingsLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*11+fontShift, color(menuStateSettings))
	}
	{
		exitLabel := "Exit"
		boundString := text.BoundString(fontFace, exitLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*13, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, exitLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*13+fontShift, color(menuStateExit))
	}
}
