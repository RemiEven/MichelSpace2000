package ui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

var boxBgColor = color.RGBA{R: 0x94, G: 0xa9, B: 0xaa, A: 0xff}

func DrawBoxAround(screen *ebiten.Image, assetLibrary *assets.Library, x, y, width, height int) {
	if int(width)%2 == 1 {
		width++
	}
	if int(height)%2 == 1 {
		height++
	}
	scale := 2
	horizontalBorderHeightPx := 6
	intermediaryImage := ebiten.NewImage(int(width/scale)+2*horizontalBorderHeightPx, int(height/scale)+2*horizontalBorderHeightPx)
	ebitenutil.DrawRect(intermediaryImage, float64(horizontalBorderHeightPx), float64(horizontalBorderHeightPx), float64(width/scale), float64(height/scale), boxBgColor)

	baseImage := assetLibrary.Images["ui/listbox"]
	baseImageWidth, baseImageHeight := baseImage.Size()
	topBorder := baseImage.SubImage(image.Rect(40, 2, 41, 2+horizontalBorderHeightPx)).(*ebiten.Image)
	bottomBorder := baseImage.SubImage(image.Rect(40, baseImageHeight-horizontalBorderHeightPx-2, 41, baseImageHeight-2)).(*ebiten.Image)
	topLeftCorner := baseImage.SubImage(image.Rect(1, 2, 1+horizontalBorderHeightPx, 2+horizontalBorderHeightPx)).(*ebiten.Image)
	topRightCorner := baseImage.SubImage(image.Rect(baseImageWidth-horizontalBorderHeightPx-2, 2, baseImageWidth-2, 2+horizontalBorderHeightPx)).(*ebiten.Image)
	bottomLeftCorner := baseImage.SubImage(image.Rect(1, baseImageHeight-2-horizontalBorderHeightPx, 1+horizontalBorderHeightPx, baseImageHeight-2)).(*ebiten.Image)
	bottomRightCorner := baseImage.SubImage(image.Rect(baseImageWidth-horizontalBorderHeightPx-2, baseImageHeight-2-horizontalBorderHeightPx, baseImageWidth-2, baseImageHeight-2)).(*ebiten.Image)

	{
		// Top border
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(horizontalBorderHeightPx), 0)
		for i := 0; i < int(width/scale); i++ {
			intermediaryImage.DrawImage(topBorder, dio)
			dio.GeoM.Translate(1, 0)
		}
	}

	{
		// Bottom border
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(horizontalBorderHeightPx), float64(height/scale+horizontalBorderHeightPx))
		for i := 0; i < int(width/scale); i++ {
			intermediaryImage.DrawImage(bottomBorder, dio)
			dio.GeoM.Translate(1, 0)
		}
	}

	{
		// Left border
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Rotate(-math.Pi / 2)
		dio.GeoM.Translate(0, 1+float64(horizontalBorderHeightPx))
		for i := 0; i < int(height/scale); i++ {
			intermediaryImage.DrawImage(topBorder, dio)
			dio.GeoM.Translate(0, 1)
		}
	}

	{
		// Right border
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Rotate(-math.Pi / 2)
		dio.GeoM.Translate(0, 1)
		dio.GeoM.Translate(float64(width/scale+horizontalBorderHeightPx), float64(horizontalBorderHeightPx))
		for i := 0; i < int(height/scale); i++ {
			intermediaryImage.DrawImage(bottomBorder, dio)
			dio.GeoM.Translate(0, 1)
		}
	}

	{
		// Top left corner
		dio := &ebiten.DrawImageOptions{}
		intermediaryImage.DrawImage(topLeftCorner, dio)
	}

	{
		// Top right corner
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(width/scale+horizontalBorderHeightPx), 0)
		intermediaryImage.DrawImage(topRightCorner, dio)
	}

	{
		// Bottom left corner
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(0, float64(height/scale+horizontalBorderHeightPx))
		intermediaryImage.DrawImage(bottomLeftCorner, dio)
	}

	{
		// Bottom right corner
		dio := &ebiten.DrawImageOptions{}
		dio.GeoM.Translate(float64(width/scale+horizontalBorderHeightPx), float64(height/scale+horizontalBorderHeightPx))
		intermediaryImage.DrawImage(bottomRightCorner, dio)
	}

	dio := &ebiten.DrawImageOptions{}
	dio.GeoM.Scale(float64(scale), float64(scale))
	dio.GeoM.Translate(float64(x-scale*horizontalBorderHeightPx), float64(y-scale*horizontalBorderHeightPx))
	screen.DrawImage(intermediaryImage, dio)
}
