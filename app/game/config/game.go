package config

import "uwwolf/app/game/types"

const (
	LowerPhaseID types.PhaseID = iota
	NightPhaseID
	DayPhaseID
	DuskPhaseID
	UpperPhaseID
)

const PreparationTime = 10 // (seconds)

const (
	Idle types.GameStatus = iota
	Waiting
	Starting
	Finished
)

const MinPollCapacity = 3
