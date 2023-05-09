package declare

import "uwwolf/internal/app/game/logic/types"

// Number of remaining times
const (
	UnlimitedTimes types.Times = iota - 1
	OutOfTimes
	Once
	Twice
)
