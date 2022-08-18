package types

type FactionId uint

const (
	UnknownFaction FactionId = iota
	VillageFaction
	WerewolfFaction
	IndependentFaction
)
