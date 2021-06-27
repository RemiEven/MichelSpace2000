package ms2k

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

// CreditScreen displays credits about the game
type CreditScreen struct {
	text string
}

func NewCreditScreen(assetLibrary *assets.Library) *CreditScreen {
	text := "Images\n"

	for name, credit := range assetLibrary.ImagesCredits {
		text += "\n" + name
		text += "\nby " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License
		text += "\n- " + credit.Source
	}

	text += "\n\nFonts\n"

	for name, credit := range assetLibrary.FontFacesCredits {
		text += "\n" + name
		text += "\nby " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License
		text += "\n- " + credit.Source
	}

	text += "\n\nMusic and sounds\n"

	for name, credit := range assetLibrary.SoundsCredits {
		text += "\n" + name
		text += "\nby " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License
		text += "\n- " + credit.Source
	}

	text += "\n\nProgramming libraries\n"
	for name, credit := range map[string]assets.Credit{
		"ebiten": {
			Authors: []string{"hajimehoshi"},
			License: "Apache-2.0",
			Source:  "https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2",
		},
		"opensimplex-go": {
			Authors: []string{"ojrac"},
			License: "Unlicense",
			Source:  "https://pkg.go.dev/github.com/ojrac/opensimplex-go",
		},
	} {
		text += "\n" + name
		text += "\nby " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License
		text += "\n- " + credit.Source
	}

	return &CreditScreen{
		text: text,
	}
}

// Update updates the credit screen
func (cs *CreditScreen) Update() int8 {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return stateInMenu
	}
	return stateInCredits
}

func (cs *CreditScreen) Draw(screen *ebiten.Image, assetLibrary *assets.Library) {
	drawSpaceBackground(screen, assetLibrary, Position{})

	ebitenutil.DebugPrint(screen, cs.text)
}
