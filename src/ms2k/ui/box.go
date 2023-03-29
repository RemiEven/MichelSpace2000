package ui

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

type BorderOption uint8

const (
	Left       = BorderOption(1 << iota)
	Right      = BorderOption(1 << iota)
	Top        = BorderOption(1 << iota)
	Bottom     = BorderOption(1 << iota)
	AllBorders = Left | Right | Top | Bottom
	NoBorder   = BorderOption(0)
)

func DrawBoxAround(screen *ebiten.Image, assetLibrary *assets.Library, x, y, width, height int, borderOption BorderOption) {
	if int(width)%2 == 1 {
		width++
	}
	if int(height)%2 == 1 {
		height++
	}
	scale := 2
	horizontalBorderHeightPx := 6

	var (
		drawLeftBorder   = (borderOption | Left) == borderOption
		drawRightBorder  = (borderOption | Right) == borderOption
		drawTopBorder    = (borderOption | Top) == borderOption
		drawBottomBorder = (borderOption | Bottom) == borderOption
	)

	intermediaryImage := ebiten.NewImage(
		int(width/scale)+numberOfTrues(drawLeftBorder, drawRightBorder)*horizontalBorderHeightPx,
		int(height/scale)+numberOfTrues(drawTopBorder, drawBottomBorder)*horizontalBorderHeightPx,
	)
	vector.DrawFilledRect(intermediaryImage, float32(numberOfTrues(drawLeftBorder)*horizontalBorderHeightPx), float32(numberOfTrues(drawTopBorder)*horizontalBorderHeightPx), float32(width/scale), float32(height/scale), BoxBgColor, false)

	baseImage, _ := assetLibrary.Images.Load("ui/listbox")
	baseImageWidth, baseImageHeight := baseImage.Bounds().Dx(), baseImage.Bounds().Dy()

	var (
		topBorder         = baseImage.SubImage(image.Rect(40, 2, 41, 2+horizontalBorderHeightPx)).(*ebiten.Image)
		bottomBorder      = baseImage.SubImage(image.Rect(40, baseImageHeight-horizontalBorderHeightPx-2, 41, baseImageHeight-2)).(*ebiten.Image)
		topLeftCorner     = baseImage.SubImage(image.Rect(1, 2, 1+horizontalBorderHeightPx, 2+horizontalBorderHeightPx)).(*ebiten.Image)
		topRightCorner    = baseImage.SubImage(image.Rect(baseImageWidth-horizontalBorderHeightPx-2, 2, baseImageWidth-2, 2+horizontalBorderHeightPx)).(*ebiten.Image)
		bottomLeftCorner  = baseImage.SubImage(image.Rect(1, baseImageHeight-2-horizontalBorderHeightPx, 1+horizontalBorderHeightPx, baseImageHeight-2)).(*ebiten.Image)
		bottomRightCorner = baseImage.SubImage(image.Rect(baseImageWidth-horizontalBorderHeightPx-2, baseImageHeight-2-horizontalBorderHeightPx, baseImageWidth-2, baseImageHeight-2)).(*ebiten.Image)
	)

	if drawTopBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(numberOfTrues(drawLeftBorder)*horizontalBorderHeightPx), 0)
		for i := 0; i < int(width/scale); i++ {
			intermediaryImage.DrawImage(topBorder, dio)
			dio.GeoM.Translate(1, 0)
		}
	}

	if drawBottomBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(numberOfTrues(drawLeftBorder)*horizontalBorderHeightPx), float64(height/scale+numberOfTrues(drawTopBorder)*horizontalBorderHeightPx))
		for i := 0; i < int(width/scale); i++ {
			intermediaryImage.DrawImage(bottomBorder, dio)
			dio.GeoM.Translate(1, 0)
		}
	}

	if drawLeftBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Rotate(-math.Pi / 2)
		dio.GeoM.Translate(0, 1+float64(numberOfTrues(drawTopBorder)*horizontalBorderHeightPx))
		for i := 0; i < int(height/scale); i++ {
			intermediaryImage.DrawImage(topBorder, dio)
			dio.GeoM.Translate(0, 1)
		}
	}

	if drawRightBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Rotate(-math.Pi / 2)
		dio.GeoM.Translate(0, 1)
		dio.GeoM.Translate(float64(width/scale+numberOfTrues(drawLeftBorder)*horizontalBorderHeightPx), float64(numberOfTrues(drawTopBorder)*horizontalBorderHeightPx))
		for i := 0; i < int(height/scale); i++ {
			intermediaryImage.DrawImage(bottomBorder, dio)
			dio.GeoM.Translate(0, 1)
		}
	}

	if drawTopBorder && drawLeftBorder {
		dio := &ebiten.DrawImageOptions{}
		intermediaryImage.DrawImage(topLeftCorner, dio)
	}

	if drawTopBorder && drawRightBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(width/scale+numberOfTrues(drawLeftBorder)*horizontalBorderHeightPx), 0)
		intermediaryImage.DrawImage(topRightCorner, dio)
	}

	if drawBottomBorder && drawLeftBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(0, float64(height/scale+numberOfTrues(drawTopBorder)*horizontalBorderHeightPx))
		intermediaryImage.DrawImage(bottomLeftCorner, dio)
	}

	if drawBottomBorder && drawRightBorder {
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(width/scale+numberOfTrues(drawLeftBorder)*horizontalBorderHeightPx), float64(height/scale+numberOfTrues(drawTopBorder)*horizontalBorderHeightPx))
		intermediaryImage.DrawImage(bottomRightCorner, dio)
	}

	dio := &ebiten.DrawImageOptions{}
	dio.GeoM.Scale(float64(scale), float64(scale))
	dio.GeoM.Translate(float64(x-numberOfTrues(drawLeftBorder)*scale*horizontalBorderHeightPx), float64(y-numberOfTrues(drawTopBorder)*scale*horizontalBorderHeightPx))
	screen.DrawImage(intermediaryImage, dio)
}

func numberOfTrues(conditions ...bool) int {
	number := 0
	for _, b := range conditions {
		if b {
			number++
		}
	}
	return number
}
