package action

import (
	"fmt"
	"time"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/util"
)

const VoteAction = "Vote"

type voteState struct {
	counter  uint
	isVoting bool
	poll     chan string
	result   map[string]uint
}

func NewVoteAction() itf.IAction {
	return &action[voteState]{
		name: VoteAction,
		state: voteState{
			poll:   make(chan string),
			result: make(map[string]uint),
		},
		kit: actionKit[voteState]{
			validate: validateVote,
			execute:  executeVote,
			skip:     skipVote,
		},
	}
}

func validateVote(instruction *typ.ActionInstruction) bool {
	return instruction.Skipped ||
		(!instruction.Skipped && len(instruction.Targets) == 1)
}

func executeVote(game itf.IGame, instruction *typ.ActionInstruction, state *voteState) bool {
	// Open election
	if !state.isVoting {
		go startVote(game, state)
	}

	if util.IsClosed(state.poll) {
		return false
	}

	state.counter++

	state.poll <- instruction.Targets[0]

	fmt.Println(instruction.Actor + " voted " + instruction.Targets[0])

	return true
}

func skipVote(game itf.IGame, instruction *typ.ActionInstruction, state *voteState) bool {
	// Open election
	if !state.isVoting {
		startVote(game, state)
	}

	state.counter++

	fmt.Println(instruction.Actor + " skipped")

	return true
}

func resetVote(state *voteState) {
	state.poll = make(chan string)

	state.isVoting = false
	state.poll = make(chan string)
	state.result = make(map[string]uint)
}

func startVote(game itf.IGame, state *voteState) {
	state.isVoting = true

	time.AfterFunc(2*time.Second, func() { close(state.poll) })

	go handleVote(game, state)
}

func handleVote(game itf.IGame, state *voteState) {
	for socketId := range state.poll {
		state.result[socketId]++

		fmt.Println("Voted: ", socketId)
	}

	resetVote(state)

	game.NextTurn()
}
