package state

import (
	"golang.org/x/exp/slices"

	"uwwolf/types"
)

type Poll struct {
	currentResultId uint
	isVoting        bool
	electorIds      []types.PlayerId
	votedElectorIds []types.PlayerId
	results         map[uint]result
}

type result = map[types.PlayerId]*record

type record struct {
	electorIds []types.PlayerId
	votes      uint
}

func (r record) ElectorIds() []types.PlayerId {
	return r.electorIds
}

func (r record) Votes() uint {
	return r.votes
}

func NewPoll(electorIds []types.PlayerId) *Poll {
	if len(electorIds) < 3 {
		return nil
	}

	return &Poll{
		electorIds:      electorIds,
		votedElectorIds: make([]types.PlayerId, len(electorIds)),
		results:         make(map[uint]result),
	}
}

func (p *Poll) IsOpen() bool {
	return p.isVoting
}

func (p *Poll) IsAllowed(electorId types.PlayerId) bool {
	return slices.Contains(p.electorIds, electorId) &&
		!slices.Contains(p.votedElectorIds, electorId)
}

func (p *Poll) GetCurrentResult() result {
	return p.results[p.currentResultId]
}

func (p *Poll) Open() bool {
	if p.IsOpen() {
		return false
	}

	p.isVoting = true
	p.currentResultId += 1
	p.results[p.currentResultId] = make(map[types.PlayerId]*record)
	p.votedElectorIds = make([]types.PlayerId, len(p.electorIds))

	return true
}

func (p *Poll) Close() map[types.PlayerId]*record {
	currentResult := p.GetCurrentResult()

	if !p.IsOpen() {
		return currentResult
	}

	p.isVoting = false

	// Store skipped votes
	currentResult[types.UnknownPlayer] = &record{}

	for _, elector := range p.electorIds {
		if !slices.Contains(p.votedElectorIds, elector) {
			currentResult[types.UnknownPlayer].votes += 1
			currentResult[types.UnknownPlayer].electorIds = append(
				currentResult[types.UnknownPlayer].electorIds,
				elector,
			)
		}
	}

	return p.GetCurrentResult()
}

func (p *Poll) Vote(electorId types.PlayerId, targetId types.PlayerId) bool {
	if !p.IsOpen() || !p.IsAllowed(electorId) {
		return false
	}

	currentResult := p.GetCurrentResult()

	if currentResult[targetId] == nil {
		currentResult[targetId] = &record{}
	}

	currentResult[targetId].votes += 1
	currentResult[targetId].electorIds = append(
		currentResult[targetId].electorIds,
		electorId,
	)

	p.votedElectorIds = append(p.votedElectorIds, electorId)

	return true
}

func (p *Poll) GetLoser() types.PlayerId {
	if p.IsOpen() || p.currentResultId == 0 {
		return types.UnknownPlayer
	}

	loserId := types.UnknownPlayer
	votes := -1
	isFiftyFifty := false

	for targetId, target := range p.GetCurrentResult() {
		if votes == int(target.votes) {
			isFiftyFifty = true
		} else if votes < int(target.votes) {
			loserId = targetId
			votes = int(target.votes)
			isFiftyFifty = false
		}
	}

	if isFiftyFifty {
		return types.UnknownPlayer
	}

	return loserId
}

func (p *Poll) RemoveElector(electorId types.PlayerId) bool {
	if i := slices.Index(p.electorIds, electorId); i == -1 {
		return false
	} else {
		p.electorIds = slices.Delete(p.electorIds, i, i+1)

		return true
	}
}
