package ms2k

const (
	menuStateNewGame = iota
	menuStateExit    = iota
)

var (
	menuStates    = []int8{menuStateNewGame, menuStateExit}
	lenMenuStates = len(menuStates)
)

// MainMenu is the main menu of the game
type MainMenu struct {
	selectedIndex int
}

func (menu *MainMenu) state() int8 {
	index := menu.selectedIndex % lenMenuStates
	if index < 0 {
		index += lenMenuStates
	}
	return menuStates[index]
}
