package ms2k

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
	"github.com/RemiEven/michelSpace2000/src/ms2k/audio"
	"github.com/RemiEven/michelSpace2000/src/ms2k/ui"
)

// CreditScreen displays credits about the game
type CreditScreen struct {
	lines         []string
	maxScroll     int
	currentScroll int

	assetLibrary *assets.Library
}

func NewCreditScreen(assetLibrary *assets.Library) *CreditScreen {
	pseudoTab := "    "
	paragraphs := []string{
		pseudoTab + pseudoTab + "Images",
	}

	addParagraph := func(paragraph string) {
		paragraphs = append(paragraphs, paragraph)
	}

	assetLibrary.ImagesCredits.Range(func(name string, credit assets.Credit) {
		addParagraph("- " + name)
		addParagraph(pseudoTab + "by " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License)
		addParagraph(pseudoTab + credit.Source)
	})

	addParagraph("")
	addParagraph(pseudoTab + pseudoTab + "Fonts")

	assetLibrary.FontFacesCredits.Range(func(name string, credit assets.Credit) {
		addParagraph("- " + name)
		addParagraph(pseudoTab + "by " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License)
		addParagraph(pseudoTab + credit.Source)
	})

	addParagraph("")
	addParagraph(pseudoTab + pseudoTab + "Music and sounds")

	assetLibrary.SoundsCredits.Range(func(name string, credit assets.Credit) {
		addParagraph("- " + name)
		addParagraph(pseudoTab + "by " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License)
		addParagraph(pseudoTab + credit.Source)
	})

	addParagraph("")
	addParagraph(pseudoTab + pseudoTab + "Programming libraries")
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
		addParagraph("- " + name)
		addParagraph(pseudoTab + "by " + strings.Join(credit.Authors, ", ") + ", used under license " + credit.License)
		addParagraph(pseudoTab + credit.Source)
	}

	lines, maxScroll := ui.SplitWallOfText(assetLibrary, 880, 640, paragraphs)

	return &CreditScreen{
		lines:        lines,
		maxScroll:    maxScroll,
		assetLibrary: assetLibrary,
	}
}

// Update updates the credit screen
func (cs *CreditScreen) Update() int8 {
	if cs.repeatingKeyPressed(ebiten.KeyDown) && cs.currentScroll < cs.maxScroll {
		cs.currentScroll++
	}
	if cs.repeatingKeyPressed(ebiten.KeyUp) && cs.currentScroll > 0 {
		cs.currentScroll--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		audio.PlaySound("click")
		return stateInMenu
	}
	return stateInCredits
}

func (cs *CreditScreen) repeatingKeyPressed(key ebiten.Key) bool {
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

func (cs *CreditScreen) Draw(screen *ebiten.Image) {
	drawSpaceBackground(screen, cs.assetLibrary, Position{})

	ui.DrawBoxAround(screen, cs.assetLibrary, 200, 80, 880, 640, ui.AllBorders)

	ui.DrawWallOfText(screen, cs.assetLibrary, 200, 80, cs.lines, cs.currentScroll, len(cs.lines)-cs.maxScroll)
}
