package vars

import "uwwolf/game/types"

// Game status
const (
	Idle types.GameStatusID = iota
	Waiting
	Starting
	Finished
)
