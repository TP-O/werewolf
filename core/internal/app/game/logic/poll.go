package logic

import (
	"fmt"
	"math"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"

	"golang.org/x/exp/slices"
)

// poll manages the voting functionality of a game.
type poll struct {
	// Round is the current poll round ID.
	Round types.Round

	// CandidateIDs is the list of all candidates in poll.
	CandidateIDs []types.PlayerId

	// RemainingCandidateIds is the candidate ID list remaining after
	// the most recent closed round.
	RemainingCandidateIds []types.PlayerId

	// ElectorIds is the list of all electors in poll.
	ElectorIds []types.PlayerId

	// RemainingElectorIds is the elector ID list remaining after
	// the most recent closed round.
	RemainingElectorIds []types.PlayerId

	// VotedElectorIds is the voted elector ID list in the opening poll
	// round. Reset every round.
	VotedElectorIds []types.PlayerId

	// Weights stores vote weight of each elector. One weight is equal
	// to one point in poll.
	Weights map[types.PlayerId]uint

	// Records stores poll results in the past.
	// Note: consider using pointer instead.
	Records map[types.Round]*contract.PollRecord
}

func NewPoll() contract.Poll {
	return &poll{
		Round:   constants.ZeroRound,
		Weights: make(map[types.PlayerId]uint),
		Records: make(map[types.Round]*contract.PollRecord),
	}
}

// IsOpen checks if a poll round is opening.
func (p poll) IsOpen() bool {
	return !util.IsZero(p.Round) && !p.Records[p.Round].IsClosed
}

// CanVote checks if the elector can vote for the current poll round.
// Returns the result and an error if any
func (p poll) CanVote(electorID types.PlayerId) (bool, error) {
	if !p.IsOpen() {
		return false, fmt.Errorf("Wait for the next poll (%v) ᕙ(⇀‸↼‶)ᕗ", p.Round)
	} else if !slices.Contains(p.RemainingElectorIds, electorID) {
		return false, fmt.Errorf("You're not allowed to vote ノ(ジ)ー'")
	} else if slices.Contains(p.VotedElectorIds, electorID) {
		return false, fmt.Errorf("Wait for the next round ಠ_ಠ")
	} else {
		return true, nil
	}
}

// Record returns the record of given round ID.
func (p poll) Record(Round types.Round) contract.PollRecord {
	if util.IsZero(Round) {
		return p.Record(p.Round)
	}

	return *p.Records[Round]
}

// Open starts a new poll round if the current one was closed.
func (p *poll) Open() (bool, error) {
	if p.IsOpen() {
		return false, fmt.Errorf("Poll is already open!")
	} else if len(p.RemainingCandidateIds) < 2 {
		return false, fmt.Errorf("Number of candidates is too small!")
	}

	p.Round++
	p.Records[p.Round] = &contract.PollRecord{
		VoteRecords: map[types.PlayerId]*contract.VoteRecord{
			// Empty vote
			"": {
				ElectorIds: []types.PlayerId{},
			},
		},
	}

	return true, nil
}

// currentWinnerID finds the winner in the current poll round.
func (p poll) currentWinnerID() types.PlayerId {
	winnerID := types.PlayerId("")
	halfVotes := uint(math.Round(float64(len(p.RemainingElectorIds)) / 2))

	for candidateID, record := range p.Records[p.Round].VoteRecords {
		if record.Weights >= halfVotes {
			if util.IsZero(winnerID) {
				winnerID = candidateID
			} else {
				// Draw if 2 candidates have overwhelming votes
				return types.PlayerId("")
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

	emptyVote := types.PlayerId("")

	// Store skipped votes
	for _, electorID := range p.RemainingElectorIds {
		if !slices.Contains(p.VotedElectorIds, electorID) {
			p.Records[p.Round].VoteRecords[emptyVote].Weights++
			p.Records[p.Round].VoteRecords[emptyVote].ElectorIds = append(
				p.Records[p.Round].VoteRecords[emptyVote].ElectorIds,
				electorID,
			)
			p.Records[p.Round].VoteRecords[emptyVote].Votes++
		}
	}

	p.Records[p.Round].WinnerId = p.currentWinnerID()
	p.Records[p.Round].IsClosed = true
	p.VotedElectorIds = make([]types.PlayerId, 0, len(p.RemainingElectorIds))

	return true
}

// AddCandidates adds new candidate to the poll.
func (p *poll) AddCandidates(candidateIDs ...types.PlayerId) {
	for _, candidateID := range candidateIDs {
		if !slices.Contains(p.RemainingCandidateIds, candidateID) {
			p.RemainingCandidateIds = append(p.RemainingCandidateIds, candidateID)

			if !slices.Contains(p.CandidateIDs, candidateID) {
				p.CandidateIDs = append(p.CandidateIDs, candidateID)
			}
		}
	}
}

// RemoveCandidate removes the candidate from the poll.
// Return true if successful
func (p *poll) RemoveCandidate(candidateID types.PlayerId) bool {
	if i := slices.Index(p.RemainingCandidateIds, candidateID); i == -1 {
		return false
	} else {
		p.RemainingCandidateIds = slices.Delete(p.RemainingCandidateIds, i, i+1)
		return true
	}
}

// AddElectors adds new electors to the poll.
func (p *poll) AddElectors(ElectorIds ...types.PlayerId) {
	for _, electorID := range ElectorIds {
		if !slices.Contains(p.RemainingElectorIds, electorID) {
			p.RemainingElectorIds = append(p.RemainingElectorIds, electorID)

			if !slices.Contains(p.ElectorIds, electorID) {
				p.SetWeight(electorID, 1)
				p.ElectorIds = append(p.ElectorIds, electorID)
			}
		}
	}
}

// RemoveElector removes the elector from the poll.
// Return true if successful
func (p *poll) RemoveElector(electorID types.PlayerId) bool {
	if i := slices.Index(p.RemainingElectorIds, electorID); i == -1 {
		return false
	} else {
		p.RemainingElectorIds = slices.Delete(p.RemainingElectorIds, i, i+1)
		return true
	}
}

// SetWeight sets the voting weight for the elector.
// Default weight is 0.
func (p *poll) SetWeight(electorID types.PlayerId, weight uint) bool {
	if !slices.Contains(p.RemainingElectorIds, electorID) {
		return false
	}

	p.Weights[electorID] = weight
	return true
}

// Vote votes or skips for the current poll round.
func (p *poll) Vote(electorID types.PlayerId, candidateID types.PlayerId) (bool, error) {
	if ok, err := p.CanVote(electorID); !ok {
		return false, err
	} else if !(util.IsZero(candidateID) ||
		slices.Contains(p.RemainingCandidateIds, candidateID)) {
		return false, fmt.Errorf("Your vote is not valid ¬_¬")
	}

	if p.Records[p.Round].VoteRecords[candidateID] == nil {
		p.Records[p.Round].VoteRecords[candidateID] = &contract.VoteRecord{}
	}

	// Empty votes always have weight of 1
	if util.IsZero(candidateID) {
		p.Records[p.Round].VoteRecords[candidateID].Weights++
	} else {
		p.Records[p.Round].VoteRecords[candidateID].Weights += p.Weights[electorID]
	}

	p.Records[p.Round].VoteRecords[candidateID].ElectorIds = append(
		p.Records[p.Round].VoteRecords[candidateID].ElectorIds,
		electorID,
	)
	p.Records[p.Round].VoteRecords[candidateID].Votes++
	p.VotedElectorIds = append(p.VotedElectorIds, electorID)

	return true, nil
}
