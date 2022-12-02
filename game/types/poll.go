package types

import "uwwolf/game/enum"

type VoteRecord struct {
	ElectorIDs []enum.PlayerID
	Votes      uint
	Weights    uint
}

type PollRecord struct {
	WinnerID    enum.PlayerID
	IsClosed    bool
	VoteRecords map[enum.PlayerID]*VoteRecord
}
