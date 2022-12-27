package assets

import (
	"bytes"
	"context"
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
	"golang.org/x/sync/errgroup"

	"github.com/RemiEven/michelSpace2000/src/ms2k/assets/internal/genericsync"
)

//go:embed files
var assetFS embed.FS

// Library loads and holds all assets of the game
type Library struct {
	Images        genericsync.Map[string, *ebiten.Image]
	ImagesCredits genericsync.Map[string, Credit]

	MP3Sounds     genericsync.Map[string, []byte]
	WavSounds     genericsync.Map[string, []byte]
	SoundsCredits genericsync.Map[string, Credit]

	FontFaces        genericsync.Map[string, font.Face]
	FontFacesCredits genericsync.Map[string, Credit]
}

// NewAssetLibrary creates a new asset library with all assets loaded
func NewAssetLibrary() (<-chan *Library, <-chan error) {
	libraryChan, errChan := make(chan *Library), make(chan error)
	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(7)

	al := &Library{
		Images:        genericsync.Map[string, *ebiten.Image]{},
		ImagesCredits: genericsync.Map[string, Credit]{},

		MP3Sounds:     genericsync.Map[string, []byte]{},
		WavSounds:     genericsync.Map[string, []byte]{},
		SoundsCredits: genericsync.Map[string, Credit]{},

		FontFaces:        genericsync.Map[string, font.Face]{},
		FontFacesCredits: genericsync.Map[string, Credit]{},
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
		path, name := path, name
		eg.Go(func() error {
			return al.loadImage(ctx, path, name)
		})
	}

	eg.Go(func() error {
		return al.loadMP3Sound(ctx, "Hardmoon_-_Deep_space.mp3", "music")
	})

	eg.Go(func() error {
		return al.loadWavSound(ctx, "click.wav", "click")
	})
	eg.Go(func() error {
		return al.loadWavSound(ctx, "click_2.wav", "click_2")
	})

	eg.Go(func() error {
		return al.loadFontFace(ctx, "Oxanium-Regular.ttf", "oxanium")
	})

	go func() {
		if err := eg.Wait(); err != nil {
			errChan <- fmt.Errorf("failed to load an asset: %w", err)
			return
		}

		ships, _ := al.Images.Load("ships")
		al.Images.Store("ship", ships.SubImage(image.Rect(80, 320, 112, 352)).(*ebiten.Image))

		libraryChan <- al
	}()

	return libraryChan, errChan
}

func (al *Library) loadImage(ctx context.Context, path, name string) error {
	absolutePath := "img/" + path
	content, err := al.getFileData(ctx, absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load image [%q]: %w", name, err)
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to decode image [%q]: %w", name, err)
	}
	al.Images.Store(name, ebiten.NewImageFromImage(img))

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.ImagesCredits.Store(name, *credit)

	return nil
}

func (al *Library) loadMP3Sound(ctx context.Context, path, name string) error {
	absolutePath := "audio/" + path
	sound, err := al.getFileData(ctx, absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load mp3 sound [%q]: %w", name, err)
	}
	al.MP3Sounds.Store(name, sound)

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.SoundsCredits.Store(name, *credit)

	return nil
}

func (al *Library) loadWavSound(ctx context.Context, path, name string) error {
	absolutePath := "audio/" + path
	sound, err := al.getFileData(ctx, absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load wav sound [%q]: %w", name, err)
	}
	al.WavSounds.Store(name, sound)

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.SoundsCredits.Store(name, *credit)

	return nil
}

func (al *Library) loadFontFace(ctx context.Context, path, name string) error {
	absolutePath := "font/" + path

	fontFileData, err := al.getFileData(ctx, absolutePath)
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

	al.FontFaces.Store(name, fontFace)

	credit, err := loadCredits(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to load credit file for [%q]: %w", name, err)
	}
	al.FontFacesCredits.Store(name, *credit)

	return nil
}

func loadCredits(absolutePath string) (*Credit, error) {
	absoluteCreditPath := strings.TrimSuffix(absolutePath, filepath.Ext(absolutePath)) + ".credit.json"
	rawCredits, err := assetFS.ReadFile("files/" + absoluteCreditPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	credit := Credit{}
	if err := json.Unmarshal(rawCredits, &credit); err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	return &credit, nil
}
