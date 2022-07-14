package ms2k

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
)

const maxSeedLength = 8

type GameCreationMenu struct {
	RNG     []rune
	counter int
}

func (menu *GameCreationMenu) RandomizeSeed() {
	menu.RNG = []rune(rng.RandomSeed())[:maxSeedLength]
}

// Update updates the game creation menu
func (menu *GameCreationMenu) Update() int8 {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return stateInGame
	}
	inputChars := ebiten.InputChars()
	charactersToAdd := make([]rune, 0, len(inputChars))
	for _, r := range inputChars {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			charactersToAdd = append(charactersToAdd, r)
		}
	}

	menu.RNG = append(menu.RNG, charactersToAdd...)
	if len(menu.RNG) > maxSeedLength {
		menu.RNG = menu.RNG[:maxSeedLength]
	}

	if menu.repeatingKeyPressed(ebiten.KeyBackspace) && len(menu.RNG) > 0 {
		menu.RNG = menu.RNG[:len(menu.RNG)-1]
	}

	menu.counter++
	menu.counter %= 60
	return stateCreatingGame
}

func (menu *GameCreationMenu) repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

// Draw draws the game creation menu
func (menu *GameCreationMenu) Draw(screen *ebiten.Image, assetLibrary *assets.Library) {
	drawSpaceBackground(screen, assetLibrary, Position{})

	screenWidth := screen.Bounds().Dx()

	fontFace := assetLibrary.FontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()
	fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

	baseRNGSeedLabel := "RNG seed: "
	largestBoundString := text.BoundString(fontFace, baseRNGSeedLabel+strings.Repeat("w", maxSeedLength))

	{
		titleLabel := "Game creation"
		boundString := text.BoundString(fontFace, titleLabel)
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*5, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*5+fontShift, textColor)
	}

	{
		baseRNGSeedLabel := "RNG seed: "
		rngSeedLabel := baseRNGSeedLabel + string(menu.RNG)
		if menu.counter < 30 && len(menu.RNG) < maxSeedLength {
			rngSeedLabel += "_"
		}
		boundString := text.BoundString(fontFace, baseRNGSeedLabel+strings.Repeat("w", maxSeedLength))
		ui.DrawBoxAround(screen, assetLibrary, (screenWidth-largestBoundString.Dx())/2, fontFaceHeight*9, largestBoundString.Dx(), fontFaceHeight, ui.AllBorders)
		text.Draw(screen, rngSeedLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*9+fontShift, textColor)
	}
}
