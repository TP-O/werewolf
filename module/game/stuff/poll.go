package stuff

import (
	"fmt"
	"strconv"
	"time"
	"uwwolf/contract/itf"
	"uwwolf/enum"
	"uwwolf/util"
)

type vote struct {
	elector int
	target  int
}

type Poll struct {
	isVoting bool
	box      chan *vote
	timeout  time.Duration
	result   map[int]int
	game     itf.IGame
	pub      chan string
	capacity int
	total    int
}

func (p *Poll) Init(game itf.IGame, factionId int, timeout time.Duration) {
	p.game = game
	p.timeout = timeout
	p.result = make(map[int]int)

	game.Pipe(&p.pub)

	if factionId == enum.VillageFaction {
		p.capacity = p.game.NumberOfVillagers()
	} else if factionId == enum.WerewolfFaction {
		p.capacity = p.game.NumberOfWerewolves()
	}
}

func (p *Poll) Start() bool {
	if p.isVoting {
		return false
	}

	p.isVoting = true
	p.total = 0
	p.box = make(chan *vote)

	time.AfterFunc(p.timeout, func() {
		if util.IsChannelNotClosed(p.box) {
			close(p.box)
		}
	})

	go p.handleVotes()

	return true
}

func (p *Poll) Vote(elector int, target int) {
	if util.IsChannelNotClosed(p.box) {
		p.total++

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

		if p.total == p.capacity {
			close(p.box)
		}
	}

	p.isVoting = false

	p.game.NextTurn()

	fmt.Println("End voting!")
}
