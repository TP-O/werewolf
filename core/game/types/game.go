package types

type RoundID uint8

func (r RoundID) IsStarted() bool {
	return r != 0
}

type PhaseID uint8

func (p PhaseID) IsUnknown() bool {
	return p == 0
}

func (p PhaseID) NextPhasePhaseID(lastPhaseID PhaseID) PhaseID {
	if p == lastPhaseID {
		return 1
	}
	return p + 1
}

func (p PhaseID) PreviousPhaseID(lastPhaseID PhaseID) PhaseID {
	if p == 1 {
		return lastPhaseID
	}
	return p - 1
}

type TurnID int8
