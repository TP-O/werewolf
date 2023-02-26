package aaa

type PlayerStatus = uint8

func IsUnknownPlayerID(p PlayerID) bool {
	return p == ""
}

type VoteRecord struct {
	ElectorIDs []PlayerID
	Votes      uint
	Weights    uint
}

type PollRecord struct {
	WinnerID    PlayerID
	IsClosed    bool
	VoteRecords map[PlayerID]*VoteRecord
}

func IsStartedRound(r Round) bool {
	return r != 0
}

type Position = int8
