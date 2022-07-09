package assets

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"path/filepath"
	"strings"

	_ "image/png" // needed to correctly load PNG files

	"github.com/hajimehoshi/ebiten/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed files
var assetFS embed.FS

// Library loads and holds all assets of the game
type Library struct {
	Images        map[string]*ebiten.Image
	ImagesCredits map[string]Credit

	Sounds        map[string][]byte
	SoundsCredits map[string]Credit

	FontFaces        map[string]font.Face
	FontFacesCredits map[string]Credit
}

// NewAssetLibrary creates a new asset library with all assets loaded
func NewAssetLibrary() (*Library, error) {
	al := &Library{
		Images:        map[string]*ebiten.Image{},
		ImagesCredits: map[string]Credit{},

		Sounds:        map[string][]byte{},
		SoundsCredits: map[string]Credit{},

		FontFaces:        map[string]font.Face{},
		FontFacesCredits: map[string]Credit{},
	}

	for name, path := range map[string]string{
		"bg":         "back.png",
		"earth":      "Earth.png",
		"moon":       "RedMoon.png",
		"planet":     "Green Gas Planet.png",
		"radar":      "Radar.png",
		"satellite":  "Satellite.png",
		"ships":      "modular_ships.png",
		"ui/listbox": "ui/listbox_default.png",
		"wormHole":   "Hurricane.png",
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
	absolutePath := "files/img/" + path
	content, err := assetFS.ReadFile(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load image [%q]: %w", name, err)
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to decode image [%q]: %w", name, err)
	}
	al.Images[name] = ebiten.NewImageFromImage(img)

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.ImagesCredits[name] = *credit

	return nil
}

func (al *Library) loadSound(path, name string) error {
	absolutePath := "files/audio/" + path
	sound, err := assetFS.ReadFile(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load sound [%q]: %w", name, err)
	}
	al.Sounds[name] = sound

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.SoundsCredits[name] = *credit

	return nil
}

func (al *Library) loadFontFace(path, name string) error {
	absolutePath := "files/font/" + path

	fontFileData, err := assetFS.ReadFile(absolutePath)
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

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.FontFacesCredits[name] = *credit

	return nil
}

func loadCredits(absolutePath string) (*Credit, error) {
	absoluteCreditPath := strings.TrimSuffix(absolutePath, filepath.Ext(absolutePath)) + ".credit.json"
	rawCredits, err := assetFS.ReadFile(absoluteCreditPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	credit := Credit{}
	if err := json.Unmarshal(rawCredits, &credit); err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	return &credit, nil
}
