package contract

import "uwwolf/app/types"

type Poll interface {
	// IsOpen returns true if current poll round is opening.
	IsOpen() bool

	// IsVoted returns true if the elector voted the current
	// poll round.
	IsVoted(electorId types.PlayerId) bool

	// IsAllowed returns true if the elector is allowed to
	// vote the current poll.
	IsAllowed(electorId types.PlayerId) bool

	// CurrentRound returns the current round poll.
	CurrentRound() types.PollRound

	// AddElectors adds the electors to the poll. This method
	// does not check unique elector id.
	AddElectors(electorIds []types.PlayerId)

	// SetWeight sets the vote weight for the elector. Default
	// weight is 0.
	SetWeight(electorId types.PlayerId, weight uint) bool

	// Open starts the new poll round if the current poll round is
	// closed. Returns fasle if current round is still opening.
	Open() bool

	// Close stops the current poll round and return closed round.
	Close() types.PollRound

	// Vote adds vote to current poll round.
	Vote(electorId types.PlayerId, targetId types.PlayerId) bool

	// Winner returns the winner of the latest closed poll round.
	Winner() types.PlayerId

	// RemoveElector removes the elector from the poll. Returns
	// false if elector does not exist.
	RemoveElector(electorId types.PlayerId) bool
}
