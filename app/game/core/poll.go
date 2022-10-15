package core

import (
	"golang.org/x/exp/slices"

	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type poll struct {
	CurrentRoundId      uint                       `json:"currentRoundId"`
	IsVoting            bool                       `json:"isVoting"`
	ElectorIds          []types.PlayerId           `json:"electorIds"`
	RemainingElectorIds []types.PlayerId           `json:"remainingElectorIds"`
	VotedElectorIds     []types.PlayerId           `json:"votedElectorIds"`
	Weights             map[types.PlayerId]uint    `json:"weights"`
	Records             map[uint]*types.PollRecord `json:"rounds"`
}

func NewPoll() contract.Poll {
	return &poll{
		Weights: make(map[types.PlayerId]uint),
		Records: make(map[uint]*types.PollRecord),
	}
}

func (p *poll) IsOpen() bool {
	return p.IsVoting
}

func (p *poll) IsVoted(electorId types.PlayerId) bool {
	return slices.Contains(p.VotedElectorIds, electorId)
}

func (p *poll) IsAllowed(electorId types.PlayerId) bool {
	return slices.Contains(p.RemainingElectorIds, electorId)
}

func (p *poll) AddElectors(electorIds []types.PlayerId) {
	p.ElectorIds = append(p.ElectorIds, electorIds...)
	p.RemainingElectorIds = append(p.RemainingElectorIds, electorIds...)
}

func (p *poll) SetWeight(electorId types.PlayerId, weight uint) bool {
	if !slices.Contains(p.RemainingElectorIds, electorId) {
		return false
	}

	p.Weights[electorId] = weight

	return true
}

func (p *poll) Open() bool {
	if p.IsOpen() || len(p.ElectorIds) <= 2 {
		return false
	}

	p.IsVoting = true
	p.CurrentRoundId += 1
	p.Records[p.CurrentRoundId] = &types.PollRecord{
		Votes: make(map[types.PlayerId]*types.VoteRecord),
	}
	p.VotedElectorIds = make([]types.PlayerId, len(p.ElectorIds))

	return true
}

func (p *poll) Winner() types.PlayerId {
	// Winner can only be found if at least one poll is completed
	// and no poll is opening
	if p.IsOpen() || p.CurrentRoundId == 0 {
		return types.PlayerId("")
	}

	if p.Records[p.CurrentRoundId].IsClosed {
		return p.Records[p.CurrentRoundId].Winner
	}

	winnerId := types.PlayerId("")
	votes := -1
	isFiftyFifty := false

	for candidateId, voteRecord := range p.Records[p.CurrentRoundId].Votes {
		if votes == int(voteRecord.Votes) {
			isFiftyFifty = true
		} else if votes < int(voteRecord.Votes) {
			winnerId = candidateId
			votes = int(voteRecord.Votes)
			isFiftyFifty = false
		}
	}

	if isFiftyFifty {
		return types.PlayerId("")
	}

	return winnerId
}

func (p *poll) Close() types.PollRecord {
	if !p.IsOpen() {
		return *p.Records[p.CurrentRoundId]
	}

	skippedVote := types.PlayerId("")

	// Store skipped votes
	for _, electorId := range p.RemainingElectorIds {
		if !slices.Contains(p.VotedElectorIds, electorId) {
			p.Records[p.CurrentRoundId].Votes[skippedVote].Votes += 1
			p.Records[p.CurrentRoundId].Votes[skippedVote].ElectorIds = append(
				p.Records[p.CurrentRoundId].Votes[skippedVote].ElectorIds,
				electorId,
			)
		}
	}

	p.IsVoting = false
	p.Records[p.CurrentRoundId].Winner = p.Winner()
	p.Records[p.CurrentRoundId].IsClosed = true

	// Remove winner
	if !p.Winner().IsUnknown() {
		winnerIndex := slices.Index(p.RemainingElectorIds, p.Winner())
		p.RemainingElectorIds = slices.Delete(p.RemainingElectorIds, winnerIndex, winnerIndex+1)
	}

	return *p.Records[p.CurrentRoundId]
}

func (p *poll) Vote(electorId types.PlayerId, targetId types.PlayerId) bool {
	if !p.IsOpen() ||
		!p.IsAllowed(electorId) ||
		p.IsVoted(electorId) ||
		!slices.Contains(p.RemainingElectorIds, targetId) {

		return false
	}

	if p.Records[p.CurrentRoundId].Votes[targetId] == nil {
		p.Records[p.CurrentRoundId].Votes[targetId] = &types.VoteRecord{}
	}

	if p.Weights[electorId] == 0 {
		p.Records[p.CurrentRoundId].Votes[targetId].Votes++
	} else {
		p.Records[p.CurrentRoundId].Votes[targetId].Votes += p.Weights[electorId]
	}

	p.Records[p.CurrentRoundId].Votes[targetId].ElectorIds = append(
		p.Records[p.CurrentRoundId].Votes[targetId].ElectorIds,
		electorId,
	)
	p.VotedElectorIds = append(p.VotedElectorIds, electorId)

	return true
}

func (p *poll) Skip(electorId types.PlayerId) bool {
	if !p.IsOpen() || !p.IsAllowed(electorId) || p.IsVoted(electorId) {
		return false
	}

	skipVote := types.PlayerId("")

	if p.Records[p.CurrentRoundId].Votes[skipVote] == nil {
		p.Records[p.CurrentRoundId].Votes[skipVote] = &types.VoteRecord{}
	}

	p.Records[p.CurrentRoundId].Votes[skipVote].Votes++
	p.Records[p.CurrentRoundId].Votes[skipVote].ElectorIds = append(
		p.Records[p.CurrentRoundId].Votes[skipVote].ElectorIds,
		electorId,
	)
	p.VotedElectorIds = append(p.VotedElectorIds, electorId)

	return true
}

func (p *poll) RemoveElector(electorId types.PlayerId) bool {
	if i := slices.Index(p.RemainingElectorIds, electorId); i == -1 {
		return false
	} else {
		p.RemainingElectorIds = slices.Delete(p.RemainingElectorIds, i, i+1)

		return true
	}
}
