package config

import "uwwolf/app/game/types"

const (
	OnlineStatus types.Status = iota + 1
	BusyStatus
	InGameStatus
)
