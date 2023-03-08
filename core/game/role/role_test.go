package role

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RoleSuite struct {
	suite.Suite
	playerID  types.PlayerID
	actionID1 types.ActionID
	actionID2 types.ActionID
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

func (rs *RoleSuite) SetupSuite() {
	rs.playerID = types.PlayerID("1")
	rs.actionID1 = types.ActionID(1)
	rs.actionID2 = types.ActionID(2)
}

func (rs RoleSuite) TestID() {
	id := vars.HunterRoleID
	role := role{
		id: id,
	}

	rs.Equal(id, role.ID())
}

// func (rs RoleSuite) TestPhaseID() {
// 	phaseID := vars.DuskPhaseID
// 	role := role{
// 		phaseID: phaseID,
// 	}

// 	rs.Equal(phaseID, role.PhaseID())
// }

func (rs RoleSuite) TestFactionID() {
	factionID := vars.WerewolfFactionID
	role := role{
		factionID: factionID,
	}

	rs.Equal(factionID, role.FactionID())
}

// func (rs RoleSuite) TestBeginRoundID() {
// 	round := types.Round(9)
// 	role := role{
// 		beginRound: round,
// 	}

// 	rs.Equal(round, role.BeginRoundID())
// }

func (rs RoleSuite) TestActiveTimes() {
	action1Limit := types.Times(1)
	action2Limit := types.Times(2)
	role := role{
		abilities: []*ability{
			{
				activeLimit: action1Limit,
			},
			{
				activeLimit: action2Limit,
			},
		},
	}

	rs.Equal(action1Limit, role.ActiveTimes(0))
	rs.Equal(action2Limit, role.ActiveTimes(1))
	rs.Equal(action1Limit+action2Limit, role.ActiveTimes(-1))
	rs.Equal(vars.OutOfTimes, role.ActiveTimes(99))
}

func (rs RoleSuite) TestBeforeDeath() {
	role := new(role)

	rs.True(role.BeforeDeath())
}

func (rs RoleSuite) TestAfterDeath() {
	//
}

func (rs RoleSuite) TestActivateAbility() {
	tests := []struct {
		name          string
		req           *types.ActivateAbilityRequest
		expectedRes   *types.ActionResponse
		expectedLimit types.Times
		setup         func(role, *gamemock.MockAction, *gamemock.MockScheduler)
	}{
		{
			name: "Failure (Invalid ability index)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 99,
			},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "This is beyond your ability (╥﹏╥)",
			},
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {},
		},
		{
			name: "Failure (Ability is out of use)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 0,
			},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Unable to use this ability anymore ¯\\(º_o)/¯",
			},
			expectedLimit: 0,
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {
				role.abilities[0].activeLimit = 0
			},
		},
		{
			name: "Ok (Skip action -> Doesn't change active limit)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 0,
				IsSkipped:    true,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: true,
			},
			expectedLimit: 2,
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {
				role.abilities[0].activeLimit = 2
				ma.EXPECT().Execute(&types.ActionRequest{
					ActorID:   rs.playerID,
					IsSkipped: true,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: true,
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Action execution is failed -> Doesn't change active limit)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 0,
			},
			expectedRes: &types.ActionResponse{
				Ok: false,
			},
			expectedLimit: 2,
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {
				role.abilities[0].activeLimit = 2
				ma.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok: false,
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Use unlimited ability -> Doesn't change active limit)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 0,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			expectedLimit: vars.UnlimitedTimes,
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {
				role.abilities[0].activeLimit = vars.UnlimitedTimes
				ma.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Action execution is successful -> Reduce active limit by 1)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 0,
			},
			expectedRes: &types.ActionResponse{
				Ok: true,
			},
			expectedLimit: 1,
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {
				role.abilities[0].activeLimit = 2
				ma.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)

			},
		},
		{
			name: "Ok (Action execution is successful -> Reach limit)",
			req: &types.ActivateAbilityRequest{
				AbilityIndex: 0,
			},
			expectedRes: &types.ActionResponse{
				Ok: true,
			},
			expectedLimit: vars.OutOfTimes,
			setup: func(role role, ma *gamemock.MockAction, ms *gamemock.MockScheduler) {
				role.abilities[0].activeLimit = 1
				ma.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)
				ms.EXPECT().RemoveSlot(&types.RemovedTurnSlot{
					PhaseID:  role.phaseID,
					RoleID:   role.id,
					PlayerID: role.player.ID(),
				})
			},
		},
	}

	for _, test := range tests {
		rs.Run(test.name, func() {
			ctrl := gomock.NewController(rs.T())
			defer ctrl.Finish()
			player := gamemock.NewMockPlayer(ctrl)
			scheduler := gamemock.NewMockScheduler(ctrl)
			game := gamemock.NewMockGame(ctrl)
			action := gamemock.NewMockAction(ctrl)

			action.EXPECT().ID().Return(rs.actionID1).AnyTimes()
			player.EXPECT().ID().Return(rs.playerID).AnyTimes()
			game.EXPECT().Scheduler().Return(scheduler).AnyTimes()

			role := role{
				game:   game,
				player: player,
				abilities: []*ability{
					{
						action: action,
					},
				},
			}
			test.setup(role, action, scheduler)
			res := role.ActivateAbility(test.req)

			rs.Equal(test.expectedRes, res)
			rs.Equal(test.expectedLimit, role.ActiveTimes(0))
		})
	}
}
