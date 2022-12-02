package role

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRoleID(t *testing.T) {
	id := config.HunterRoleID
	role := role{
		id: id,
	}

	assert.Equal(t, id, role.ID())
}

func TestRolePhaseID(t *testing.T) {
	phaseID := config.DuskPhaseID
	role := role{
		phaseID: phaseID,
	}

	assert.Equal(t, phaseID, role.PhaseID())
}

func TestRoleFactionID(t *testing.T) {
	factionID := config.WerewolfFactionID
	role := role{
		factionID: factionID,
	}

	assert.Equal(t, factionID, role.FactionID())
}

func TestRolePriority(t *testing.T) {
	priority := config.VillagerTurnPriority
	role := role{
		priority: priority,
	}

	assert.Equal(t, priority, role.Priority())
}

func TestRoleBeginRound(t *testing.T) {
	round := types.Round(9)
	role := role{
		beginRound: round,
	}

	assert.Equal(t, round, role.BeginRound())
}

func TestRoleActiveLimit(t *testing.T) {
	actionID := config.KillActionID
	expectedActionLimit := types.Limit(2)
	role := role{
		abilities: map[types.ActionID]*ability{
			actionID: {
				action:      nil,
				activeLimit: expectedActionLimit,
			},
		},
	}

	assert.Equal(t, expectedActionLimit, role.ActiveLimit(actionID))
	assert.Equal(t, config.ReachedLimit, role.ActiveLimit(types.ActionID(99)))
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

	playerID := types.PlayerID("1")
	actionID1 := types.ActionID(1)
	actionID2 := types.ActionID(2)

	mockAction1.EXPECT().ID().Return(actionID1).AnyTimes()
	mockAction2.EXPECT().ID().Return(actionID2).AnyTimes()

	tests := []struct {
		name                string
		req                 *types.UseRoleRequest
		res                 *types.ActionResponse
		expectedActiveLimit types.Limit
		setup               func(role)
	}{
		{
			name: "Invalid ability",
			req: &types.UseRoleRequest{
				ActionID: types.ActionID(99),
			},
			res: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "This is beyond your ability (╥﹏╥)",
			},
			expectedActiveLimit: -99,
			setup:               func(role role) {},
		},
		{
			name: "Ability is out of use",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			res: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "Unable to use this ability anymore ¯\\(º_o)/¯",
			},
			expectedActiveLimit: 0,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = 0
			},
		},
		{
			name: "Ok but don't change limit (Skip action)",
			req: &types.UseRoleRequest{
				ActionID:  actionID1,
				IsSkipped: true,
			},
			res: &types.ActionResponse{
				Ok:        true,
				IsSkipped: true,
			},
			expectedActiveLimit: 5,
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
			name: "Ok but don't change limit (Action execution is failed)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			res: &types.ActionResponse{
				Ok: false,
			},
			expectedActiveLimit: 5,
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
			name: "Ok but don't change limit (Unlimited ability)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			res: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			expectedActiveLimit: config.Unlimited,
			setup: func(role role) {
				role.abilities[actionID1].activeLimit = config.Unlimited
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
			name: "Ok but reduce limit by 1 (action execution is successful)",
			req: &types.UseRoleRequest{
				ActionID: actionID1,
			},
			res: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			expectedActiveLimit: 4,
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
				abilities: map[types.ActionID]*ability{
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

			assert.Equal(t, test.res, res)

			// Invalid ability
			if test.expectedActiveLimit == -99 {
				assert.Nil(t, role.abilities[test.req.ActionID])
			} else {
				assert.Equal(t, test.expectedActiveLimit, role.abilities[test.req.ActionID].activeLimit)
			}
		})
	}
}
