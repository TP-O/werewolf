package stuff

import (
	"fmt"
	"strconv"
	"time"
	"uwwolf/contract/itf"
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
	game     itf.IGame
	pub      chan string
}

func (p *Poll) Init(game itf.IGame, timeout time.Duration) {
	p.game = game
	p.timeout = timeout
	p.result = make(map[uint]uint)

	game.Pipe(&p.pub)
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
		p.pub <- strconv.Itoa(int(vote.elector)) + " voted " + strconv.Itoa(int(vote.target))
	}

	p.isVoting = false

	p.game.NextTurn()

	fmt.Println("End voting!")
}
