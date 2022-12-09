package enum

type GameID = string

func IsValidGameID(g GameID) bool {
	return len(g) == 28
}

type PhaseID = uint8

const (
	LowerPhaseID PhaseID = iota
	NightPhaseID
	DayPhaseID
	DuskPhaseID
	UpperPhaseID
)

func IsUnknownPhaseID(p PhaseID) bool {
	return p == 0
}

func NextPhasePhaseID(p PhaseID) PhaseID {
	if p+1 >= UpperPhaseID {
		return LowerPhaseID + 1
	}

	return p + 1
}

func PreviousPhaseID(p PhaseID) PhaseID {
	if p-1 <= LowerPhaseID {
		return UpperPhaseID - 1
	}

	return p - 1
}

type GameStatus = uint8

const (
	Idle GameStatus = iota
	Waiting
	Running
	Finished
)
