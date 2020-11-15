package ms2k

import (
	"fmt"
	"image"
	"io/ioutil"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// AssetLibrary loads and holds all assets of the game
type AssetLibrary struct {
	images map[string]*ebiten.Image
	sounds map[string][]byte
}

// NewAssetLibrary creates a new asset library with all assets loaded
func NewAssetLibrary() (*AssetLibrary, error) {
	al := &AssetLibrary{
		images: map[string]*ebiten.Image{},
		sounds: map[string][]byte{},
	}

	for name, path := range map[string]string{
		"ships":  "modular_ships.png",
		"planet": "Green Gas Planet.png",
		"bg":     "back.png",
	} {
		if err := al.loadImage(path, name); err != nil {
			return nil, err
		}
	}

	al.images["ship"] = al.images["ships"].SubImage(image.Rect(80, 320, 112, 352)).(*ebiten.Image)

	al.loadSound("Hardmoon_-_Deep_space.mp3", "music")

	return al, nil
}

func (al *AssetLibrary) loadImage(path, name string) error {
	img, _, err := ebitenutil.NewImageFromFile("./assets/img/" + path)
	if err != nil {
		return fmt.Errorf("failed to load image [%q]: %w", name, err)
	}
	al.images[name] = img
	return nil
}

func (al *AssetLibrary) loadSound(path, name string) error {
	sound, err := ioutil.ReadFile("./assets/audio/" + path)
	if err != nil {
		return fmt.Errorf("failed to load sound [%q]: %w", name, err)
	}

	al.sounds[name] = sound
	return nil
}
