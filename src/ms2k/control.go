package ms2k

import "github.com/hajimehoshi/ebiten/v2"

const (
	keyboardLayoutQwerty = "QWERTY"
	keyboardLayoutAzerty = "AZERTY"
)

var keyboardLayouts = []string{keyboardLayoutAzerty, keyboardLayoutQwerty}

// KeyMapping holds a mapping
type KeyMapping struct {
	ID string

	PreviousShip, NextShip ebiten.Key

	ZoomIn, ZoomOut ebiten.Key

	Up, Down, Left, Right ebiten.Key
}

var keyMapping = &KeyMapping{
	PreviousShip: ebiten.KeyA,
	NextShip:     ebiten.KeyD,

	ZoomIn:  ebiten.KeyW,
	ZoomOut: ebiten.KeyS,

	Up:    ebiten.KeyUp,
	Down:  ebiten.KeyDown,
	Left:  ebiten.KeyLeft,
	Right: ebiten.KeyRight,
}

func ebitenKeyToString(keyboardLayout string, key ebiten.Key) string {
	switch keyboardLayout {
	case keyboardLayoutAzerty:
		return map[ebiten.Key]string{
			ebiten.KeyA:     "Q",
			ebiten.KeyD:     "D",
			ebiten.KeyW:     "Z",
			ebiten.KeyS:     "S",
			ebiten.KeyUp:    "Up",
			ebiten.KeyDown:  "Down",
			ebiten.KeyLeft:  "Left",
			ebiten.KeyRight: "Right",
		}[key]
	case keyboardLayoutQwerty:
		return map[ebiten.Key]string{
			ebiten.KeyA:     "A",
			ebiten.KeyD:     "D",
			ebiten.KeyW:     "W",
			ebiten.KeyS:     "S",
			ebiten.KeyUp:    "Up",
			ebiten.KeyDown:  "Down",
			ebiten.KeyLeft:  "Left",
			ebiten.KeyRight: "Right",
		}[key]
	}

	return ""
}
