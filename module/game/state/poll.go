package state

import (
	"golang.org/x/exp/slices"
)

type pollStatistics map[int]*candidateStatistics

type vote struct {
	elector int
	target  int
	value   uint
}

type candidateStatistics struct {
	electors []int
	votes    uint
}

type electionResult struct {
	loser int
	stats pollStatistics
}

type poll struct {
	stageId       int
	isVoting      bool
	electors      []int
	votedElectors []int
	history       map[int]pollStatistics
}

func NewPoll(electors []int) *poll {
	return &poll{
		electors:      electors,
		votedElectors: make([]int, 0),
		history:       make(map[int]pollStatistics, 0),
	}
}

func (p *poll) IsVoting() bool {
	return p.isVoting
}

func (p *poll) History() map[int]pollStatistics {
	return p.history
}

func (p *poll) Start() bool {
	if p.isVoting {
		return false
	}

	p.isVoting = true
	p.history[p.stageId] = make(pollStatistics)

	return true
}

func (p *poll) Vote(elector int, target int, value uint) bool {
	if !p.IsVoting() ||
		len(p.votedElectors) >= len(p.electors) ||
		slices.Contains(p.votedElectors, elector) {

		return false
	}

	p.history[p.stageId][target].votes += value
	p.history[p.stageId][target].electors = append(p.history[p.stageId][target].electors, elector)
	p.votedElectors = append(p.votedElectors, elector)

	return true
}

func (p *poll) Clear() {
	p.isVoting = false
	p.votedElectors = make([]int, 0)
	p.history[p.stageId] = make(pollStatistics)
}

func (p *poll) Close() *electionResult {
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

func (p *poll) eliminateLoser(stats pollStatistics) int {
	loser := 0
	totalVotes := uint(0)

	for target, cdStats := range stats {
		if cdStats.votes > stats[loser].votes {
			loser = target
		}

		totalVotes += cdStats.votes
	}

	if stats[loser].votes > totalVotes/2 {
		return loser
	}

	return 0
}
