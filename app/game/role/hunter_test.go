package role_test

import (
	"testing"
	"uwwolf/app/game/role"
	"uwwolf/app/game/state"
	"uwwolf/app/types"
	"uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHunterAfterDeath(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockPlayer := game.NewMockPlayer(ctrl)

	//=============================================================
	round := state.NewRound()
	round.AddTurn(&types.TurnSetting{
		PhaseId: types.DayPhase,
		RoleId:  types.VillagerRole,
	})
	round.AddTurn(&types.TurnSetting{
		PhaseId: types.NightPhase,
		RoleId:  types.WerewolfRole,
	})

	mockGame.
		EXPECT().
		Player(gomock.Any()).
		Return(mockPlayer)
	mockGame.
		EXPECT().
		Round().
		Return(round).
		Times(4)

	mockPlayer.
		EXPECT().
		Id().
		Return(types.PlayerId("1")).
		Times(2)

	roleSetting := &types.RoleSetting{
		Id:         types.HunterRole,
		PhaseId:    types.DayPhase,
		Expiration: types.OneTimes,
	}
	h := role.NewHunterRole(mockGame, roleSetting)

	//=============================================================
	// Die at other phases (current is night)
	h.AfterDeath()

	round.NextTurn() // hunter turn

	assert.Equal(t, roleSetting.Id, round.CurrentTurn().RoleId())

	//=============================================================
	// Die at his phase
	round.NextTurn() // villager turn

	h.AfterDeath()

	round.NextTurn()

	assert.Equal(t, roleSetting.Id, round.CurrentTurn().RoleId())
}
