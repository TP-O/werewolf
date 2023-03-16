package contract

import "uwwolf/game/types"

// Poll manages the voting functionality of a game.
type Poll interface {
	// IsOpen checks if a poll round is opening.
	IsOpen() bool

	// CanVote checks if the elector can vote for the current poll round.
	// Returns the result and an error if any
	CanVote(electorID types.PlayerID) (bool, error)

	// Record returns the record of given round ID.
	// Retun latest round record if `roundID` is 0.
	Record(roundID types.RoundID) *types.PollRecord

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
