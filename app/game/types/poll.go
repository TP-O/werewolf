package types

type VoteRecord struct {
	ElectorIDs []PlayerID `json:"electorIDs"`
	Votes      uint       `json:"votes"`
	Weights    uint       `json:"weight"`
}

type PollRecord struct {
	WinnerID    PlayerID                 `json:"winnerID"`
	IsClosed    bool                     `json:"isClosed"`
	VoteRecords map[PlayerID]*VoteRecord `json:"voteRecords"`
}
