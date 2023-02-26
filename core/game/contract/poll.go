package contract

import (
	"uwwolf/game/types"
)

type Poll interface {
	// IsOpen returns true if latest poll round is open.
	IsOpen() (bool, types.Round)

	// CanVote returns true if the elector can vote for the current
	// poll round.
	CanVote(electorID types.PlayerID) (bool, error)

	// Record returns the specific poll round's record.
	Record(round types.Round) *types.PollRecord

	// Open opens the new poll round if the current one was
	// closed. Returns open status and latest round number.
	Open() (bool, types.Round)

	// Close closes the current poll round. Returns close
	// status and recent closed poll round's record.
	Close() (bool, *types.PollRecord)

	// AddCandidates adds new candidate to the poll.
	AddCandidates(candidateIDs ...types.PlayerID)

	// RemoveCandidate removes the candidate from the poll.
	// Return true if successful
	RemoveCandidate(candidateID types.PlayerID) bool

	// AddElectors adds new electors to the poll.
	// Returns false if the poll's capacity is overloaded.
	AddElectors(electorIDs ...types.PlayerID) bool

	// RemoveElector removes the elector from the poll.
	// Return true if successful
	RemoveElector(electorID types.PlayerID) bool

	// SetWeight sets the vote weight for the elector.
	// Default weight is 0.
	SetWeight(electorID types.PlayerID, weight uint) bool

	// Vote votes or skips for the current poll round.
	Vote(electorID types.PlayerID, candidateID types.PlayerID) (bool, error)
}
