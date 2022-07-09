package ms2k

// Radar is the result of a radar scan at a given position
type Radar struct {
	Position Position

	Scan Operation
}
