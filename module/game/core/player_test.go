package core_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"uwwolf/mock/game"
	"uwwolf/module/game/core"
	"uwwolf/module/game/role"
	"uwwolf/module/game/state"
	"uwwolf/types"
)

func TestPlayerId(t *testing.T) {
	id := types.PlayerId("1")
	p := core.NewPlayer(nil, id)

	assert.Equal(t, id, p.Id())
}

func TestPlayerRoleIds(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)

	//=============================================================
	mockGame.
		EXPECT().
		Player(gomock.Any()).
		Return(nil).
		Times(2)

	r1 := role.NewSeerRole(mockGame, &types.RoleSetting{
		Id: 98,
	})
	r2 := role.NewSeerRole(mockGame, &types.RoleSetting{
		Id: 99,
	})
	p := core.NewPlayer(nil, types.PlayerId("1"))

	p.AssignRoles(r1, r2)

	assert.Contains(t, p.RoleIds(), r1.Id(), r2.Id())
}

func TestPlayerFactionId(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRole := game.NewMockRole(ctrl)

	//=============================================================
	factionId := types.FactionId(99)
	mockRole.
		EXPECT().
		Id().
		Return(types.RoleId(1)).
		Times(2)
	mockRole.
		EXPECT().
		FactionId().
		Return(factionId)

	p := core.NewPlayer(nil, types.PlayerId("1"))

	p.AssignRoles(mockRole)

	assert.Equal(t, factionId, p.FactionId())
}

func TestPlayerAssignRoles(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRole1 := game.NewMockRole(ctrl)
	mockRole2 := game.NewMockRole(ctrl)

	//=============================================================
	role1Id := types.RoleId(98)
	role1FactionId := types.FactionId(98)
	mockRole1.
		EXPECT().
		Id().
		Return(role1Id).
		Times(7)
	mockRole1.
		EXPECT().
		FactionId().
		Return(role1FactionId).
		Times(3)

	role2Id := types.RoleId(99)
	role2FactionId := types.FactionId(99)
	mockRole2.
		EXPECT().
		Id().
		Return(role2Id).
		Times(4)
	mockRole2.
		EXPECT().
		FactionId().
		Return(role2FactionId).
		Times(2)

	//=============================================================
	// Assign new roles with replacing faction id
	p := core.NewPlayer(nil, types.PlayerId("1"))
	p.AssignRoles(mockRole1, mockRole2)

	assert.Contains(t, p.RoleIds(), role1Id, role2Id)
	assert.Equal(t, role2FactionId, p.FactionId())

	//=============================================================
	// Asign new roles without replacing faction id
	p = core.NewPlayer(nil, types.PlayerId("1"))
	p.AssignRoles(mockRole2, mockRole1)

	assert.Contains(t, p.RoleIds(), role1Id, role2Id)
	assert.Equal(t, role2FactionId, p.FactionId())

	//=============================================================
	// Assign duplicate role
	p = core.NewPlayer(nil, types.PlayerId("1"))
	p.AssignRoles(mockRole1, mockRole1)

	assert.Len(t, p.RoleIds(), 1)
	assert.Contains(t, p.RoleIds(), role1Id)
	assert.Equal(t, role1FactionId, p.FactionId())
}

func TestUseSkill(t *testing.T) {
	//========================MOCK================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGame := game.NewMockGame(ctrl)
	mockRole1 := game.NewMockRole(ctrl)
	mockRole2 := game.NewMockRole(ctrl)

	//=============================================================
	round := state.NewRound()
	round.AddTurn(&types.TurnSetting{
		PhaseId:  types.DayPhase,
		RoleId:   types.VillagerRole,
		Position: types.NextPosition,
	})
	mockGame.
		EXPECT().
		Round().
		Return(round).
		Times(2)

	roleId := round.CurrentTurn().RoleId()
	mockRole1.
		EXPECT().
		Id().
		Return(roleId).
		Times(2)
	mockRole1.
		EXPECT().
		FactionId().
		Return(types.FactionId(1))
	mockRole1.
		EXPECT().
		ActivateSkill(gomock.Any()).
		Return(&types.ActionResponse{Ok: true})

	mockRole2.
		EXPECT().
		Id().
		Return(types.RoleId(99)).
		Times(2)
	mockRole2.
		EXPECT().
		FactionId().
		Return(types.FactionId(1))

	//=============================================================
	// Turn of player
	p := core.NewPlayer(mockGame, types.PlayerId("1"))
	p.AssignRoles(mockRole1)

	res := p.UseSkill(&types.ActionRequest{})

	assert.True(t, res.Ok)

	//=============================================================
	// Not turn of player
	p = core.NewPlayer(mockGame, types.PlayerId("1"))
	p.AssignRoles(mockRole2)

	res = p.UseSkill(&types.ActionRequest{})

	assert.False(t, res.Ok)
	assert.NotNil(t, res.Error)
}
