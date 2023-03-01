package types

type VoteRecord struct {
	ElectorIDs []PlayerID
	Votes      uint
	Weights    uint
}

type PollRecord struct {
	WinnerID    PlayerID
	IsClosed    bool
	VoteRecords map[PlayerID]VoteRecord
}
