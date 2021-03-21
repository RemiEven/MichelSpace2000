package ms2k

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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
func (settings *Settings) Draw(screen *ebiten.Image, assetLibrary *AssetLibrary) {
	fontFace := assetLibrary.fontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()

	{
		titleLabel := "MichelSpace2000 - Controls"
		boundString := text.BoundString(fontFace, titleLabel)
		text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*2, textColor)
	}
	{
		newGameLabel := "Key mapping: < " + settings.keyboardLayout + " >"
		boundString := text.BoundString(fontFace, newGameLabel)
		text.Draw(screen, newGameLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*4, textColor)
	}
	{
		previousShipLabel := "Select previous ship: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.PreviousShip)
		boundString := text.BoundString(fontFace, previousShipLabel)
		text.Draw(screen, previousShipLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*6, textColor)
	}
	{
		nextShipLabel := "Select next ship: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.NextShip)
		boundString := text.BoundString(fontFace, nextShipLabel)
		text.Draw(screen, nextShipLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*7, textColor)
	}
	{
		zoomInLabel := "Zoom in: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.ZoomIn)
		boundString := text.BoundString(fontFace, zoomInLabel)
		text.Draw(screen, zoomInLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*8, textColor)
	}
	{
		zoomOutLabel := "Zoom out: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.ZoomOut)
		boundString := text.BoundString(fontFace, zoomOutLabel)
		text.Draw(screen, zoomOutLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*9, textColor)
	}
	{
		upLabel := "Go up: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Up)
		boundString := text.BoundString(fontFace, upLabel)
		text.Draw(screen, upLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*10, textColor)
	}
	{
		downLabel := "Go down: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Down)
		boundString := text.BoundString(fontFace, downLabel)
		text.Draw(screen, downLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*11, textColor)
	}
	{
		leftLabel := "Go left: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Left)
		boundString := text.BoundString(fontFace, leftLabel)
		text.Draw(screen, leftLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*12, textColor)
	}
	{
		rightLabel := "Go right: " + ebitenKeyToString(settings.keyboardLayout, keyMapping.Right)
		boundString := text.BoundString(fontFace, rightLabel)
		text.Draw(screen, rightLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*13, textColor)
	}
}
