package ms2k

import "time"

// Operation holds data about something ongoing
type Operation struct {
	completedPercentage float64
	lastUpdate          time.Time
	speed               float64
}

// Update updates an operation so that its completedPercentage is correct according to the given time
func (op *Operation) Update(now time.Time) {
	elapsedSeconds := now.Sub(op.lastUpdate).Seconds()
	op.completedPercentage += elapsedSeconds * op.speed
	op.lastUpdate = now
}

// IsCompleted returns whether the operation is completed
func (op *Operation) IsCompleted() bool {
	return op.completedPercentage >= 100
}
