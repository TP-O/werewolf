package enum

type GameID string

func (g GameID) IsValid() bool {
	return len(g) == 28
}

type PhaseID uint

const (
	LowerPhaseID PhaseID = iota
	NightPhaseID
	DayPhaseID
	DuskPhaseID
	UpperPhaseID
)

func (p PhaseID) IsUnknown() bool {
	return p == 0
}

func (p PhaseID) NextPhase() PhaseID {
	if p+1 >= UpperPhaseID {
		return LowerPhaseID + 1
	}

	return p + 1
}

func (p PhaseID) PreviousPhase() PhaseID {
	if p-1 <= LowerPhaseID {
		return UpperPhaseID - 1
	}

	return p - 1
}

type GameStatus uint

const (
	Idle GameStatus = iota
	Waiting
	Starting
	Finished
)
