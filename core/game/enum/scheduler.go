package enum

type Round = int8

func IsStartedRound(r Round) bool {
	return r != 0
}

const (
	LastRound Round = iota - 1
	_
	FirstRound
)

type Position = int8

const (
	NextPosition Position = iota - 3
	SortedPosition
	LastPosition
	FirstPosition
)
