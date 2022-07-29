package state

import (
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"uwwolf/module/game/role"
	"uwwolf/types"
)

type turn struct {
	players []types.PlayerId
	timeout time.Duration
	role    role.Role
}

type Phase struct {
	currentPhase types.Phase
	currentTurn  int
	next         chan bool
	waitNext     sync.WaitGroup
	phases       map[types.Phase][]*turn
}

func NewPhase() Phase {
	return Phase{
		currentPhase: types.NightPhase,
		phases:       make(map[types.Phase][]*turn),
	}
}

func (p *Phase) Export() map[types.Phase][]*turn {
	return p.phases
}

func (p *Phase) Start() {
	if p.IsEmpty() {
		panic("Phase is empty!")
	}

	go func() {
		for {
			// Wait until time out or turn over
			select {
			case <-p.next:
				p.nextTurn()

				p.waitNext.Done()
			case <-time.After(p.GetTurn().timeout):
				p.nextTurn()
			}
		}
	}()
}

func (p *Phase) nextTurn() bool {
	if p.IsEmpty() {
		return false
	}

	if int(p.currentTurn) < len(p.phases[p.currentPhase])-1 {
		p.currentTurn++
	} else {
		p.currentTurn = 0
		p.currentPhase = (p.currentPhase + 1) % (types.NightPhase + 1)

		if p.currentPhase == 0 {
			p.currentPhase = 1
		}
	}

	if len(p.phases[p.currentPhase]) == 0 {
		return p.nextTurn()
	}

	return true
}

func (p *Phase) IsEmpty() bool {
	for _, p := range p.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

func (p *Phase) NextTurn() {
	p.waitNext.Add(1)

	p.next <- true

	p.waitNext.Wait()
}

func (p *Phase) GetTurn() *turn {
	return p.phases[p.currentPhase][p.currentTurn]
}

func (p *Phase) AddTurn(phase types.Phase, timeout time.Duration, role role.Role, players []types.PlayerId) bool {
	if phase > types.NightPhase {
		return false
	}

	p.phases[phase] = append(p.phases[phase], &turn{
		role:    role,
		timeout: timeout,
		players: players,
	})

	return true
}

func (p *Phase) RemoveTurn(phase types.Phase, roleName string) bool {
	for i := 0; i < len(p.phases[phase]); i++ {
		if p.phases[phase][i].role.GetName() == roleName {
			slices.Delete(p.phases[phase], i, i+1)

			return true
		}
	}

	return false
}

func (p *Phase) AddPlayer(phase types.Phase, turnIndex int, players ...types.PlayerId) bool {
	if phase > types.NightPhase || p.phases[phase][turnIndex] == nil {
		return false
	}

	p.phases[phase][turnIndex].players = append(p.phases[phase][turnIndex].players, players...)

	return true
}

func (p *Phase) RemovePlayer(player types.PlayerId) {
	for _, phase := range p.phases {
		for _, turn := range phase {
			deletedIndex := slices.Index(turn.players, player)

			if deletedIndex != -1 {
				turn.players = slices.Delete(turn.players, deletedIndex, deletedIndex+1)
			}
		}
	}
}

func (p *Phase) IsValidPlayer(player types.PlayerId) bool {
	return slices.Contains(p.GetTurn().players, player)
}

func (p *Phase) UseSkill(data *types.ActionData) bool {
	return p.GetTurn().role.UseSkill(data)
}
