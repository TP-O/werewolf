package stuff

import (
	"fmt"
	"time"
	"uwwolf/util"
)

type vote struct {
	elector uint
	target  uint
}

type Poll struct {
	isVoting bool
	box      chan *vote
	timeout  time.Duration
	result   map[uint]uint
}

func (p *Poll) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

func (p *Poll) Start() bool {
	if p.isVoting {
		return false
	}

	p.isVoting = true
	p.box = make(chan *vote)

	time.AfterFunc(p.timeout, func() {
		close(p.box)
	})

	go p.handleVotes()

	return true
}

func (p *Poll) Vote(elector uint, target uint) {
	if !util.IsClosed(p.box) {
		p.box <- &vote{
			elector: elector,
			target:  target,
		}
	}
}

func (p *Poll) IsVoting() bool {
	return p.isVoting
}

func (p *Poll) handleVotes() {
	for vote := range p.box {
		fmt.Println(vote.elector, " voted ", vote.target)
	}

	p.isVoting = false

	fmt.Println("End voting!")
}
