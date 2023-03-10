package types

// VoteRecord contains voting information of a candidate.
type VoteRecord struct {
	// ElectorIDs is player ID list voting for the candidate.
	ElectorIDs []PlayerID

	// Votes is number of votes.
	Votes uint

	// Weights is the score of the candidate.
	Weights uint
}

// PollRecord contains poll information of a round.
type PollRecord struct {
	// WinnerID is player ID of the winner.
	WinnerID PlayerID

	// IsClosed marks that is poll is closed.
	IsClosed bool

	// VoteRecords contains voting information of all candidates.
	VoteRecords map[PlayerID]*VoteRecord
}
