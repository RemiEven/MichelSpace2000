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

	assetLibrary *assets.Library

	tt *TricklingText
}

func NewLongTricklingText(texts []string, timeNow time.Time, frequency time.Duration, assetLibrary *assets.Library) *LongTricklingText {
	return &LongTricklingText{
		texts:        texts,
		tt:           NewTricklingText(texts[0], timeNow, frequency, assetLibrary),
		frequency:    frequency,
		assetLibrary: assetLibrary,
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
		ltt.tt = NewTricklingText(ltt.texts[ltt.displayedTextIndex], timeNow, ltt.frequency, ltt.assetLibrary)
		return true, false
	}
	return false, true
}

func (ltt *LongTricklingText) Draw(screen *ebiten.Image, x, y, width, height int) {
	ltt.tt.Draw(screen, x, y, width, height)
}
