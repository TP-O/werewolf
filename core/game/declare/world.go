package declare

import "uwwolf/game/types"

const (
	Idle types.GameStatusID = iota
	Waiting
	Starting
	Finished
)
