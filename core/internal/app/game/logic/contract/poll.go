package contract

import "uwwolf/internal/app/game/logic/types"

// VoteRecord contains voting information of a candidate.
type VoteRecord struct {
	// ElectorIDs is player ID list voting for the candidate.
	ElectorIds []types.PlayerId

	// Votes is number of votes.
	Votes uint

	// Weights is the score of the candidate.
	Weights uint
}

// PollRecord contains poll information of a round.
type PollRecord struct {
	// WinnerID is player ID of the winner.
	WinnerId types.PlayerId

	// IsClosed marks that is poll is closed.
	IsClosed bool

	// VoteRecords contains voting information of all candidates.
	VoteRecords map[types.PlayerId]*VoteRecord
}

// Poll manages the voting functionality of a game.
type Poll interface {
	// IsOpen checks if a poll round is opening.
	IsOpen() bool

	// CanVote checks if the elector can vote for the current poll round.
	// Returns the result and an error if any
	CanVote(electorID types.PlayerId) (bool, error)

	// Record returns the record of given round ID.
	// Retun latest round record if `roundID` is 0.
	Record(round types.Round) PollRecord

	// Open starts a new poll round if the current one was closed.
	Open() (bool, error)

	// Close ends the current poll round.
	Close() bool

	// AddCandidates adds new candidate to the poll.
	AddCandidates(candidateIDs ...types.PlayerId)

	// RemoveCandidate removes the candidate from the poll.
	// Return true if successful
	RemoveCandidate(candidateID types.PlayerId) bool

	// AddElectors adds new electors to the poll.
	AddElectors(electorIDs ...types.PlayerId)

	// RemoveElector removes the elector from the poll.
	// Return true if successful
	RemoveElector(electorID types.PlayerId) bool

	// SetWeight sets the voting weight for the elector.
	// Default weight is 0.
	SetWeight(electorID types.PlayerId, weight uint) bool

	// Vote votes or skips for the current poll round.
	Vote(electorID types.PlayerId, candidateID types.PlayerId) (bool, error)
}
