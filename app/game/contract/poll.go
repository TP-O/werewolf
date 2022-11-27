package contract

import "uwwolf/app/game/types"

type Poll interface {
	// IsOpen returns true if current poll round is opening.
	IsOpen() (bool, types.Round)

	// CanVote returns true if the elector voted the current
	// poll round.
	CanVote(electorID types.PlayerID) bool

	// Winner returns the winner of the latest closed poll round.
	Record(round types.Round) *types.PollRecord

	// Open starts the new poll round if the current poll round is
	// closed. Returns fasle if current round is still opening.
	Open() (bool, types.Round)

	// Close stops the current poll round and return closed round.
	Close() (bool, *types.PollRecord)

	AddCandidates(candidateIDs ...types.PlayerID)

	RemoveCandidate(candidateID types.PlayerID) bool

	// AddElectors adds the electors to the poll. This method
	// does not check unique elector id.
	AddElectors(electorIDs ...types.PlayerID) bool

	// RemoveElector removes the elector from the poll. Returns
	// false if elector does not exist.
	RemoveElector(electorID types.PlayerID) bool

	// SetWeight sets the vote weight for the elector. Default
	// weight is 0.
	SetWeight(electorID types.PlayerID, weight uint) bool

	// Vote adds vote to current poll round.
	Vote(electorID types.PlayerID, candidateID types.PlayerID) bool
}
