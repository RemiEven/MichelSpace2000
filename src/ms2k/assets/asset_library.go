package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"

	_ "image/png" // needed to correctly load PNG files

	"github.com/hajimehoshi/ebiten/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed files
var assetFS embed.FS

// Library loads and holds all assets of the game
type Library struct {
	Images    map[string]*ebiten.Image
	Sounds    map[string][]byte
	FontFaces map[string]font.Face
}

// NewAssetLibrary creates a new asset library with all assets loaded
func NewAssetLibrary() (*Library, error) {
	al := &Library{
		Images:    map[string]*ebiten.Image{},
		Sounds:    map[string][]byte{},
		FontFaces: map[string]font.Face{},
	}

	for name, path := range map[string]string{
		"ships":      "modular_ships.png",
		"planet":     "Green Gas Planet.png",
		"bg":         "back.png",
		"earth":      "Earth.png",
		"moon":       "RedMoon.png",
		"wormHole":   "Hurricane.png",
		"satellite":  "Satellite.png",
		"ui/listbox": "ui/listbox_default.png",
	} {
		if err := al.loadImage(path, name); err != nil {
			return nil, err
		}
	}

	al.Images["ship"] = al.Images["ships"].SubImage(image.Rect(80, 320, 112, 352)).(*ebiten.Image)

	if err := al.loadSound("Hardmoon_-_Deep_space.mp3", "music"); err != nil {
		return nil, err
	}

	if err := al.loadFontFace("Oxanium-Regular.ttf", "oxanium"); err != nil {
		return nil, err
	}

	return al, nil
}

func (al *Library) loadImage(path, name string) error {
	content, err := assetFS.ReadFile("files/img/" + path)
	if err != nil {
		return fmt.Errorf("failed to load image [%q]: %w", name, err)
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to decode image [%q]: %w", name, err)
	}
	al.Images[name] = ebiten.NewImageFromImage(img)
	return nil
}

func (al *Library) loadSound(path, name string) error {
	sound, err := assetFS.ReadFile("files/audio/" + path)
	if err != nil {
		return fmt.Errorf("failed to load sound [%q]: %w", name, err)
	}

	al.Sounds[name] = sound
	return nil
}

func (al *Library) loadFontFace(path, name string) error {
	fontFileData, err := assetFS.ReadFile("files/font/" + path)
	if err != nil {
		return fmt.Errorf("failed to read font [%q]: %w", name, err)
	}
	parsedFont, err := opentype.Parse(fontFileData)
	if err != nil {
		return fmt.Errorf("failed to parse font [%q]: %w", name, err)
	}

	const dpi = 72
	fontFace, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("failed to create face from parsed font [%q]: %w", name, err)
	}

	al.FontFaces[name] = fontFace
	return nil
}
