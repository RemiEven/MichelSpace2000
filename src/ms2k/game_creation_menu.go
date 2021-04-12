package ms2k

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/rng"
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

	if menu.repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(menu.RNG) > 0 {
			menu.RNG = menu.RNG[:len(menu.RNG)-1]
		}
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
func (menu *GameCreationMenu) Draw(screen *ebiten.Image, assetLibrary *AssetLibrary) {
	fontFace := assetLibrary.fontFaces["oxanium"]
	fontFaceHeight := fontFace.Metrics().Height.Ceil()

	{
		titleLabel := "Game creation"
		boundString := text.BoundString(fontFace, titleLabel)
		text.Draw(screen, titleLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*2, textColor)
	}

	{
		baseRNGSeedLabel := "RNG seed: "
		rngSeedLabel := baseRNGSeedLabel + string(menu.RNG)
		if menu.counter < 30 && len(menu.RNG) < maxSeedLength {
			rngSeedLabel += "_"
		}
		boundString := text.BoundString(fontFace, baseRNGSeedLabel+strings.Repeat("9", maxSeedLength))
		text.Draw(screen, rngSeedLabel, fontFace, (screenWidth-boundString.Dx())/2, fontFaceHeight*4, textColor)
	}
}
