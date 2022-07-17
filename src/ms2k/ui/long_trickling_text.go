package ui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

type LongTricklingText struct {
	texts              []string
	displayedTextIndex int
	frequency          time.Duration

	tt *TricklingText
}

func NewLongTricklingText(texts []string, timeNow time.Time, frequency time.Duration) *LongTricklingText {
	return &LongTricklingText{
		texts:     texts,
		tt:        NewTricklingText(texts[0], timeNow, frequency),
		frequency: frequency,
	}
}

func (ltt *LongTricklingText) Update(timeNow time.Time) (addedRune, allShown bool) {
	addedRune, stepAllShown := ltt.tt.Update(timeNow)
	if !stepAllShown {
		return addedRune, false
	}
	if ltt.displayedTextIndex == len(ltt.texts)-1 {
		return addedRune, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ltt.displayedTextIndex++
		ltt.tt = NewTricklingText(ltt.texts[ltt.displayedTextIndex], timeNow, ltt.frequency)
		return true, false
	}
	return false, true
}

func (ltt *LongTricklingText) Draw(screen *ebiten.Image, assetLibrary *assets.Library, x, y, width, height int) {
	ltt.tt.Draw(screen, assetLibrary, x, y, width, height)
}
