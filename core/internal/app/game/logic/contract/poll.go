package contract

import "uwwolf/internal/app/game/logic/types"

// VoteRecord contains voting information of a candidate.
type VoteRecord struct {
	// ElectorIDs is player ID list voting for the candidate.
	ElectorIDs []types.PlayerID

	// Votes is number of votes.
	Votes uint

	// Weights is the score of the candidate.
	Weights uint
}

// PollRecord contains poll information of a round.
type PollRecord struct {
	// WinnerID is player ID of the winner.
	WinnerID types.PlayerID

	// IsClosed marks that is poll is closed.
	IsClosed bool

	// VoteRecords contains voting information of all candidates.
	VoteRecords map[types.PlayerID]*VoteRecord
}

// Poll manages the voting functionality of a game.
type Poll interface {
	// IsOpen checks if a poll round is opening.
	IsOpen() bool

	// CanVote checks if the elector can vote for the current poll round.
	// Returns the result and an error if any
	CanVote(electorID types.PlayerID) (bool, error)

	// Record returns the record of given round ID.
	// Retun latest round record if `roundID` is 0.
	Record(roundID types.RoundID) PollRecord

	// Open starts a new poll round if the current one was closed.
	Open() (bool, error)

	// Close ends the current poll round.
	Close() bool

	// AddCandidates adds new candidate to the poll.
	AddCandidates(candidateIDs ...types.PlayerID)

	// RemoveCandidate removes the candidate from the poll.
	// Return true if successful
	RemoveCandidate(candidateID types.PlayerID) bool

	// AddElectors adds new electors to the poll.
	AddElectors(electorIDs ...types.PlayerID)

	// RemoveElector removes the elector from the poll.
	// Return true if successful
	RemoveElector(electorID types.PlayerID) bool

	// SetWeight sets the voting weight for the elector.
	// Default weight is 0.
	SetWeight(electorID types.PlayerID, weight uint) bool

	// Vote votes or skips for the current poll round.
	Vote(electorID types.PlayerID, candidateID types.PlayerID) (bool, error)
}
