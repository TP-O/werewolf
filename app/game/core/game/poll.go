package game

import (
	"errors"
	"math"
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/types"

	"golang.org/x/exp/slices"
)

const MinPollCapacity = 3

type poll struct {
	Round                 types.Round                       `json:"round"`
	CandidateIDs          []types.PlayerID                  `json:"candidateIDs"`
	RemainingCandidateIDs []types.PlayerID                  `json:"remainingCandidateIDs"`
	ElectorIDs            []types.PlayerID                  `json:"electorIDs"`
	RemainingElectorIDs   []types.PlayerID                  `json:"remainingElectorIDs"`
	VotedElectorIDs       []types.PlayerID                  `json:"votedElectorIDs"`
	Capacity              uint                              `json:"capacity"`
	Weights               map[types.PlayerID]uint           `json:"weights"`
	Records               map[types.Round]*types.PollRecord `json:"rounds"`
}

func NewPoll(capacity uint) (contract.Poll, error) {
	if capacity < MinPollCapacity {
		return nil, errors.New("The capacity is too small.")
	}

	return &poll{
		Capacity: capacity,
		Weights:  make(map[types.PlayerID]uint),
		Records:  make(map[types.Round]*types.PollRecord),
	}, nil
}

func (p *poll) IsOpen() (bool, types.Round) {
	isOpen := p.Round.IsStarted() && !p.Records[p.Round].IsClosed

	return isOpen, p.Round
}

func (p *poll) CanVote(electorID types.PlayerID) bool {
	isOpen, _ := p.IsOpen()

	return isOpen &&
		slices.Contains(p.RemainingElectorIDs, electorID) &&
		!slices.Contains(p.VotedElectorIDs, electorID)
}

func (p *poll) Record(round types.Round) *types.PollRecord {
	if !p.Round.IsStarted() || round > p.Round {
		return nil
	} else if round == config.LastRound {
		return p.Records[p.Round]
	} else {
		return p.Records[round]
	}
}

func (p *poll) Open() (bool, types.Round) {
	if isOpen, _ := p.IsOpen(); isOpen || len(p.ElectorIDs) < int(p.Capacity) {
		return false, p.Round
	}

	p.Round++
	p.Records[p.Round] = &types.PollRecord{
		VoteRecords: make(map[types.PlayerID]*types.VoteRecord),
	}
	p.VotedElectorIDs = make([]types.PlayerID, len(p.RemainingElectorIDs))

	return true, p.Round
}

func (p *poll) currentWinnerID() types.PlayerID {
	winnerID := types.PlayerID("")
	halfVotes := uint(math.Round(float64(len(p.RemainingElectorIDs)) / 2))

	for candidateID, record := range p.Records[p.Round].VoteRecords {
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
	if len(p.ElectorIDs)+len(electorIDs) > int(p.Capacity) {
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

func (p *poll) Vote(electorID types.PlayerID, candidateID types.PlayerID) bool {
	if isOpen, _ := p.IsOpen(); !isOpen ||
		p.CanVote(electorID) ||
		!(candidateID.IsUnknown() ||
			slices.Contains(p.RemainingCandidateIDs, candidateID)) {

		return false
	}

	if p.Records[p.Round].VoteRecords[candidateID] == nil {
		p.Records[p.Round].VoteRecords[candidateID] = &types.VoteRecord{}
	}

	p.Records[p.Round].VoteRecords[candidateID].Weights += p.Weights[electorID]
	p.Records[p.Round].VoteRecords[candidateID].ElectorIDs = append(
		p.Records[p.Round].VoteRecords[candidateID].ElectorIDs,
		electorID,
	)
	p.Records[p.Round].VoteRecords[candidateID].Votes++
	p.VotedElectorIDs = append(p.VotedElectorIDs, electorID)

	return true
}
