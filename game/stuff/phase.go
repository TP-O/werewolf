package stuff

import (
	"time"

	"golang.org/x/exp/slices"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/enum"
)

type turn struct {
	playerIds []uint
	timeout   time.Duration
	role      itf.IRole
}

type Phase struct {
	currentPhaseId   uint8
	currentTurnIndex uint8
	nextTurnSignal   chan bool
	phases           map[uint8][]*turn
}

func (p *Phase) Data() map[uint8][]*turn {
	return p.phases
}

func (p *Phase) Init() {
	p.currentPhaseId = enum.NightPhase
	p.phases = make(map[uint8][]*turn)
}

func (p *Phase) Start() {
	for {
		// Wait until time out or turn over
		select {
		case <-p.nextTurnSignal:
		case <-time.After(p.GetTurn().timeout):
		}

		p.nextTurn()
	}
}

func (p *Phase) NextTurn() {
	p.nextTurnSignal <- true
}

func (p *Phase) AddTurn(phaseId uint8, role itf.IRole, playerIds []uint) bool {
	if phaseId >= enum.EndPhase {
		return false
	}

	p.phases[phaseId] = append(p.phases[phaseId], &turn{
		role:      role,
		playerIds: playerIds,
	})

	return true
}

func (p *Phase) RemoveTurn(phaseId uint8, roleName string) bool {
	for i := 0; i < len(p.phases[phaseId]); i++ {
		if p.phases[phaseId][i].role.GetName() == roleName {
			slices.Delete(p.phases[phaseId], i, i+1)

			return true
		}
	}

	return false
}

func (p *Phase) GetTurn() *turn {
	return p.phases[p.currentPhaseId][p.currentTurnIndex]
}

func (p *Phase) AddPlayer(phaseId uint8, turnIndex uint8, playerIds ...uint) bool {
	if phaseId >= enum.EndPhase || p.phases[phaseId][turnIndex] == nil {
		return false
	}

	p.phases[phaseId][turnIndex].playerIds = append(p.phases[phaseId][turnIndex].playerIds, playerIds...)

	return true
}

func (p *Phase) RemovePlayer(playerId uint) {
	for _, phase := range p.phases {
		for _, turn := range phase {
			deletedIndex := slices.Index(turn.playerIds, playerId)

			if deletedIndex != -1 {
				turn.playerIds = slices.Delete(turn.playerIds, deletedIndex, deletedIndex+1)
			}
		}
	}
}

func (p *Phase) IsValidPlayer(playerId uint) bool {
	return slices.Contains(p.GetTurn().playerIds, playerId)
}

func (p *Phase) UseSkill(instruction *typ.ActionInstruction) bool {
	return p.GetTurn().role.UseSkill(instruction)
}

func (p *Phase) nextTurn() *turn {
	if int(p.currentTurnIndex) < len(p.phases[p.currentPhaseId])-1 {
		p.currentTurnIndex++
	} else {
		p.currentTurnIndex = 0
		p.currentPhaseId = (p.currentPhaseId + 1) % enum.EndPhase

		if p.currentPhaseId == 0 {
			p.currentPhaseId = 1
		}
	}

	if p.GetTurn() == nil {
		return p.nextTurn()
	}

	return p.GetTurn()
}
