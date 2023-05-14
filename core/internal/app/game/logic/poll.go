package logic

import (
	"errors"
	"fmt"
	"math"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/pkg/util"

	"golang.org/x/exp/slices"
)

// poll manages the voting mechanism of the game.
type poll struct {
	// Round is the current poll round ID.
	Round types.Round

	// CandidateIDs is the list of all candidates in poll.
	CandidateIds []types.PlayerId

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
	Records map[types.Round]*types.PollRecord
}

func NewPoll() contract.Poll {
	return &poll{
		Round:   constants.ZeroRound,
		Weights: make(map[types.PlayerId]uint),
		Records: make(map[types.Round]*types.PollRecord),
	}
}

// IsOpen checks if a poll round is opening.
func (p poll) IsOpen() bool {
	return !util.IsZero(p.Round) && !p.Records[p.Round].IsClosed
}

// CanVote checks if the elector can vote for the current poll round.
func (p poll) CanVote(electorId types.PlayerId) (bool, error) {
	if !p.IsOpen() {
		return false, fmt.Errorf("Wait for the next poll (%v) ᕙ(⇀‸↼‶)ᕗ", p.Round)
	} else if !slices.Contains(p.RemainingElectorIds, electorId) {
		return false, errors.New("You're not allowed to vote ノ(ジ)ー'")
	} else if slices.Contains(p.VotedElectorIds, electorId) {
		return false, errors.New("Wait for the next round ಠ_ಠ")
	} else {
		return true, nil
	}
}

// Record returns the record of given round ID.
// Retun latest round record if the given`round` is 0.
func (p poll) Record(round types.Round) *types.PollRecord {
	if util.IsZero(round) {
		return p.Record(p.Round)
	}

	return p.Records[round]
}

// Open starts a new poll round if the current one was closed.
func (p *poll) Open() (bool, error) {
	if p.IsOpen() {
		return false, errors.New("Poll is already open!")
	} else if len(p.RemainingCandidateIds) < 2 {
		return false, errors.New("Number of candidates is too small!")
	}

	p.Round++
	p.Records[p.Round] = &types.PollRecord{
		VoteRecords: map[types.PlayerId]*types.VoteRecord{
			// Empty vote
			"": {
				ElectorIds: []types.PlayerId{},
			},
		},
	}

	return true, nil
}

// Close ends the current poll round.
func (p *poll) Close() bool {
	if !p.IsOpen() {
		return false
	}

	emptyVote := types.PlayerId("")

	// Store skipped votes
	for _, electorId := range p.RemainingElectorIds {
		if !slices.Contains(p.VotedElectorIds, electorId) {
			p.Records[p.Round].VoteRecords[emptyVote].Weights++
			p.Records[p.Round].VoteRecords[emptyVote].ElectorIds = append(
				p.Records[p.Round].VoteRecords[emptyVote].ElectorIds,
				electorId,
			)
			p.Records[p.Round].VoteRecords[emptyVote].Votes++
		}
	}

	p.Records[p.Round].WinnerId = p.currentWinnerId()
	p.Records[p.Round].IsClosed = true
	p.VotedElectorIds = make([]types.PlayerId, 0, len(p.RemainingElectorIds))

	return true
}

// currentWinnerID finds the winner in the current poll round.
func (p poll) currentWinnerId() types.PlayerId {
	winnerId := types.PlayerId("")
	halfVotes := uint(math.Round(float64(len(p.RemainingElectorIds)) / 2))

	for candidateId, record := range p.Records[p.Round].VoteRecords {
		if record.Weights >= halfVotes {
			if util.IsZero(winnerId) {
				winnerId = candidateId
			} else {
				// Return draw result if 2 candidates have major votes
				return types.PlayerId("")
			}
		}
	}

	return winnerId
}

// AddCandidates adds new candidate to the poll.
func (p *poll) AddCandidates(candidateIds ...types.PlayerId) {
	for _, candidateId := range candidateIds {
		if !slices.Contains(p.RemainingCandidateIds, candidateId) {
			p.RemainingCandidateIds = append(p.RemainingCandidateIds, candidateId)

			if !slices.Contains(p.CandidateIds, candidateId) {
				p.CandidateIds = append(p.CandidateIds, candidateId)
			}
		}
	}
}

// RemoveCandidate removes the candidate from the poll.
func (p *poll) RemoveCandidate(candidateId types.PlayerId) bool {
	if i := slices.Index(p.RemainingCandidateIds, candidateId); i == -1 {
		return false
	} else {
		p.RemainingCandidateIds = slices.Delete(p.RemainingCandidateIds, i, i+1)
		return true
	}
}

// AddElectors adds new electors to the poll.
func (p *poll) AddElectors(electorIds ...types.PlayerId) {
	for _, electorID := range electorIds {
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
func (p *poll) RemoveElector(electorId types.PlayerId) bool {
	if i := slices.Index(p.RemainingElectorIds, electorId); i == -1 {
		return false
	} else {
		p.RemainingElectorIds = slices.Delete(p.RemainingElectorIds, i, i+1)
		return true
	}
}

// SetWeight sets the voting weight for the elector.
//
// Default weight is 0.
func (p *poll) SetWeight(electorId types.PlayerId, weight uint) bool {
	if !slices.Contains(p.RemainingElectorIds, electorId) {
		return false
	}

	p.Weights[electorId] = weight
	return true
}

// Vote votes or skips for the current poll round.
func (p *poll) Vote(electorId types.PlayerId, candidateId types.PlayerId) (bool, error) {
	if ok, err := p.CanVote(electorId); !ok {
		return false, err
	} else if !(util.IsZero(candidateId) ||
		slices.Contains(p.RemainingCandidateIds, candidateId)) {
		return false, errors.New("Your vote is not valid ¬_¬")
	}

	if p.Records[p.Round].VoteRecords[candidateId] == nil {
		p.Records[p.Round].VoteRecords[candidateId] = &types.VoteRecord{}
	}

	// Empty votes always have weight of 1
	if util.IsZero(candidateId) {
		p.Records[p.Round].VoteRecords[candidateId].Weights++
	} else {
		p.Records[p.Round].VoteRecords[candidateId].Weights += p.Weights[electorId]
	}

	p.Records[p.Round].VoteRecords[candidateId].ElectorIds = append(
		p.Records[p.Round].VoteRecords[candidateId].ElectorIds,
		electorId,
	)
	p.Records[p.Round].VoteRecords[candidateId].Votes++
	p.VotedElectorIds = append(p.VotedElectorIds, electorId)

	return true, nil
}
