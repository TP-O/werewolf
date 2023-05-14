package types

// VoteRecord contains voting information of a candidate.
type VoteRecord struct {
	// ElectorIds is player ID list voting for the candidate.
	ElectorIds []PlayerId

	// Votes is number of votes.
	Votes uint

	// Weights is the score of the candidate.
	Weights uint
}

// PollRecord contains poll information of a round.
type PollRecord struct {
	// WinnerId is player ID of the winner.
	WinnerId PlayerId

	// IsClosed marks that is poll is closed.
	IsClosed bool

	// VoteRecords contains voting information of all candidates.
	VoteRecords map[PlayerId]*VoteRecord
}
