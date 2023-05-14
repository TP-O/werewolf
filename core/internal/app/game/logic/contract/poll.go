package contract

import "uwwolf/internal/app/game/logic/types"

// Poll manages the voting mechanism of the game.
type Poll interface {
	// IsOpen checks if a poll round is opening.
	IsOpen() bool

	// CanVote checks if the elector can vote for the current poll round.
	CanVote(electorId types.PlayerId) (bool, error)

	// Record returns the record of given round ID.
	// Retun latest round record if the given`round` is 0.
	Record(round types.Round) *types.PollRecord

	// Open starts a new poll round if the current one was closed.
	Open() (bool, error)

	// Close ends the current poll round.
	Close() bool

	// AddCandidates adds new candidate to the poll.
	AddCandidates(candidateIds ...types.PlayerId)

	// RemoveCandidate removes the candidate from the poll.
	RemoveCandidate(candidateId types.PlayerId) bool

	// AddElectors adds new electors to the poll.
	AddElectors(electorIds ...types.PlayerId)

	// RemoveElector removes the elector from the poll.
	RemoveElector(electorId types.PlayerId) bool

	// SetWeight sets the voting weight for the elector.
	//
	// Default weight is 0.
	SetWeight(electorId types.PlayerId, weight uint) bool

	// Vote votes or skips for the current poll round.
	Vote(electorId types.PlayerId, candidateId types.PlayerId) (bool, error)
}
