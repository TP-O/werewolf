package role

import (
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"
	mock_game_logic "uwwolf/test/mock/app/game/logic"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RoleSuite struct {
	suite.Suite
	playerId types.PlayerId
	actionId types.ActionId
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

func (rs *RoleSuite) SetupSuite() {
	rs.playerId = types.PlayerId("1")
	rs.actionId = types.ActionId(1)
}

func (rs RoleSuite) TestId() {
	Id := constants.HunterRoleId
	role := role{
		id: Id,
	}

	rs.Equal(Id, role.Id())
}

func (rs RoleSuite) TestFactionId() {
	factionId := constants.WerewolfFactionId
	role := role{
		factionId: factionId,
	}

	rs.Equal(factionId, role.FactionId())
}

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
	rs.Equal(constants.OutOfTimes, role.ActiveTimes(99))
}

func (rs RoleSuite) TestOnAfterAssign() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	scheduler := mock_game_logic.NewMockScheduler(ctrl)

	moderator.EXPECT().Scheduler().Return(scheduler)

	r := role{
		id:         constants.HunterRoleId,
		moderator:  moderator,
		playerId:   rs.playerId,
		beginRound: constants.SecondRound,
		phaseId:    constants.DayPhaseId,
		turn:       constants.PostTurn,
	}

	scheduler.EXPECT().AddSlot(types.NewTurnSlot{
		PhaseId:    r.phaseId,
		Turn:       r.turn,
		BeginRound: r.beginRound,
		PlayerId:   rs.playerId,
		RoleId:     r.Id(),
	})

	r.OnAfterAssign()
}

func (rs RoleSuite) TestOnAfterRevoke() {
	ctrl := gomock.NewController(rs.T())
	defer ctrl.Finish()
	moderator := mock_game_logic.NewMockModerator(ctrl)
	scheduler := mock_game_logic.NewMockScheduler(ctrl)

	moderator.EXPECT().Scheduler().Return(scheduler)

	r := role{
		id:        constants.HunterRoleId,
		phaseId:   constants.DayPhaseId,
		moderator: moderator,
		playerId:  rs.playerId,
	}

	scheduler.EXPECT().RemoveSlot(types.RemovedTurnSlot{
		PhaseId:  r.phaseId,
		PlayerId: rs.playerId,
		RoleId:   r.Id(),
	})

	r.OnAfterRevoke()
}

func (rs RoleSuite) TestOnBeforeDeath() {
	role := new(role)
	isDead := role.OnBeforeDeath()

	rs.True(isDead)
	rs.True(role.isBeforeDeathTriggered)
}

func (rs RoleSuite) TestAfterDeath() {
	//
}

func (rs RoleSuite) TestUse() {
	id := types.RoleId(1)
	round := constants.FirstRound
	phaseId := constants.DayPhaseId
	turn := constants.PostTurn

	tests := []struct {
		name          string
		req           types.RoleRequest
		expectedRes   types.RoleResponse
		expectedLimit types.Times
		setup         func(
			role,
			*mock_game_logic.MockAction,
			*mock_game_logic.MockScheduler,
			*mock_game_logic.MockModerator,
		)
	}{
		{
			name: "Failure (InvalId ability index)",
			req: types.RoleRequest{
				AbilityIndex: 99,
			},
			expectedRes: types.RoleResponse{
				ActionResponse: types.ActionResponse{
					Message: "This action is beyond your ability (╥﹏╥)",
				},
			},
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
			},
		},
		{
			name: "Failure (Ability is out of use)",
			req: types.RoleRequest{
				AbilityIndex: 0,
			},
			expectedRes: types.RoleResponse{
				ActionResponse: types.ActionResponse{
					Message: "Unable to use this action anymore ¯\\(º_o)/¯",
				},
			},
			expectedLimit: constants.OutOfTimes,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.OutOfTimes
			},
		},
		{
			name: "Ok (Skip action -> Doesn't change active limit)",
			req: types.RoleRequest{
				AbilityIndex: 0,
				IsSkipped:    true,
			},
			expectedRes: types.RoleResponse{
				Round:   round,
				PhaseId: phaseId,
				Turn:    turn,
				RoleId:  id,
				ActionResponse: types.ActionResponse{
					Ok: true,
					ActionRequest: types.ActionRequest{
						IsSkipped: true,
					},
				},
			},
			expectedLimit: constants.Twice,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.Twice

				ms.EXPECT().Round().Return(round).Times(1)
				ms.EXPECT().PhaseId().Return(phaseId).Times(1)
				ms.EXPECT().Turn().Return(turn).Times(1)
				ma.EXPECT().Execute(types.ActionRequest{
					ActorId:   role.playerId,
					IsSkipped: true,
				}).
					Return(types.ActionResponse{
						Ok: true,
						ActionRequest: types.ActionRequest{
							IsSkipped: true,
						},
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Action execution is failed -> Doesn't change active limit)",
			req: types.RoleRequest{
				AbilityIndex: 0,
			},
			expectedRes: types.RoleResponse{
				Round:   round,
				PhaseId: phaseId,
				Turn:    turn,
				RoleId:  id,
				ActionResponse: types.ActionResponse{
					Ok: false,
				},
			},
			expectedLimit: constants.Twice,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.Twice
				role.abilities[0].isImmediate = true

				ms.EXPECT().Round().Return(round).Times(1)
				ms.EXPECT().PhaseId().Return(phaseId).Times(1)
				ms.EXPECT().Turn().Return(turn).Times(1)
				ma.EXPECT().Execute(types.ActionRequest{
					ActorId: role.playerId,
				}).
					Return(types.ActionResponse{
						Ok: false,
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Use unlimited ability -> Doesn't change active limit)",
			req: types.RoleRequest{
				AbilityIndex: 0,
			},
			expectedRes: types.RoleResponse{
				Round:   round,
				PhaseId: phaseId,
				Turn:    turn,
				RoleId:  id,
				ActionResponse: types.ActionResponse{
					Ok: true,
				},
			},
			expectedLimit: constants.UnlimitedTimes,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.UnlimitedTimes
				role.abilities[0].isImmediate = true

				ms.EXPECT().Round().Return(round).Times(1)
				ms.EXPECT().PhaseId().Return(phaseId).Times(1)
				ms.EXPECT().Turn().Return(turn).Times(1)
				ma.EXPECT().Execute(types.ActionRequest{
					ActorId: role.playerId,
				}).
					Return(types.ActionResponse{
						Ok: true,
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Action execution is successful -> Reduce active limit by 1)",
			req: types.RoleRequest{
				AbilityIndex: 0,
			},
			expectedRes: types.RoleResponse{
				Round:   round,
				PhaseId: phaseId,
				Turn:    turn,
				RoleId:  id,
				ActionResponse: types.ActionResponse{
					Ok: true,
				},
			},
			expectedLimit: constants.Once,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.Twice
				role.abilities[0].isImmediate = true

				ms.EXPECT().Round().Return(round).Times(1)
				ms.EXPECT().PhaseId().Return(phaseId).Times(1)
				ms.EXPECT().Turn().Return(turn).Times(1)
				ma.EXPECT().Execute(types.ActionRequest{
					ActorId: role.playerId,
				}).
					Return(types.ActionResponse{
						Ok: true,
					}).
					Times(1)
			},
		},
		{
			name: "Ok (Action execution is successful -> Reach limit)",
			req: types.RoleRequest{
				AbilityIndex: 0,
			},
			expectedRes: types.RoleResponse{
				Round:   round,
				PhaseId: phaseId,
				Turn:    turn,
				RoleId:  id,
				ActionResponse: types.ActionResponse{
					Ok: true,
				},
			},
			expectedLimit: constants.OutOfTimes,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.Once
				role.abilities[0].isImmediate = true

				ms.EXPECT().Round().Return(round).Times(1)
				ms.EXPECT().PhaseId().Return(phaseId).Times(1)
				ms.EXPECT().Turn().Return(turn).Times(1)
				ma.EXPECT().Execute(types.ActionRequest{
					ActorId: role.playerId,
				}).
					Return(types.ActionResponse{
						Ok: true,
					}).
					Times(1)
				ms.EXPECT().RemoveSlot(types.RemovedTurnSlot{
					PhaseId:  role.phaseId,
					RoleId:   role.id,
					PlayerId: role.playerId,
				})
			},
		},
		{
			name: "Ok (Action will be executed later)",
			req: types.RoleRequest{
				AbilityIndex: 0,
			},
			expectedRes: types.RoleResponse{
				Round:   round,
				PhaseId: phaseId,
				Turn:    turn,
				RoleId:  id,
				ActionResponse: types.ActionResponse{
					Ok:       true,
					ActionId: rs.actionId,
					ActionRequest: types.ActionRequest{
						ActorId: rs.playerId,
					},
					Message: "Action is registered!",
				},
			},
			expectedLimit: constants.Once,
			setup: func(
				role role,
				ma *mock_game_logic.MockAction,
				ms *mock_game_logic.MockScheduler,
				mm *mock_game_logic.MockModerator,
			) {
				role.abilities[0].activeLimit = constants.Twice
				role.abilities[0].isImmediate = false

				ms.EXPECT().Round().Return(round).Times(1)
				ms.EXPECT().PhaseId().Return(phaseId).Times(1)
				ms.EXPECT().Turn().Return(turn).Times(1)
				mm.EXPECT().RegisterActionExecution(gomock.Any())
			},
		},
	}

	for _, test := range tests {
		rs.Run(test.name, func() {
			ctrl := gomock.NewController(rs.T())
			defer ctrl.Finish()
			scheduler := mock_game_logic.NewMockScheduler(ctrl)
			moderator := mock_game_logic.NewMockModerator(ctrl)
			action := mock_game_logic.NewMockAction(ctrl)

			action.EXPECT().Id().Return(rs.actionId).AnyTimes()
			moderator.EXPECT().Scheduler().Return(scheduler).AnyTimes()

			role := role{
				id:        id,
				moderator: moderator,
				playerId:  rs.playerId,
				abilities: []*ability{
					{
						action: action,
					},
				},
			}
			test.setup(role, action, scheduler, moderator)
			res := role.Use(test.req)

			rs.Equal(test.expectedRes, res)
			rs.Equal(test.expectedLimit, role.ActiveTimes(0))
		})
	}
}
