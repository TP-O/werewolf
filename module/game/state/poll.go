package state

import (
	"golang.org/x/exp/slices"

	"uwwolf/types"
)

type Poll struct {
	CurrentRoundId  uint                    `json:"current_round_id"`
	IsVoting        bool                    `json:"is_voting"`
	ElectorIds      []types.PlayerId        `json:"elector_ids"`
	VotedElectorIds []types.PlayerId        `json:"voted_elector_ids"`
	Weights         map[types.PlayerId]uint `json:"weights"`
	Rounds          map[uint]round          `json:"rounds"`
}

type round = map[types.PlayerId]*record

type record struct {
	ElectorIds []types.PlayerId `json:"elector_ids"`
	Votes      uint             `json:"votes"`
}

func NewPoll() *Poll {
	return &Poll{
		Weights: make(map[types.PlayerId]uint),
		Rounds:  make(map[uint]round),
	}
}

func (p *Poll) IsOpen() bool {
	return p.IsVoting
}

func (p *Poll) IsVoted(electorId types.PlayerId) bool {
	return slices.Contains(p.VotedElectorIds, electorId)
}

func (p *Poll) IsAllowed(electorId types.PlayerId) bool {
	return slices.Contains(p.ElectorIds, electorId)
}

func (p *Poll) CurrentRound() round {
	return p.Rounds[p.CurrentRoundId]
}

func (p *Poll) AddElectors(electorIds []types.PlayerId) {
	p.ElectorIds = append(p.ElectorIds, electorIds...)
}

func (p *Poll) SetWeight(electorId types.PlayerId, weight uint) bool {
	if !slices.Contains(p.ElectorIds, electorId) {
		return false
	}

	p.Weights[electorId] = weight

	return true
}

func (p *Poll) Open() bool {
	if p.IsOpen() || len(p.ElectorIds) <= 2 {
		return false
	}

	p.IsVoting = true
	p.CurrentRoundId += 1
	p.Rounds[p.CurrentRoundId] = make(map[types.PlayerId]*record)
	p.VotedElectorIds = make([]types.PlayerId, len(p.ElectorIds))

	return true
}

func (p *Poll) Close() map[types.PlayerId]*record {
	currentRound := p.CurrentRound()

	if !p.IsOpen() {
		return currentRound
	}

	p.IsVoting = false

	// Store skipped votes
	currentRound[types.UnknownPlayer] = &record{}

	for _, elector := range p.ElectorIds {
		if !slices.Contains(p.VotedElectorIds, elector) {
			currentRound[types.UnknownPlayer].Votes += 1
			currentRound[types.UnknownPlayer].ElectorIds = append(
				currentRound[types.UnknownPlayer].ElectorIds,
				elector,
			)
		}
	}

	return p.CurrentRound()
}

func (p *Poll) Vote(electorId types.PlayerId, targetId types.PlayerId) bool {
	if !p.IsOpen() || !p.IsAllowed(electorId) || p.IsVoted(electorId) {
		return false
	}

	currentRound := p.CurrentRound()

	if currentRound[targetId] == nil {
		currentRound[targetId] = &record{}
	}

	if targetId.IsUnknown() ||
		(!targetId.IsUnknown() && p.Weights[electorId] == 0) {

		currentRound[targetId].Votes++
	} else {
		currentRound[targetId].Votes += p.Weights[electorId]
	}

	currentRound[targetId].ElectorIds = append(
		currentRound[targetId].ElectorIds,
		electorId,
	)

	p.VotedElectorIds = append(p.VotedElectorIds, electorId)

	return true
}

func (p *Poll) Winner() types.PlayerId {
	if p.IsOpen() || p.CurrentRoundId == 0 {
		return types.UnknownPlayer
	}

	loserId := types.UnknownPlayer
	votes := -1
	isFiftyFifty := false

	for targetId, target := range p.CurrentRound() {
		if votes == int(target.Votes) {
			isFiftyFifty = true
		} else if votes < int(target.Votes) {
			loserId = targetId
			votes = int(target.Votes)
			isFiftyFifty = false
		}
	}

	if isFiftyFifty {
		return types.UnknownPlayer
	}

	return loserId
}

func (p *Poll) RemoveElector(electorId types.PlayerId) bool {
	if i := slices.Index(p.ElectorIds, electorId); i == -1 {
		return false
	} else {
		p.ElectorIds = slices.Delete(p.ElectorIds, i, i+1)

		return true
	}
}
