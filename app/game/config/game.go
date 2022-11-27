package config

import "uwwolf/app/game/types"

const (
	NightPhaseID types.PhaseID = iota + 1
	DayPhaseID
	DuskPhaseID
)

const PreparationTime = 10 // (seconds)

const (
	Idle types.GameStatus = iota
	Waiting
	Starting
	Finished
)
