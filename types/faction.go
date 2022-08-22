package types

type FactionId uint

const (
	UnknownFaction FactionId = iota
	VillagerFaction
	WerewolfFaction
)
