package enum

type Round int

func (r Round) IsStarted() bool {
	return r != 0
}

const (
	LastRound Round = iota - 1
	_
	FirstRound
)

type Position int

const (
	NextPosition Position = iota - 3
	SortedPosition
	LastPosition
	FirstPosition
)
