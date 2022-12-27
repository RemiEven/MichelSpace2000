package ui

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

const Lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func SplitWallOfText(assetLibrary *assets.Library, width, height int, paragraphs []string) ([]string, int) {
	fontFace, _ := assetLibrary.FontFaces.Load("oxanium")

	lines := []string{}
	for _, paragraph := range paragraphs {
		words := strings.Split(paragraph, " ")
		currentLine := ""
		for i := range words {
			boundString := text.BoundString(fontFace, currentLine+words[i]+" ") // TODO: handle linefeeds
			if currentLine != "" && boundString.Dx() > width {
				lines = append(lines, currentLine)
				currentLine = ""
			}

			currentLine += words[i] + " "
		}
		lines = append(lines, currentLine)
	}

	numberOfDrawableLines := height / fontFace.Metrics().Height.Ceil()
	return lines, max(0, len(lines)-numberOfDrawableLines)
}

func DrawWallOfText(screen *ebiten.Image, assetLibrary *assets.Library, x, y int, lines []string, lineOffset, lineNumber int) {
	fontFace, _ := assetLibrary.FontFaces.Load("oxanium")
	fontShift := (fontFace.Metrics().Ascent + (fontFace.Metrics().Height-fontFace.Metrics().Ascent-fontFace.Metrics().Descent)/2).Ceil()

	numberOfLinesToDraw := min(len(lines)-lineOffset, lineNumber)

	text.Draw(screen, strings.Join(lines[lineOffset:lineOffset+numberOfLinesToDraw], "\n"), fontFace, x, y+fontShift, TextColor)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
