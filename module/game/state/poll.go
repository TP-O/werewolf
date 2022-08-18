package state

import (
	"uwwolf/types"

	"golang.org/x/exp/slices"
)

type pollStatistics map[int]*candidateStatistics

type vote struct {
	elector types.PlayerId
	target  types.PlayerId
	value   uint
}

type candidateStatistics struct {
	electors []types.PlayerId
	votes    uint
}

type electionResult struct {
	loser types.PlayerId
	stats pollStatistics
}

type Poll struct {
	stageId       int
	isVoting      bool
	electors      []types.PlayerId
	votedElectors []types.PlayerId
	history       map[int]pollStatistics
}

func NewPoll(electors []types.PlayerId) *Poll {
	return &Poll{
		electors:      electors,
		votedElectors: make([]types.PlayerId, 0),
		history:       make(map[int]pollStatistics, 0),
	}
}

func (p *Poll) IsVoting() bool {
	return p.isVoting
}

func (p *Poll) IsVoted(playerId types.PlayerId) bool {
	//
}

func (p *Poll) History() map[int]pollStatistics {
	return p.history
}

func (p *Poll) Start() bool {
	if p.isVoting {
		return false
	}

	p.isVoting = true
	p.history[p.stageId] = make(pollStatistics)

	return true
}

func (p *Poll) Vote(elector types.PlayerId, target types.PlayerId) bool {
	if !p.IsVoting() ||
		len(p.votedElectors) >= len(p.electors) ||
		slices.Contains(p.votedElectors, elector) {

		return false
	}

	p.history[p.stageId][int(target)].votes += value
	p.history[p.stageId][int(target)].electors = append(
		p.history[p.stageId][int(target)].electors,
		elector,
	)
	p.votedElectors = append(p.votedElectors, elector)

	return true
}

func (p *Poll) Clear() {
	p.isVoting = false
	p.votedElectors = make([]types.PlayerId, 0)
	p.history[p.stageId] = make(pollStatistics)
}

func (p *Poll) Close() *electionResult {
	if !p.IsVoting() {
		return nil
	}

	result := electionResult{
		stats: make(pollStatistics),
		loser: p.eliminateLoser(p.history[p.stageId]),
	}

	p.stageId++
	p.isVoting = false

	return &result
}

func (p *Poll) eliminateLoser(stats pollStatistics) types.PlayerId {
	loser := 0
	totalVotes := uint(0)

	for target, cdStats := range stats {
		if cdStats.votes > stats[loser].votes {
			loser = target
		}

		totalVotes += cdStats.votes
	}

	if stats[loser].votes > totalVotes/2 {
		return types.PlayerId(loser)
	}

	return 0
}
