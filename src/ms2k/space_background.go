package ms2k

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets"
)

func drawSpaceBackground(screen *ebiten.Image, assetLibrary *assets.Library, position Position) {
	screenWidth, screenHeight := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	scale := 1.0

	parallaxFactor := math.Pow(3.0, zoomFactor)
	imageWidth, imageHeight := assetLibrary.Images["bg"].Size()
	topLeftBackgroundTileX := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*position.X - screenWidth/2 /*/zoomFactor*/) / (float64(imageWidth) * scale)))
	topLeftBackgroundTileY := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*position.Y - screenHeight/2 /*/zoomFactor*/) / (float64(imageHeight) * scale)))
	bottomRightBackgroundTileX := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*position.X + screenWidth/2 /*/zoomFactor*/) / (float64(imageWidth) * scale)))
	bottomRightBackgroundTileY := int(math.Floor((((parallaxFactor-1.0)/parallaxFactor)*position.Y + screenHeight/2 /*/zoomFactor*/) / (float64(imageHeight) * scale)))
	x := topLeftBackgroundTileX
	for x <= bottomRightBackgroundTileX {
		y := topLeftBackgroundTileY
		for y <= bottomRightBackgroundTileY {
			dio := &ebiten.DrawImageOptions{}
			dio.GeoM.Translate(float64(x)*scale*float64(imageWidth)+screenWidth/2.0 /*/zoomFactor*/ -(parallaxFactor-1.0)*position.X/parallaxFactor, float64(y)*scale*float64(imageHeight)+screenHeight/2.0 /*/zoomFactor*/ -(parallaxFactor-1.0)*position.Y/parallaxFactor)
			screen.DrawImage(assetLibrary.Images["bg"], dio)
			y++
		}
		x++
	}
}
