package role

import (
	"testing"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRoleID(t *testing.T) {
	id := enum.HunterRoleID
	role := role{
		id: id,
	}

	assert.Equal(t, id, role.ID())
}

func TestRolePhaseID(t *testing.T) {
	phaseID := enum.DuskPhaseID
	role := role{
		phaseID: phaseID,
	}

	assert.Equal(t, phaseID, role.PhaseID())
}

func TestRoleFactionID(t *testing.T) {
	factionID := enum.WerewolfFactionID
	role := role{
		factionID: factionID,
	}

	assert.Equal(t, factionID, role.FactionID())
}

func TestRolePriority(t *testing.T) {
	priority := enum.VillagerTurnPriority
	role := role{
		priority: priority,
	}

	assert.Equal(t, priority, role.Priority())
}

func TestRoleBeginRound(t *testing.T) {
	round := enum.Round(9)
	role := role{
		beginRound: round,
	}

	assert.Equal(t, round, role.BeginRound())
}

func TestRoleActiveLimit(t *testing.T) {
	actionID := enum.KillActionID
	expectedActionLimit := enum.Limit(2)
	role := role{
		abilities: map[enum.ActionID]*ability{
			actionID: {
				action:      nil,
				activeLimit: expectedActionLimit,
			},
		},
	}

	assert.Equal(t, expectedActionLimit, role.ActiveLimit(actionID))
	assert.Equal(t, enum.ReachedLimit, role.ActiveLimit(enum.ActionID(99)))
}

func TestRoleBeforeDeath(t *testing.T) {
	role := role{}

	assert.True(t, role.BeforeDeath())
}

func TestRoleAfterDeath(t *testing.T) {
	//
}

func TestRoleUseAbility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPlayer := gamemock.NewMockPlayer(ctrl)
	mockAction1 := gamemock.NewMockAction(ctrl)
	mockAction2 := gamemock.NewMockAction(ctrl)

	playerID := enum.PlayerID("1")
	actionID1 := enum.ActionID(1)
	actionID2 := enum.ActionID(2)

	mockAction1.EXPECT().ID().Return(actionID1).AnyTimes()
	mockAction2.EXPECT().ID().Return(actionID2).AnyTimes()

	tests := []struct {
		name           string
		req            *types.UseRoleRequest
		expectedRes    *types.ActionResponse
		isInvalidRole  bool
		newActiveLimit enum.Limit
		setup          func(role)
	}{
		{
			name: "Failure (Invalid ability)",
			req: &types.UseRoleRequest{
				ActionID: 99,
			},
			expectedRes: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "This is beyond your ability (╥﹏╥)",
			},
			isInvalidRole: true,
			setup:         func(role role) {},
		},
		{
			name: "Failure (Ability is out of use)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "Unable to use this ability anymore ¯\\(º_o)/¯",
			},
			newActiveLimit: 0,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = 0
			},
		},
		{
			name: "Ok (Skip action -> Doesn't change active limit)",
			req: &types.UseRoleRequest{
				ActionID:  actionID1,
				IsSkipped: true,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: true,
			},
			newActiveLimit: 5,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = 5
				mockAction1.EXPECT().Execute(&types.ActionRequest{
					ActorID:   playerID,
					IsSkipped: true,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: true,
					}).
					Times(1)
				mockPlayer.EXPECT().ID().Return(playerID).Times(1)
			},
		},
		{
			name: "Ok (Action execution is failed -> Doesn't change active limit)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok: false,
			},
			newActiveLimit: 5,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = 5
				mockAction1.EXPECT().Execute(&types.ActionRequest{
					ActorID: playerID,
				}).
					Return(&types.ActionResponse{
						Ok: false,
					}).
					Times(1)
				mockPlayer.EXPECT().ID().Return(playerID).Times(1)
			},
		},
		{
			name: "Ok (Use unlimited ability -> Doesn't change active limit)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			newActiveLimit: enum.Unlimited,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = enum.Unlimited
				mockAction1.EXPECT().Execute(&types.ActionRequest{
					ActorID: playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)
				mockPlayer.EXPECT().ID().Return(playerID).Times(1)
			},
		},
		{
			name: "Ok (Action execution is successful -> reduce active limit by 1)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			newActiveLimit: 4,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = 5
				mockAction1.EXPECT().Execute(&types.ActionRequest{
					ActorID: playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)
				mockPlayer.EXPECT().ID().Return(playerID).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			role := role{
				player: mockPlayer,
				abilities: map[enum.ActionID]*ability{
					actionID1: {
						action: mockAction1,
					},
					actionID2: {
						action: mockAction2,
					},
				},
			}
			test.setup(role)
			res := role.UseAbility(test.req)

			assert.Equal(t, test.expectedRes, res)

			// Invalid ability
			if test.isInvalidRole {
				assert.Nil(t, role.abilities[test.req.ActionID])
			} else {
				assert.Equal(t, test.newActiveLimit, role.abilities[test.req.ActionID].activeLimit)
			}
		})
	}
}
