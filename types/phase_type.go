package types

type Phase int

const (
	UnknownPhase Phase = iota
	DayPhase
	DuskPhase
	NightPhase
)
