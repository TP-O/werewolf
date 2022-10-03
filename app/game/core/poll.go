package core

import (
	"golang.org/x/exp/slices"

	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type poll struct {
	CurrentpollRoundId uint                     `json:"currentRoundId"`
	IsVoting           bool                     `json:"isVoting"`
	ElectorIds         []types.PlayerId         `json:"electorIds"`
	VotedElectorIds    []types.PlayerId         `json:"votedElectorIds"`
	Weights            map[types.PlayerId]uint  `json:"weights"`
	Rounds             map[uint]types.PollRound `json:"rounds"`
}

func NewPoll() contract.Poll {
	return &poll{
		Weights: make(map[types.PlayerId]uint),
		Rounds:  make(map[uint]types.PollRound),
	}
}

func (p *poll) IsOpen() bool {
	return p.IsVoting
}

func (p *poll) IsVoted(electorId types.PlayerId) bool {
	return slices.Contains(p.VotedElectorIds, electorId)
}

func (p *poll) IsAllowed(electorId types.PlayerId) bool {
	return slices.Contains(p.ElectorIds, electorId)
}

func (p *poll) CurrentRound() types.PollRound {
	return p.Rounds[p.CurrentpollRoundId]
}

func (p *poll) AddElectors(electorIds []types.PlayerId) {
	p.ElectorIds = append(p.ElectorIds, electorIds...)
}

func (p *poll) SetWeight(electorId types.PlayerId, weight uint) bool {
	if !slices.Contains(p.ElectorIds, electorId) {
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
	p.CurrentpollRoundId += 1
	p.Rounds[p.CurrentpollRoundId] = make(types.PollRound)
	p.VotedElectorIds = make([]types.PlayerId, len(p.ElectorIds))

	return true
}

func (p *poll) Close() types.PollRound {
	currentpollRound := p.CurrentRound()

	if !p.IsOpen() {
		return currentpollRound
	}

	p.IsVoting = false

	// Store skipped votes
	currentpollRound[types.PlayerId("")] = &types.PollRecord{}

	for _, elector := range p.ElectorIds {
		if !slices.Contains(p.VotedElectorIds, elector) {
			currentpollRound[types.PlayerId("")].Votes += 1
			currentpollRound[types.PlayerId("")].ElectorIds = append(
				currentpollRound[types.PlayerId("")].ElectorIds,
				elector,
			)
		}
	}

	return p.CurrentRound()
}

func (p *poll) Vote(electorId types.PlayerId, targetId types.PlayerId) bool {
	if !p.IsOpen() || !p.IsAllowed(electorId) || p.IsVoted(electorId) {
		return false
	}

	currentpollRound := p.CurrentRound()

	if currentpollRound[targetId] == nil {
		currentpollRound[targetId] = &types.PollRecord{}
	}

	if targetId.IsUnknown() ||
		(!targetId.IsUnknown() && p.Weights[electorId] == 0) {

		currentpollRound[targetId].Votes++
	} else {
		currentpollRound[targetId].Votes += p.Weights[electorId]
	}

	currentpollRound[targetId].ElectorIds = append(
		currentpollRound[targetId].ElectorIds,
		electorId,
	)

	p.VotedElectorIds = append(p.VotedElectorIds, electorId)

	return true
}

func (p *poll) Winner() types.PlayerId {
	if p.IsOpen() || p.CurrentpollRoundId == 0 {
		return types.PlayerId("")
	}

	loserId := types.PlayerId("")
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
		return types.PlayerId("")
	}

	return loserId
}

func (p *poll) RemoveElector(electorId types.PlayerId) bool {
	if i := slices.Index(p.ElectorIds, electorId); i == -1 {
		return false
	} else {
		p.ElectorIds = slices.Delete(p.ElectorIds, i, i+1)

		return true
	}
}
