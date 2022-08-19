package state

import (
	"golang.org/x/exp/slices"

	"uwwolf/types"
)

type Poll struct {
	recordId      uint
	isVoting      bool
	electors      []types.PlayerId
	votedElectors []types.PlayerId
	records       map[uint]map[types.PlayerId]*record
}

type record = struct {
	electors []types.PlayerId
	votes    uint
}

// type vote struct {
// 	elector types.PlayerId
// 	target  types.PlayerId
// 	weight  uint
// }

// type electionResult struct {
// 	loser types.PlayerId
// 	stats pollStatistics
// }

func NewPoll(electors []types.PlayerId) *Poll {
	return &Poll{
		electors:      electors,
		votedElectors: make([]types.PlayerId, len(electors)),
		records:       make(map[uint]map[types.PlayerId]*record),
	}
}

func (p *Poll) IsOpen() bool {
	return p.isVoting
}

func (p *Poll) IsAllowed(playerId types.PlayerId) bool {
	return slices.Contains(p.electors, playerId) &&
		!slices.Contains(p.votedElectors, playerId)
}

func (p *Poll) GetRecord() map[types.PlayerId]*record {
	return p.records[p.recordId]
}

func (p *Poll) Open() bool {
	if p.IsOpen() {
		return false
	}

	p.isVoting = true
	p.recordId += 1
	p.records[p.recordId] = make(map[types.PlayerId]*record)
	p.votedElectors = make([]types.PlayerId, len(p.electors))

	return true
}

func (p *Poll) Close() map[types.PlayerId]*record {
	if !p.IsOpen() {
		return p.records[p.recordId]
	}

	p.isVoting = false
	p.records[p.recordId][types.UnknownPlayer] = &record{}

	for _, elector := range p.electors {
		if !slices.Contains(p.votedElectors, elector) {
			p.records[p.recordId][types.UnknownPlayer].votes += 1
			p.records[p.recordId][types.UnknownPlayer].electors = append(
				p.records[p.recordId][types.UnknownPlayer].electors,
				elector,
			)
		}
	}

	return p.records[p.recordId]
}

func (p *Poll) Vote(playerId types.PlayerId, target types.PlayerId) {
	if !p.IsAllowed(playerId) {
		return
	}

	if p.records[p.recordId][target] == nil {
		p.records[p.recordId][target] = &record{}
	}

	p.records[p.recordId][target].votes += 1
	p.records[p.recordId][target].electors = append(
		p.records[p.recordId][target].electors,
		playerId,
	)

	p.votedElectors = append(p.votedElectors, playerId)
}

func (p *Poll) RemoveElector(playerId types.PlayerId) bool {
	if i := slices.Index(p.electors, playerId); i == -1 {
		return false
	} else {
		p.electors = slices.Delete(p.electors, i, i+1)

		return true
	}
}
