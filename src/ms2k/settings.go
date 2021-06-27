package ms2k

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
)

// Settings holds the settings of the game
type Settings struct {
	keyboardLayout          string
	selectedKeyMappingIndex int
}

// Update updates the settings
func (settings *Settings) Update() int8 {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return stateInMenu
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		settings.selectedKeyMappingIndex++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		settings.selectedKeyMappingIndex--
	}
	settings.selectedKeyMappingIndex %= len(keyboardLayouts)
	if settings.selectedKeyMappingIndex < 0 {
		settings.selectedKeyMappingIndex += len(keyboardLayouts)
	}
	settings.keyboardLayout = keyboardLayouts[settings.selectedKeyMappingIndex]

	return stateInSettings
}

// Draw draws the settings
func (settings *Settings) Draw(screen *ebiten.Image, assetLibrary *assets.Library) {
	drawSpaceBackground(screen, assetLibrary, Position{})

	fontFace := assetLibrary.FontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()
	fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

	longestLabel := "MichelSpace2000 - Controls"
	largestBoundString := text.BoundString(fontFace, longestLabel)

	{
		titleLabel := "MichelSpace2000 - Controls"
		boundString := text.BoundString(fontFace, titleLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*5, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*5+fontShift, textColor)
	}

	{
		keyMappingLabel := "Key mapping: < " + settings.keyboardLayout + " >"
		boundString := text.BoundString(fontFace, keyMappingLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*9, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, keyMappingLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*9+fontShift, textColor)
	}

	ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*11, largestBoundString.Dx(), fontFaceHeight*8, ui.AllBorders)

	{
		previousShipLabel := "Select previous ship: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.PreviousShip)
		boundString := text.BoundString(fontFace, previousShipLabel)
		text.Draw(screen, previousShipLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*11+fontShift, textColor)
	}
	{
		nextShipLabel := "Select next ship: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.NextShip)
		boundString := text.BoundString(fontFace, nextShipLabel)
		text.Draw(screen, nextShipLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*12+fontShift, textColor)
	}
	{
		zoomInLabel := "Zoom in: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.ZoomIn)
		boundString := text.BoundString(fontFace, zoomInLabel)
		text.Draw(screen, zoomInLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*13+fontShift, textColor)
	}
	{
		zoomOutLabel := "Zoom out: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.ZoomOut)
		boundString := text.BoundString(fontFace, zoomOutLabel)
		text.Draw(screen, zoomOutLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*14+fontShift, textColor)
	}
	{
		upLabel := "Go up: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Up)
		boundString := text.BoundString(fontFace, upLabel)
		text.Draw(screen, upLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*15+fontShift, textColor)
	}
	{
		downLabel := "Go down: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Down)
		boundString := text.BoundString(fontFace, downLabel)
		text.Draw(screen, downLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*16+fontShift, textColor)
	}
	{
		leftLabel := "Go left: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Left)
		boundString := text.BoundString(fontFace, leftLabel)
		text.Draw(screen, leftLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*17+fontShift, textColor)
	}
	{
		rightLabel := "Go right: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Right)
		boundString := text.BoundString(fontFace, rightLabel)
		text.Draw(screen, rightLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*18+fontShift, textColor)
	}
}
