package ui

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

type TricklingText struct {
	text string

	lines         []string
	trickledLines []string
	width, height int

	start     time.Time
	frequency time.Duration

	numberOfRunesToShow int
	showAll             bool

	addedRunes bool
}

func NewTricklingText(text string, timeNow time.Time, frequency time.Duration) *TricklingText {
	return &TricklingText{
		text:      text,
		start:     timeNow,
		frequency: frequency,
	}
}

func (tt *TricklingText) Update(timeNow time.Time) (addedRune, allShown bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		tt.showAll = true
	}
	if tt.showAll {
		tt.numberOfRunesToShow = math.MaxInt
	} else {
		tt.numberOfRunesToShow = int(timeNow.Sub(tt.start) / tt.frequency)
	}

	if tt.addedRunes {
		tt.addedRunes = false
		return true, false
	}
	return false, tt.allShown()
}

func (tt *TricklingText) Draw(screen *ebiten.Image, assetLibrary *assets.Library, x, y, width, height int) {
	if tt.lines == nil || tt.width != width || tt.height != height {
		tt.lines, _ = SplitWallOfText(assetLibrary, width, height, []string{tt.text})
		tt.trickledLines = make([]string, len(tt.lines))
		tt.width, tt.height = width, height
	}

	totalNumberOfRunes := numberOfRunesIn(tt.lines)
	numberOfRunesToShow := min(tt.numberOfRunesToShow, totalNumberOfRunes)
	if numberOfRunesToShow == totalNumberOfRunes {
		tt.trickledLines = tt.lines
	} else {
		numberOfShownRunes := numberOfRunesIn(tt.trickledLines)
		numberOfRunesToAdd := numberOfRunesToShow - numberOfShownRunes

		currentLine := max(0, firstEmptyLineIndexIn(tt.trickledLines)-1)
		for numberOfRunesToAdd > 0 {
			tt.addedRunes = true
			numberOfShownRunesInLine := runeLen(tt.trickledLines[currentLine])
			numberOfRunesToAddToLine := min(numberOfRunesToAdd, runeLen(tt.lines[currentLine])-numberOfShownRunesInLine)
			tt.trickledLines[currentLine] += tt.lines[currentLine][numberOfShownRunesInLine : numberOfShownRunesInLine+numberOfRunesToAddToLine]
			currentLine++
			numberOfRunesToAdd -= numberOfRunesToAddToLine
		}
	}

	DrawBoxAround(screen, assetLibrary, x, y, width, height, AllBorders)
	DrawWallOfText(screen, assetLibrary, x, y, tt.trickledLines, 0, len(tt.trickledLines))
}

func (tt *TricklingText) allShown() bool {
	linesLength := len(tt.lines)
	switch {
	case linesLength != len(tt.trickledLines):
		return false
	case linesLength == 0:
		return true
	default:
		return tt.lines[linesLength-1] == tt.trickledLines[linesLength-1]
	}
}

func numberOfRunesIn(lines []string) int {
	result := 0
	for _, line := range lines {
		result += runeLen(line)
	}
	return result
}

func runeLen(s string) int {
	return len([]rune(s))
}

func firstEmptyLineIndexIn(lines []string) int {
	for i := range lines {
		if lines[i] == "" {
			return i
		}
	}
	return len(lines)
}
