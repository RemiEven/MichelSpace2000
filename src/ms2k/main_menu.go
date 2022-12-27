package ms2k

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/audio"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
)

const (
	menuStateNewGame = iota
	menuStateSettings
	menuStateCredits
	menuStateExit
)

// MainMenu is the main menu of the game
type MainMenu struct {
	states        []int8
	selectedIndex int

	assetLibrary *assets.Library
}

func NewMainMenu(assetLibrary *assets.Library, allowExit bool) *MainMenu {
	states := []int8{menuStateNewGame, menuStateSettings, menuStateCredits, menuStateExit}
	if !allowExit {
		states = states[0:3]
	}
	return &MainMenu{
		states:       states,
		assetLibrary: assetLibrary,
	}
}

func (menu *MainMenu) state() int8 {
	index := menu.selectedIndex % len(menu.states)
	if index < 0 {
		index += len(menu.states)
	}
	return menu.states[index]
}

func (menu *MainMenu) Reset() {
	menu.selectedIndex = menuStateNewGame
}

// Update updates the MainMenu
func (menu *MainMenu) Update() int8 {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		audio.PlaySound("click")
		switch menu.state() {
		case menuStateNewGame:
			return stateCreatingGame
		case menuStateSettings:
			return stateInSettings
		case menuStateCredits:
			return stateInCredits
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
func (menu *MainMenu) Draw(screen *ebiten.Image) {
	drawSpaceBackground(screen, menu.assetLibrary, Position{})

	screenWidth := screen.Bounds().Dx()

	fontFace, _ := menu.assetLibrary.FontFaces.Load("oxanium")
	fontFaceHeight := fontFace.Metrics().Height.Ceil()
	fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

	longestLabel := "MichelSpace2000"
	largestBoundString := text.BoundString(fontFace, longestLabel)

	{
		titleLabel := "MichelSpace2000"
		boundString := text.BoundString(fontFace, titleLabel)
		ui.DrawBoxAround(screen, menu.assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*5, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*5+fontShift, ui.TextColor)
	}

	color := func(menuOption int8) color.Color {
		if menuOption == menu.state() {
			return ui.SelectedTextColor
		}
		return ui.TextColor
	}

	{
		newGameLabel := "New game"
		boundString := text.BoundString(fontFace, newGameLabel)
		ui.DrawBoxAround(screen, menu.assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*9, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, newGameLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*9+fontShift, color(menuStateNewGame))
	}
	{
		settingsLabel := "Controls"
		boundString := text.BoundString(fontFace, settingsLabel)
		ui.DrawBoxAround(screen, menu.assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*11, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, settingsLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*11+fontShift, color(menuStateSettings))
	}
	{
		creditsLabel := "Credits"
		boundString := text.BoundString(fontFace, creditsLabel)
		ui.DrawBoxAround(screen, menu.assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*13, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, creditsLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*13+fontShift, color(menuStateCredits))
	}
	if len(menu.states) == 4 {
		exitLabel := "Exit"
		boundString := text.BoundString(fontFace, exitLabel)
		ui.DrawBoxAround(screen, menu.assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*15, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, exitLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*15+fontShift, color(menuStateExit))
	}
}
