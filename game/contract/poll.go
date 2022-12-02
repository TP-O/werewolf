package contract

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type Poll interface {
	// IsOpen returns true if latest poll round is open.
	IsOpen() (bool, enum.Round)

	// CanVote returns true if the elector can vote for the current
	// poll round.
	CanVote(electorID enum.PlayerID) bool

	// Record returns the specific poll round's record.
	Record(round enum.Round) *types.PollRecord

	// Open opens the new poll round if the current one was
	// closed. Returns open status and latest round number.
	Open() (bool, enum.Round)

	// Close closes the current poll round. Returns close
	// status and recent closed poll round's record.
	Close() (bool, *types.PollRecord)

	// AddCandidates adds new candidate to the poll.
	AddCandidates(candidateIDs ...enum.PlayerID)

	// RemoveCandidate removes the candidate from the poll.
	// Return true if successful
	RemoveCandidate(candidateID enum.PlayerID) bool

	// AddElectors adds new electors to the poll.
	// Returns false if the poll's capacity is overloaded.
	AddElectors(electorIDs ...enum.PlayerID) bool

	// RemoveElector removes the elector from the poll.
	// Return true if successful
	RemoveElector(electorID enum.PlayerID) bool

	// SetWeight sets the vote weight for the elector.
	// Default weight is 0.
	SetWeight(electorID enum.PlayerID, weight uint) bool

	// Vote votes or skips for the current poll round.
	Vote(electorID enum.PlayerID, candidateID enum.PlayerID) bool
}
