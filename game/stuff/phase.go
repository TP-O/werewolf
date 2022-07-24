package stuff

import (
	"golang.org/x/exp/slices"

	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/enum"
)

type turn struct {
	playerIds []uint
	role      itf.IRole
}

type Phase struct {
	currentPhaseIndex uint8
	currentTurnIndex  uint8
	phases            map[uint8][]*turn
}

func (p *Phase) Init() {
	p.currentPhaseIndex = enum.NightPhase
	p.phases = make(map[uint8][]*turn)
}

func (p *Phase) AddTurn(phase uint8, role itf.IRole, playerIds []uint) {
	p.phases[phase] = append(p.phases[phase], &turn{
		role:      role,
		playerIds: playerIds,
	})
}

func (p *Phase) RemoveTurn(phase uint8, roleName string) bool {
	for i := 0; i < len(p.phases[phase]); i++ {
		if p.phases[phase][i].role.GetName() == roleName {
			slices.Delete(p.phases[phase], i, i+1)

			return true
		}
	}

	return false
}

func (p *Phase) GetTurn() *turn {
	return p.phases[p.currentPhaseIndex][p.currentTurnIndex]
}

func (p *Phase) NextTurn() *turn {
	if int(p.currentTurnIndex) < len(p.phases[p.currentPhaseIndex]) {
		p.currentTurnIndex++
	} else {
		p.currentTurnIndex = 0
		p.currentPhaseIndex = (p.currentPhaseIndex + 1) % enum.EndPhase

		if p.currentPhaseIndex == 0 {
			p.currentPhaseIndex = 1
		}
	}

	if p.GetTurn() == nil {
		return p.NextTurn()
	}

	return p.GetTurn()
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
