package game

import (
	"fmt"
	"math"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/util"

	"golang.org/x/exp/slices"
)

type poll struct {
	Round                 types.Round
	CandidateIDs          []types.PlayerID
	RemainingCandidateIDs []types.PlayerID
	ElectorIDs            []types.PlayerID
	RemainingElectorIDs   []types.PlayerID
	VotedElectorIDs       []types.PlayerID
	Capacity              uint8
	Weights               map[types.PlayerID]uint
	Records               map[types.Round]*types.PollRecord
}

func NewPoll(capacity uint8) (contract.Poll, error) {
	if capacity < util.Config().Game.MinPollCapacity {
		return nil, fmt.Errorf("The capacity is too small ¬_¬")
	}

	return &poll{
		Capacity: capacity,
		Weights:  make(map[types.PlayerID]uint),
		Records:  make(map[types.Round]*types.PollRecord),
	}, nil
}

func (p *poll) IsOpen() (bool, types.Round) {
	isOpen := types.IsStartedRound(p.Round) &&
		p.Records[p.Round] != nil &&
		!p.Records[p.Round].IsClosed

	return isOpen, p.Round
}

func (p *poll) CanVote(electorID types.PlayerID) (bool, error) {
	if isOpen, round := p.IsOpen(); !isOpen {
		return false, fmt.Errorf("Poll (%v) is closed ᕙ(⇀‸↼‶)ᕗ", round)
	} else if !slices.Contains(p.RemainingElectorIDs, electorID) {
		return false, fmt.Errorf("You're not allowed to vote ノ(ジ)ー'")
	} else if slices.Contains(p.VotedElectorIDs, electorID) {
		return false, fmt.Errorf("Wait for the next round ಠ_ಠ")
	} else {
		return true, nil
	}
}

func (p *poll) Record(round types.Round) *types.PollRecord {
	if !types.IsStartedRound(p.Round) || round > p.Round {
		return nil
	} else if round == LastRound {
		return p.Records[p.Round]
	} else {
		return p.Records[round]
	}
}

func (p *poll) Open() (bool, types.Round) {
	if isOpen, _ := p.IsOpen(); isOpen ||
		len(p.RemainingElectorIDs) < int(p.Capacity) {
		return false, p.Round
	}

	p.Round++
	p.Records[p.Round] = &types.PollRecord{
		VoteRecords: map[types.PlayerID]*types.VoteRecord{
			// Empty vote
			"": {
				ElectorIDs: []types.PlayerID{},
			},
		},
	}
	p.VotedElectorIDs = make([]types.PlayerID, 0, len(p.RemainingElectorIDs))

	return true, p.Round
}

func (p *poll) currentWinnerID() types.PlayerID {
	winnerID := types.PlayerID("")
	halfVotes := uint(math.Round(float64(len(p.RemainingElectorIDs)) / 2))

	for candidateID, record := range p.Records[p.Round].VoteRecords {
		if record.Weights >= halfVotes {
			if types.IsUnknownPlayerID(winnerID) {
				winnerID = candidateID
			} else {
				// Draw if 2 candidates have overwhelming votes
				return types.PlayerID("")
			}
		}
	}

	return winnerID
}

func (p *poll) Close() (bool, *types.PollRecord) {
	if isOpen, _ := p.IsOpen(); !isOpen {
		return false, nil
	}

	emptyVote := types.PlayerID("")

	// Store skipped votes
	for _, electorID := range p.RemainingElectorIDs {
		if !slices.Contains(p.VotedElectorIDs, electorID) {
			p.Records[p.Round].VoteRecords[emptyVote].Weights++
			p.Records[p.Round].VoteRecords[emptyVote].ElectorIDs = append(
				p.Records[p.Round].VoteRecords[emptyVote].ElectorIDs,
				electorID,
			)
			p.Records[p.Round].VoteRecords[emptyVote].Votes++
			p.VotedElectorIDs = append(p.VotedElectorIDs, electorID)
		}
	}

	p.Records[p.Round].WinnerID = p.currentWinnerID()
	p.Records[p.Round].IsClosed = true

	return true, p.Records[p.Round]
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

func (p *poll) AddElectors(electorIDs ...types.PlayerID) bool {
	if len(p.RemainingElectorIDs)+len(electorIDs) > int(p.Capacity) {
		return false
	}

	for _, electorID := range electorIDs {
		if !slices.Contains(p.RemainingElectorIDs, electorID) {
			p.RemainingElectorIDs = append(p.RemainingElectorIDs, electorID)

			if !slices.Contains(p.ElectorIDs, electorID) {
				p.SetWeight(electorID, 1)
				p.ElectorIDs = append(p.ElectorIDs, electorID)
			}
		}
	}

	return true
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
	if can, err := p.CanVote(electorID); !can {
		return false, err
	} else if !(types.IsUnknownPlayerID(candidateID) ||
		slices.Contains(p.RemainingCandidateIDs, candidateID)) {
		return false, fmt.Errorf("Your vote is not valid ¬_¬")
	}

	if p.Records[p.Round].VoteRecords[candidateID] == nil {
		p.Records[p.Round].VoteRecords[candidateID] = &types.VoteRecord{}
	}

	// Empty votes always have weight of 1
	if types.IsUnknownPlayerID(candidateID) {
		p.Records[p.Round].VoteRecords[candidateID].Weights++
	} else {
		p.Records[p.Round].VoteRecords[candidateID].Weights += p.Weights[electorID]
	}

	p.Records[p.Round].VoteRecords[candidateID].ElectorIDs = append(
		p.Records[p.Round].VoteRecords[candidateID].ElectorIDs,
		electorID,
	)
	p.Records[p.Round].VoteRecords[candidateID].Votes++
	p.VotedElectorIDs = append(p.VotedElectorIDs, electorID)

	return true, nil
}