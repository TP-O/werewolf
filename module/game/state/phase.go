package state

import (
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"uwwolf/module/game/role"
	"uwwolf/types"
)

type turn struct {
	playerIds []int
	timeout   time.Duration
	role      role.Role
}

type phase struct {
	currentPhase types.Phase
	currentTurn  int
	next         chan bool
	waitNext     sync.WaitGroup
	phases       map[types.Phase][]*turn
}

func NewPhase() *phase {
	return &phase{
		currentPhase: types.NightPhase,
		phases:       make(map[types.Phase][]*turn),
	}
}

func (p *phase) Export() map[types.Phase][]*turn {
	return p.phases
}

func (p *phase) Start() {
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

func (p *phase) nextTurn() bool {
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

func (p *phase) IsEmpty() bool {
	for _, p := range p.phases {
		if len(p) != 0 {
			return false
		}
	}

	return true
}

func (p *phase) NextTurn() {
	p.waitNext.Add(1)

	p.next <- true

	p.waitNext.Wait()
}

func (p *phase) GetTurn() *turn {
	return p.phases[p.currentPhase][p.currentTurn]
}

func (p *phase) AddTurn(phase types.Phase, timeout time.Duration, role role.Role, playerIds []int) bool {
	if phase > types.NightPhase {
		return false
	}

	p.phases[phase] = append(p.phases[phase], &turn{
		role:      role,
		timeout:   timeout,
		playerIds: playerIds,
	})

	return true
}

func (p *phase) RemoveTurn(phase types.Phase, roleName string) bool {
	for i := 0; i < len(p.phases[phase]); i++ {
		if p.phases[phase][i].role.GetName() == roleName {
			slices.Delete(p.phases[phase], i, i+1)

			return true
		}
	}

	return false
}

func (p *phase) AddPlayer(phase types.Phase, turnIndex int, playerIds ...int) bool {
	if phase > types.NightPhase || p.phases[phase][turnIndex] == nil {
		return false
	}

	p.phases[phase][turnIndex].playerIds = append(p.phases[phase][turnIndex].playerIds, playerIds...)

	return true
}

func (p *phase) RemovePlayer(playerId int) {
	for _, phase := range p.phases {
		for _, turn := range phase {
			deletedIndex := slices.Index(turn.playerIds, playerId)

			if deletedIndex != -1 {
				turn.playerIds = slices.Delete(turn.playerIds, deletedIndex, deletedIndex+1)
			}
		}
	}
}

func (p *phase) IsValidPlayer(playerId int) bool {
	return slices.Contains(p.GetTurn().playerIds, playerId)
}

func (p *phase) UseSkill(data *types.ActionData) bool {
	return p.GetTurn().role.UseSkill(data)
}
