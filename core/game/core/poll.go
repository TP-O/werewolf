package core

import (
	"fmt"
	"math"
	"uwwolf/game/contract"
	"uwwolf/game/types"

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
	Records map[types.RoundID]types.PollRecord
}

var _ contract.Poll = (*poll)(nil)

func NewPoll() contract.Poll {
	return &poll{
		RoundID: types.RoundID(0),
		Weights: make(map[types.PlayerID]uint),
		Records: make(map[types.RoundID]types.PollRecord),
	}
}

func (p poll) IsOpen() bool {
	return !p.Records[p.RoundID].IsClosed
}

func (p poll) CanVote(electorID types.PlayerID) (bool, error) {
	if !p.IsOpen() {
		return false, fmt.Errorf("Poll (%v) is closed ᕙ(⇀‸↼‶)ᕗ", p.RoundID)
	} else if !slices.Contains(p.RemainingElectorIDs, electorID) {
		return false, fmt.Errorf("You're not allowed to vote ノ(ジ)ー'")
	} else if slices.Contains(p.VotedElectorIDs, electorID) {
		return false, fmt.Errorf("Wait for the next round ಠ_ಠ")
	} else {
		return true, nil
	}
}

func (p poll) Record(roundID types.RoundID) types.PollRecord {
	return p.Records[roundID]
}

func (p *poll) Open() bool {
	if p.IsOpen() {
		return false
	}

	p.RoundID++
	p.Records[p.RoundID] = types.PollRecord{
		VoteRecords: map[types.PlayerID]types.VoteRecord{
			// Empty vote
			"": {
				ElectorIDs: []types.PlayerID{},
			},
		},
	}
	p.VotedElectorIDs = make([]types.PlayerID, 0, len(p.RemainingElectorIDs))

	return true
}

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

func (p *poll) Close() bool {
	if !p.IsOpen() {
		return false
	}

	// Record the skipped votes
	pollRecord := p.Records[p.RoundID]
	emptyVote := types.PlayerID("")
	emptyVoteRecord := pollRecord.VoteRecords[emptyVote]
	for _, electorID := range p.RemainingElectorIDs {
		if !slices.Contains(p.VotedElectorIDs, electorID) {
			emptyVoteRecord.Weights++
			emptyVoteRecord.ElectorIDs = append(
				emptyVoteRecord.ElectorIDs,
				electorID,
			)
			emptyVoteRecord.Votes++
			p.VotedElectorIDs = append(p.VotedElectorIDs, electorID)
		}
	}
	pollRecord.VoteRecords[emptyVote] = emptyVoteRecord
	pollRecord.WinnerID = p.currentWinnerID()
	pollRecord.IsClosed = true
	p.Records[p.RoundID] = pollRecord

	return true
}

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

func (p *poll) RemoveCandidate(candidateID types.PlayerID) bool {
	if i := slices.Index(p.RemainingCandidateIDs, candidateID); i == -1 {
		return false
	} else {
		p.RemainingCandidateIDs = slices.Delete(p.RemainingCandidateIDs, i, i+1)
		return true
	}
}

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

func (p *poll) RemoveElector(electorID types.PlayerID) bool {
	if i := slices.Index(p.RemainingElectorIDs, electorID); i == -1 {
		return false
	} else {
		p.RemainingElectorIDs = slices.Delete(p.RemainingElectorIDs, i, i+1)
		return true
	}
}

func (p *poll) SetWeight(electorID types.PlayerID, weight uint) bool {
	if !slices.Contains(p.RemainingElectorIDs, electorID) {
		return false
	}

	p.Weights[electorID] = weight
	return true
}

func (p *poll) Vote(electorID types.PlayerID, candidateID types.PlayerID) (bool, error) {
	if ok, err := p.CanVote(electorID); !ok {
		return false, err
	} else if !(candidateID.IsUnknown() ||
		slices.Contains(p.RemainingCandidateIDs, candidateID)) {
		return false, fmt.Errorf("Your vote is not valid ¬_¬")
	}

	voteRecord := p.Records[p.RoundID].VoteRecords[candidateID]
	// Empty votes always have weight of 1
	if candidateID.IsUnknown() {
		voteRecord.Weights++
	} else {
		voteRecord.Weights += p.Weights[electorID]
	}

	voteRecord.ElectorIDs = append(
		voteRecord.ElectorIDs,
		electorID,
	)
	voteRecord.Votes++
	p.Records[p.RoundID].VoteRecords[candidateID] = voteRecord
	p.VotedElectorIDs = append(p.VotedElectorIDs, electorID)

	return true, nil
}
