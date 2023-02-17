package role

import (
	"testing"
	"uwwolf/game/enum"
	"uwwolf/game/types"
	gamemock "uwwolf/mock/game"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RoleSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	game      *gamemock.MockGame
	player    *gamemock.MockPlayer
	action1   *gamemock.MockAction
	action2   *gamemock.MockAction
	playerID  enum.PlayerID
	actionID1 enum.ActionID
	actionID2 enum.ActionID
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

func (rs *RoleSuite) SetupSuite() {
	rs.playerID = "1"
	rs.actionID1 = 1
	rs.actionID2 = 2
}

func (rs *RoleSuite) BeforeTest(_, test string) {
	if test == "TestRoleUseAbility" {
		rs.ctrl = gomock.NewController(rs.T())
		rs.player = gamemock.NewMockPlayer(rs.ctrl)
		rs.action1 = gamemock.NewMockAction(rs.ctrl)
		rs.action2 = gamemock.NewMockAction(rs.ctrl)
		rs.action1.EXPECT().ID().Return(rs.actionID1).AnyTimes()
		rs.action2.EXPECT().ID().Return(rs.actionID2).AnyTimes()
	}

}

func (rs *RoleSuite) AfterTest(_, test string) {
	if test == "TestRoleUseAbility" {
		rs.ctrl.Finish()
	}
}

func (rs *RoleSuite) TestRoleID() {
	id := enum.HunterRoleID
	role := role{
		id: id,
	}

	rs.Equal(id, role.ID())
}

func (rs *RoleSuite) TestRolePhaseID() {
	phaseID := enum.DuskPhaseID
	role := role{
		phaseID: phaseID,
	}

	rs.Equal(phaseID, role.PhaseID())
}

func (rs *RoleSuite) TestRoleFactionID() {
	factionID := enum.WerewolfFactionID
	role := role{
		factionID: factionID,
	}

	rs.Equal(factionID, role.FactionID())
}

func (rs *RoleSuite) TestRolePriority() {
	priority := enum.VillagerTurnPriority
	role := role{
		priority: priority,
	}

	rs.Equal(priority, role.Priority())
}

func (rs *RoleSuite) TestRoleBeginRound() {
	round := enum.Round(9)
	role := role{
		beginRound: round,
	}

	rs.Equal(round, role.BeginRound())
}

func (rs *RoleSuite) TestRoleActiveLimit() {
	killLimit := enum.Limit(1)
	predictLimit := enum.Limit(1)
	role := role{
		abilities: map[enum.ActionID]*ability{
			enum.KillActionID: {
				action:      nil,
				activeLimit: killLimit,
			},
			enum.PredictActionID: {
				action:      nil,
				activeLimit: predictLimit,
			},
		},
	}

	rs.Equal(killLimit, role.ActiveLimit(enum.KillActionID))
	rs.Equal(killLimit+predictLimit, role.ActiveLimit(0))
	rs.Equal(enum.ReachedLimit, role.ActiveLimit(99))
}

func (rs *RoleSuite) TestRoleBeforeDeath() {
	role := new(role)

	rs.True(role.BeforeDeath())
}

func (rs *RoleSuite) TestRoleAfterDeath() {
	//
}

func (rs *RoleSuite) TestRoleUseAbility() {
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
				ActionID: rs.actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "Unable to use this ability anymore ¯\\(º_o)/¯",
			},
			newActiveLimit: 0,
			setup: func(role role) {
				role.abilities[rs.actionID1].activeLimit = 0
			},
		},
		{
			name: "Ok (Skip action -> Doesn't change active limit)",
			req: &types.UseRoleRequest{
				ActionID:  rs.actionID1,
				IsSkipped: true,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: true,
			},
			newActiveLimit: 5,
			setup: func(role role) {
				role.abilities[rs.actionID1].activeLimit = 5
				rs.action1.EXPECT().Execute(&types.ActionRequest{
					ActorID:   rs.playerID,
					IsSkipped: true,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: true,
					}).
					Times(1)
				rs.player.EXPECT().ID().Return(rs.playerID).Times(1)
			},
		},
		{
			name: "Ok (Action execution is failed -> Doesn't change active limit)",
			req: &types.UseRoleRequest{
				ActionID: rs.actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok: false,
			},
			newActiveLimit: 5,
			setup: func(role role) {
				role.abilities[rs.actionID1].activeLimit = 5
				rs.action1.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok: false,
					}).
					Times(1)
				rs.player.EXPECT().ID().Return(rs.playerID).Times(1)
			},
		},
		{
			name: "Ok (Use unlimited ability -> Doesn't change active limit)",
			req: &types.UseRoleRequest{
				ActionID: rs.actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			newActiveLimit: enum.Unlimited,
			setup: func(role role) {
				role.abilities[rs.actionID1].activeLimit = enum.Unlimited
				rs.action1.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)
				rs.player.EXPECT().ID().Return(rs.playerID).Times(1)
			},
		},
		{
			name: "Ok (Action execution is successful -> reduce active limit by 1)",
			req: &types.UseRoleRequest{
				ActionID: rs.actionID1,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: false,
			},
			newActiveLimit: 4,
			setup: func(role role) {
				role.abilities[rs.actionID1].activeLimit = 5
				rs.action1.EXPECT().Execute(&types.ActionRequest{
					ActorID: rs.playerID,
				}).
					Return(&types.ActionResponse{
						Ok:        true,
						IsSkipped: false,
					}).
					Times(1)
				rs.player.EXPECT().ID().Return(rs.playerID).Times(1)
			},
		},
	}

	for _, test := range tests {
		rs.Run(test.name, func() {
			role := role{
				player: rs.player,
				abilities: map[enum.ActionID]*ability{
					rs.actionID1: {
						action: rs.action1,
					},
					rs.actionID2: {
						action: rs.action2,
					},
				},
			}
			test.setup(role)
			res := role.UseAbility(test.req)

			rs.Equal(test.expectedRes, res)

			// Invalid ability
			if test.isInvalidRole {
				rs.Nil(role.abilities[test.req.ActionID])
			} else {
				rs.Equal(test.newActiveLimit, role.abilities[test.req.ActionID].activeLimit)
			}
		})
	}
}
