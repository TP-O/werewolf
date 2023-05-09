package declare

import "uwwolf/internal/app/game/logic/types"

const (
	Idle types.GameStatusID = iota
	Waiting
	Starting
	Finished
)
