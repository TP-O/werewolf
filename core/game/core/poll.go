package core

import (
	"fmt"
	"math"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"

	"golang.org/x/exp/slices"
)

// poll manages the voting functionality of a game.
type poll struct {
	// RoundID is the current poll round ID.
	RoundID types.RoundID

	// CandidateIDs is the list of all candidates in poll.
	CandidateIDs []types.PlayerID

	// RemainingCandidateIDs is the candidate ID list remaining after
	// the most recent closed round.
	RemainingCandidateIDs []types.PlayerID

	// ElectorIDs is the list of all electors in poll.
	ElectorIDs []types.PlayerID

	// RemainingElectorIDs is the elector ID list remaining after
	// the most recent closed round.
	RemainingElectorIDs []types.PlayerID

	// VotedElectorIDs is the voted elector ID list in the opening poll
	// round. Reset every round.
	VotedElectorIDs []types.PlayerID

	// Weights stores vote weight of each elector. One weight is equal
	// to one point in poll.
	Weights map[types.PlayerID]uint

	// Records stores poll results in the past.
	// Note: consider using pointer instead.
	Records map[types.RoundID]*types.PollRecord
}

func NewPoll() contract.Poll {
	return &poll{
		RoundID: vars.ZeroRound,
		Weights: make(map[types.PlayerID]uint),
		Records: make(map[types.RoundID]*types.PollRecord),
	}
}

// IsOpen checks if a poll round is opening.
func (p poll) IsOpen() bool {
	return p.RoundID != vars.ZeroRound && !p.Records[p.RoundID].IsClosed
}

// CanVote checks if the elector can vote for the current poll round.
// Returns the result and an error if any
func (p poll) CanVote(electorID types.PlayerID) (bool, error) {
	if !p.IsOpen() {
		return false, fmt.Errorf("Wait for the next poll (%v) ᕙ(⇀‸↼‶)ᕗ", p.RoundID)
	} else if !slices.Contains(p.RemainingElectorIDs, electorID) {
		return false, fmt.Errorf("You're not allowed to vote ノ(ジ)ー'")
	} else if slices.Contains(p.VotedElectorIDs, electorID) {
		return false, fmt.Errorf("Wait for the next round ಠ_ಠ")
	} else {
		return true, nil
	}
}

// Record returns the record of given round ID.
func (p poll) Record(roundID types.RoundID) *types.PollRecord {
	return p.Records[roundID]
}

// Open starts a new poll round if the current one was closed.
func (p *poll) Open() (bool, error) {
	if p.IsOpen() {
		return false, fmt.Errorf("Poll is already open!")
	} else if len(p.RemainingCandidateIDs) < 2 {
		return false, fmt.Errorf("Number of candidates is too small!")
	}

	p.RoundID++
	p.Records[p.RoundID] = &types.PollRecord{
		VoteRecords: map[types.PlayerID]*types.VoteRecord{
			// Empty vote
			"": {
				ElectorIDs: []types.PlayerID{},
			},
		},
	}

	return true, nil
}

// currentWinnerID finds the winner in the current poll round.
func (p poll) currentWinnerID() types.PlayerID {
	winnerID := types.PlayerID("")
	halfVotes := uint(math.Round(float64(len(p.RemainingElectorIDs)) / 2))

	for candidateID, record := range p.Records[p.RoundID].VoteRecords {
		if record.Weights >= halfVotes {
			if winnerID.IsUnknown() {
				winnerID = candidateID
			} else {
				// Draw if 2 candidates have overwhelming votes
				return types.PlayerID("")
			}
		}
	}

	return winnerID
}

// Close finishes the current poll round.
func (p *poll) Close() bool {
	if !p.IsOpen() {
		return false
	}

	emptyVote := types.PlayerID("")

	// Store skipped votes
	for _, electorID := range p.RemainingElectorIDs {
		if !slices.Contains(p.VotedElectorIDs, electorID) {
			p.Records[p.RoundID].VoteRecords[emptyVote].Weights++
			p.Records[p.RoundID].VoteRecords[emptyVote].ElectorIDs = append(
				p.Records[p.RoundID].VoteRecords[emptyVote].ElectorIDs,
				electorID,
			)
			p.Records[p.RoundID].VoteRecords[emptyVote].Votes++
		}
	}

	p.Records[p.RoundID].WinnerID = p.currentWinnerID()
	p.Records[p.RoundID].IsClosed = true
	p.VotedElectorIDs = make([]types.PlayerID, 0, len(p.RemainingElectorIDs))

	return true
}

// AddCandidates adds new candidate to the poll.
func (p *poll) AddCandidates(candidateIDs ...types.PlayerID) {
	for _, candidateID := range candidateIDs {
		if !slices.Contains(p.RemainingCandidateIDs, candidateID) {
			p.RemainingCandidateIDs = append(p.RemainingCandidateIDs, candidateID)

			if !slices.Contains(p.CandidateIDs, candidateID) {
				p.CandidateIDs = append(p.CandidateIDs, candidateID)
			}
		}
	}
}

// RemoveCandidate removes the candidate from the poll.
// Return true if successful
func (p *poll) RemoveCandidate(candidateID types.PlayerID) bool {
	if i := slices.Index(p.RemainingCandidateIDs, candidateID); i == -1 {
		return false
	} else {
		p.RemainingCandidateIDs = slices.Delete(p.RemainingCandidateIDs, i, i+1)
		return true
	}
}

// AddElectors adds new electors to the poll.
func (p *poll) AddElectors(electorIDs ...types.PlayerID) {
	for _, electorID := range electorIDs {
		if !slices.Contains(p.RemainingElectorIDs, electorID) {
			p.RemainingElectorIDs = append(p.RemainingElectorIDs, electorID)

			if !slices.Contains(p.ElectorIDs, electorID) {
				p.SetWeight(electorID, 1)
				p.ElectorIDs = append(p.ElectorIDs, electorID)
			}
		}
	}
}

// RemoveElector removes the elector from the poll.
// Return true if successful
func (p *poll) RemoveElector(electorID types.PlayerID) bool {
	if i := slices.Index(p.RemainingElectorIDs, electorID); i == -1 {
		return false
	} else {
		p.RemainingElectorIDs = slices.Delete(p.RemainingElectorIDs, i, i+1)
		return true
	}
}

// SetWeight sets the voting weight for the elector.
// Default weight is 0.
func (p *poll) SetWeight(electorID types.PlayerID, weight uint) bool {
	if !slices.Contains(p.RemainingElectorIDs, electorID) {
		return false
	}

	p.Weights[electorID] = weight
	return true
}

// Vote votes or skips for the current poll round.
func (p *poll) Vote(electorID types.PlayerID, candidateID types.PlayerID) (bool, error) {
	if ok, err := p.CanVote(electorID); !ok {
		return false, err
	} else if !(candidateID.IsUnknown() ||
		slices.Contains(p.RemainingCandidateIDs, candidateID)) {
		return false, fmt.Errorf("Your vote is not valid ¬_¬")
	}

	if p.Records[p.RoundID].VoteRecords[candidateID] == nil {
		p.Records[p.RoundID].VoteRecords[candidateID] = &types.VoteRecord{}
	}

	// Empty votes always have weight of 1
	if candidateID.IsUnknown() {
		p.Records[p.RoundID].VoteRecords[candidateID].Weights++
	} else {
		p.Records[p.RoundID].VoteRecords[candidateID].Weights += p.Weights[electorID]
	}

	p.Records[p.RoundID].VoteRecords[candidateID].ElectorIDs = append(
		p.Records[p.RoundID].VoteRecords[candidateID].ElectorIDs,
		electorID,
	)
	p.Records[p.RoundID].VoteRecords[candidateID].Votes++
	p.VotedElectorIDs = append(p.VotedElectorIDs, electorID)

	return true, nil
}
