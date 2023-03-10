package game

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"
	gamemock "uwwolf/mock/game"
	"uwwolf/util"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type GameSuite struct {
	suite.Suite
	player1ID types.PlayerID
	player2ID types.PlayerID
	player3ID types.PlayerID
	role1ID   types.RoleID
	role2ID   types.RoleID
	role3ID   types.RoleID
}

func (gs *GameSuite) SetupSuite() {
	gs.player1ID = types.PlayerID("1")
	gs.player2ID = types.PlayerID("2")
	gs.player3ID = types.PlayerID("3")
	gs.role1ID = types.RoleID(1)
	gs.role2ID = types.RoleID(2)
	gs.role3ID = types.RoleID(3)
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}

func (gs GameSuite) TestNewGame() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)

	setting := &types.GameSetting{
		NumberWerewolves: 10,
		RoleIDs:          []types.RoleID{gs.role1ID, gs.role2ID, gs.role3ID},
		RequiredRoleIDs:  []types.RoleID{gs.role1ID, gs.role2ID},
		PlayerIDs:        []types.PlayerID{gs.player1ID, gs.player2ID, gs.player3ID},
	}

	g := NewGame(scheduler, setting).(*game)

	gs.Equal(setting.NumberWerewolves, g.numberWerewolves)
	gs.Equal(setting.RoleIDs, g.roleIDs)
	gs.Equal(setting.RequiredRoleIDs, g.requiredRoleIDs)
	gs.Equal(vars.Idle, g.statusID)
	gs.Equal(scheduler, g.scheduler)
	gs.NotNil(g.players)
	gs.ElementsMatch(setting.PlayerIDs, maps.Keys(g.players))
	for _, p := range g.players {
		gs.NotNil(p)
	}
	gs.NotNil(g.polls)
	gs.NotNil(g.polls[vars.VillagerFactionID])
	gs.NotNil(g.polls[vars.WerewolfFactionID])
}

func (gs GameSuite) TestStatusID() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).statusID = vars.Starting

	gs.Equal(vars.Starting, g.StatusID())
}

func (gs GameSuite) TestScheduler() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)

	g := NewGame(scheduler, &types.GameSetting{})

	gs.Same(scheduler, g.Scheduler())
}

func (gs GameSuite) TestPoll() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	poll := gamemock.NewMockPoll(ctrl)

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).polls[vars.VillagerFactionID] = poll

	gs.Same(poll, g.Poll(vars.VillagerFactionID))
}

func (gs GameSuite) TestPlayer() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	player := gamemock.NewMockPlayer(ctrl)

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).players[gs.player1ID] = player

	gs.Same(player, g.Player(gs.player1ID))
}

func (gs GameSuite) TestPlayers() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	player1 := gamemock.NewMockPlayer(ctrl)
	player2 := gamemock.NewMockPlayer(ctrl)
	player3 := gamemock.NewMockPlayer(ctrl)

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).players[gs.player1ID] = player1
	g.(*game).players[gs.player2ID] = player2
	g.(*game).players[gs.player3ID] = player3

	gs.Len(g.Players(), 3)
	gs.Same(player1, g.Players()[gs.player1ID])
	gs.Same(player2, g.Players()[gs.player2ID])
	gs.Same(player3, g.Players()[gs.player3ID])
}

func (gs GameSuite) TestAlivePlayerIDsWithRoleID() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	player1 := gamemock.NewMockPlayer(ctrl)
	player2 := gamemock.NewMockPlayer(ctrl)
	player3 := gamemock.NewMockPlayer(ctrl)

	player1.EXPECT().IsDead().Return(true)
	player2.EXPECT().IsDead().Return(false)
	player2.EXPECT().RoleIDs().Return([]types.RoleID{gs.role1ID, gs.role2ID})
	player3.EXPECT().IsDead().Return(false)
	player3.EXPECT().RoleIDs().Return([]types.RoleID{gs.role1ID, gs.role3ID})

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).players[gs.player1ID] = player1
	g.(*game).players[gs.player2ID] = player2
	g.(*game).players[gs.player3ID] = player3

	pIDs := g.AlivePlayerIDsWithRoleID(gs.role2ID)

	gs.ElementsMatch([]types.PlayerID{gs.player2ID}, pIDs)
}

func (gs GameSuite) TestAlivePlayerIDsWithFactionID() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	player1 := gamemock.NewMockPlayer(ctrl)
	player2 := gamemock.NewMockPlayer(ctrl)
	player3 := gamemock.NewMockPlayer(ctrl)

	player1.EXPECT().IsDead().Return(true)
	player2.EXPECT().IsDead().Return(false)
	player2.EXPECT().FactionID().Return(vars.WerewolfFactionID)
	player3.EXPECT().IsDead().Return(false)
	player3.EXPECT().FactionID().Return(vars.VillagerFactionID)

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).players[gs.player1ID] = player1
	g.(*game).players[gs.player2ID] = player2
	g.(*game).players[gs.player3ID] = player3

	pIDs := g.AlivePlayerIDsWithFactionID(vars.WerewolfFactionID)

	gs.ElementsMatch([]types.PlayerID{gs.player2ID}, pIDs)
}

func (gs GameSuite) TestAlivePlayerIDsWithoutFactionID() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	player1 := gamemock.NewMockPlayer(ctrl)
	player2 := gamemock.NewMockPlayer(ctrl)
	player3 := gamemock.NewMockPlayer(ctrl)

	player1.EXPECT().IsDead().Return(true)
	player2.EXPECT().IsDead().Return(false)
	player2.EXPECT().FactionID().Return(vars.WerewolfFactionID)
	player3.EXPECT().IsDead().Return(false)
	player3.EXPECT().FactionID().Return(vars.VillagerFactionID)

	g := NewGame(scheduler, &types.GameSetting{})
	g.(*game).players[gs.player1ID] = player1
	g.(*game).players[gs.player2ID] = player2
	g.(*game).players[gs.player3ID] = player3

	pIDs := g.AlivePlayerIDsWithoutFactionID(vars.WerewolfFactionID)

	gs.ElementsMatch([]types.PlayerID{gs.player3ID}, pIDs)
}

func (gs GameSuite) TestSelectRoleID() {
	tests := []struct {
		name                       string
		setting                    *types.GameSetting
		roleID                     types.RoleID
		werewolfCounter            int
		nonWerewolfCounter         int
		expectedStatus             bool
		expectedWerewolfCounter    int
		expectedNonWerewolfCounter int
		expectedSelectedRoleIDs    []types.RoleID
	}{
		{
			name: "False (Selected role id list is enough)",
			setting: &types.GameSetting{
				NumberWerewolves: 1,
				PlayerIDs: []types.PlayerID{
					gs.player1ID,
					gs.player2ID,
					gs.player3ID,
				},
			},
			roleID:                     vars.HunterRoleID,
			expectedStatus:             false,
			werewolfCounter:            1,
			nonWerewolfCounter:         2,
			expectedWerewolfCounter:    1,
			expectedNonWerewolfCounter: 2,
			expectedSelectedRoleIDs:    []types.RoleID{},
		},
		{
			name: "True (Non-Werewolf role)",
			setting: &types.GameSetting{
				NumberWerewolves: 1,
				PlayerIDs: []types.PlayerID{
					gs.player1ID,
					gs.player2ID,
					gs.player3ID,
				},
			},
			roleID:                     vars.HunterRoleID,
			expectedStatus:             true,
			werewolfCounter:            1,
			nonWerewolfCounter:         1,
			expectedWerewolfCounter:    1,
			expectedNonWerewolfCounter: 2,
			expectedSelectedRoleIDs:    []types.RoleID{vars.HunterRoleID},
		},
		// Update if new werewolf role added
		{
			name: "True (Werewolf role)",
			setting: &types.GameSetting{
				NumberWerewolves: 1,
				PlayerIDs: []types.PlayerID{
					gs.player1ID,
					gs.player2ID,
					gs.player3ID,
				},
			},
			roleID:                     vars.WerewolfRoleID,
			expectedStatus:             true,
			werewolfCounter:            0,
			nonWerewolfCounter:         2,
			expectedWerewolfCounter:    1,
			expectedNonWerewolfCounter: 2,
			expectedSelectedRoleIDs:    []types.RoleID{vars.WerewolfRoleID},
		},
	}

	for _, test := range tests {
		gs.Run(test.name, func() {
			ctrl := gomock.NewController(gs.T())
			defer ctrl.Finish()
			scheduler := gamemock.NewMockScheduler(ctrl)

			// Make werewolf role can be used in selectRoleID
			vars.RoleSets[vars.WerewolfRoleID] = 1
			defer func() {
				vars.RoleSets[vars.WerewolfRoleID] = vars.UnlimitedTimes
			}()

			g := NewGame(scheduler, test.setting)
			ok := g.(*game).selectRoleID(&test.werewolfCounter, &test.nonWerewolfCounter, test.roleID)

			gs.Equal(test.expectedStatus, ok)
			gs.Equal(test.expectedWerewolfCounter, test.werewolfCounter)
			gs.Equal(test.expectedNonWerewolfCounter, test.nonWerewolfCounter)
			gs.ElementsMatch(test.expectedSelectedRoleIDs, g.(*game).selectedRoleIDs)
		})
	}
}

// Update if new werewolf role added
func (gs GameSuite) TestSelectRoleIDs() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)

	roleIDs := []types.RoleID{
		vars.HunterRoleID,
		vars.SeerRoleID,
		vars.TwoSistersRoleID,
	}
	requiredRoleIDs := []types.RoleID{
		vars.SeerRoleID,
	}

	// Check random
	for i := 0; i < 10; i++ {
		g := NewGame(scheduler, &types.GameSetting{
			NumberWerewolves: 1,
			RoleIDs:          roleIDs,
			RequiredRoleIDs:  requiredRoleIDs,
			PlayerIDs: []types.PlayerID{
				gs.player1ID,
				gs.player2ID,
				gs.player3ID,
			},
		})
		g.(*game).selectRoleIDs()

		for _, roleID := range requiredRoleIDs {
			gs.Contains(g.(*game).selectedRoleIDs, roleID)
		}
		gs.Len(g.(*game).selectedRoleIDs, 2)
		gs.False(util.IsDuplicateSlice(g.(*game).selectedRoleIDs))
		gs.Condition(func() (success bool) {
			return slices.Contains(g.(*game).selectedRoleIDs, vars.HunterRoleID) ||
				slices.Contains(g.(*game).selectedRoleIDs, vars.TwoSistersRoleID)
		})
	}
}

func (gs GameSuite) TestAssignRoles() {
	tests := []struct {
		name  string
		setup func(*game, *gamemock.MockPlayer)
	}{
		{
			name: "Only assign default villager faction's role",
			setup: func(g *game, mp1 *gamemock.MockPlayer) {
				g.selectedRoleIDs = []types.RoleID{}

				mp1.EXPECT().AssignRole(vars.VillagerRoleID)
			},
		},
		{
			name: "Assign selected role and default villager faction's role",
			setup: func(g *game, mp1 *gamemock.MockPlayer) {
				g.selectedRoleIDs = []types.RoleID{vars.HunterRoleID}

				mp1.EXPECT().AssignRole(vars.VillagerRoleID)
				mp1.EXPECT().AssignRole(vars.HunterRoleID)
			},
		},
		// Update if new werewolf role added
		{
			name: "Assign selected role and default villager and werewolf faction's role",
			setup: func(g *game, mp1 *gamemock.MockPlayer) {
				g.selectedRoleIDs = []types.RoleID{vars.WerewolfRoleID}

				mp1.EXPECT().AssignRole(vars.VillagerRoleID)
				mp1.EXPECT().AssignRole(vars.WerewolfRoleID)
				mp1.EXPECT().AssignRole(vars.WerewolfRoleID)
			},
		},
	}

	for _, test := range tests {
		gs.Run(test.name, func() {
			ctrl := gomock.NewController(gs.T())
			defer ctrl.Finish()
			scheduler := gamemock.NewMockScheduler(ctrl)
			player1 := gamemock.NewMockPlayer(ctrl)

			player1.EXPECT().ID().Return(gs.player1ID)

			g := NewGame(scheduler, &types.GameSetting{})

			g.(*game).players[gs.player1ID] = player1
			test.setup(g.(*game), player1)

			g.(*game).assignRoles()
		})
	}
}

func (gs GameSuite) TestPrepare() {
	ctrl := gomock.NewController(gs.T())
	defer ctrl.Finish()
	scheduler := gamemock.NewMockScheduler(ctrl)
	player1 := gamemock.NewMockPlayer(ctrl)
	player2 := gamemock.NewMockPlayer(ctrl)
	player3 := gamemock.NewMockPlayer(ctrl)

	// assignRoles() is called
	player1.EXPECT().ID().Return(gs.player1ID)
	player2.EXPECT().ID().Return(gs.player2ID)
	player3.EXPECT().ID().Return(gs.player3ID)

	player1.EXPECT().AssignRole(vars.VillagerRoleID)
	player2.EXPECT().AssignRole(vars.VillagerRoleID)
	player3.EXPECT().AssignRole(vars.VillagerRoleID)

	player1.EXPECT().AssignRole(vars.WerewolfRoleID).MaxTimes(1)
	player2.EXPECT().AssignRole(vars.WerewolfRoleID).MaxTimes(1)
	player3.EXPECT().AssignRole(vars.WerewolfRoleID).MaxTimes(1)

	player1.EXPECT().AssignRole(gomock.Any()).MaxTimes(1)
	player2.EXPECT().AssignRole(gomock.Any()).MaxTimes(1)
	player3.EXPECT().AssignRole(gomock.Any()).MaxTimes(1)

	roleIDs := []types.RoleID{
		vars.HunterRoleID,
		vars.SeerRoleID,
		vars.TwoSistersRoleID,
	}
	requiredRoleIDs := []types.RoleID{
		vars.SeerRoleID,
	}

	g := NewGame(scheduler, &types.GameSetting{
		NumberWerewolves: 1,
		RoleIDs:          roleIDs,
		RequiredRoleIDs:  requiredRoleIDs,
		PlayerIDs: []types.PlayerID{
			gs.player1ID,
			gs.player2ID,
			gs.player3ID,
		},
	})
	g.(*game).players[gs.player1ID] = player1
	g.(*game).players[gs.player2ID] = player2
	g.(*game).players[gs.player3ID] = player3

	//============================ Current status isn't Idle
	g.(*game).statusID = vars.Starting
	t := g.Prepare()

	gs.Equal(vars.Starting, g.StatusID())
	gs.Negative(t)

	//============================ Current status isn't Idle
	g.(*game).statusID = vars.Idle
	t = g.Prepare()

	gs.Equal(vars.Waiting, g.StatusID())
	gs.Positive(t)

	// selectRoleIDs() is called
	for _, roleID := range requiredRoleIDs {
		gs.Contains(g.(*game).selectedRoleIDs, roleID)
	}
	gs.Len(g.(*game).selectedRoleIDs, 2)
	gs.False(util.IsDuplicateSlice(g.(*game).selectedRoleIDs))
	gs.Condition(func() (success bool) {
		return slices.Contains(g.(*game).selectedRoleIDs, vars.HunterRoleID) ||
			slices.Contains(g.(*game).selectedRoleIDs, vars.TwoSistersRoleID)
	})
}

func (gs GameSuite) TestStart() {
	tests := []struct {
		name           string
		statusID       types.GameStatusID
		expectedStatus bool
	}{
		{
			name:           "False (Current status isn't Waiting)",
			statusID:       vars.Starting,
			expectedStatus: false,
		},
		{
			name:           "Ok",
			statusID:       vars.Waiting,
			expectedStatus: true,
		},
	}

	for _, test := range tests {
		gs.Run(test.name, func() {
			ctrl := gomock.NewController(gs.T())
			defer ctrl.Finish()
			scheduler := gamemock.NewMockScheduler(ctrl)

			g := NewGame(scheduler, &types.GameSetting{})
			g.(*game).statusID = test.statusID

			ok := g.Start()

			gs.Equal(test.expectedStatus, ok)
		})
	}
}

func (gs GameSuite) TestFinish() {
	tests := []struct {
		name           string
		statusID       types.GameStatusID
		expectedStatus bool
	}{
		{
			name:           "False (Current status is already Finished)",
			statusID:       vars.Finished,
			expectedStatus: false,
		},
		{
			name:           "Ok",
			statusID:       vars.Waiting,
			expectedStatus: true,
		},
	}

	for _, test := range tests {
		gs.Run(test.name, func() {
			ctrl := gomock.NewController(gs.T())
			defer ctrl.Finish()
			scheduler := gamemock.NewMockScheduler(ctrl)

			g := NewGame(scheduler, &types.GameSetting{})
			g.(*game).statusID = test.statusID

			ok := g.Start()

			gs.Equal(test.expectedStatus, ok)
		})
	}
}

func (gs GameSuite) TestPlay() {
	tests := []struct {
		name        string
		req         *types.ActivateAbilityRequest
		expectedRes *types.ActionResponse
		setup       func(*game, *gamemock.MockPlayer)
	}{
		{
			name: "Failure (Game isn't starting)",
			req:  &types.ActivateAbilityRequest{},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "The game is about to start ノ(ジ)ー'",
			},
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.statusID = vars.Waiting
			},
		},
		{
			name: "Failure (Game isn't starting)",
			req:  &types.ActivateAbilityRequest{},
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Unable to play this game (╥﹏╥)",
			},
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.statusID = vars.Starting
				g.players[gs.player1ID] = nil
			},
		},
		{
			name: "Ok",
			req:  &types.ActivateAbilityRequest{},
			expectedRes: &types.ActionResponse{
				Ok:      true,
				Message: "Ok msg",
			},
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.statusID = vars.Starting
				g.players[gs.player1ID] = mp

				mp.EXPECT().ActivateAbility(&types.ActivateAbilityRequest{}).
					Return(&types.ActionResponse{
						Ok:      true,
						Message: "Ok msg",
					})
			},
		},
	}

	for _, test := range tests {
		gs.Run(test.name, func() {
			ctrl := gomock.NewController(gs.T())
			defer ctrl.Finish()
			scheduler := gamemock.NewMockScheduler(ctrl)
			player1 := gamemock.NewMockPlayer(ctrl)

			g := NewGame(scheduler, &types.GameSetting{})
			test.setup(g.(*game), player1)

			res := g.Play(gs.player1ID, test.req)

			gs.Equal(test.expectedRes, res)
		})
	}
}

func (gs GameSuite) TestKillPlayer() {
	tests := []struct {
		name                string
		playerID            types.PlayerID
		isExited            bool
		expectedIsNilPlayer bool
		setup               func(*game, *gamemock.MockPlayer)
	}{
		{
			name:                "Failure (Non-existent player)",
			playerID:            gs.player1ID,
			isExited:            false,
			expectedIsNilPlayer: true,
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.players[gs.player1ID] = nil
			},
		},
		{
			name:                "Failure (Play is already dead)",
			playerID:            gs.player1ID,
			isExited:            false,
			expectedIsNilPlayer: true,
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.players[gs.player1ID] = mp

				mp.EXPECT().IsDead().Return(true)
			},
		},
		{
			name:                "Failure (Play is protected by role)",
			playerID:            gs.player1ID,
			isExited:            true,
			expectedIsNilPlayer: true,
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.players[gs.player1ID] = mp

				mp.EXPECT().IsDead().Return(false)
				mp.EXPECT().Die(true).Return(false)
			},
		},
		{
			name:                "Ok",
			playerID:            gs.player1ID,
			isExited:            false,
			expectedIsNilPlayer: false,
			setup: func(g *game, mp *gamemock.MockPlayer) {
				g.players[gs.player1ID] = mp

				mp.EXPECT().IsDead().Return(false)
				mp.EXPECT().Die(false).Return(true)
			},
		},
	}

	for _, test := range tests {
		gs.Run(test.name, func() {
			ctrl := gomock.NewController(gs.T())
			defer ctrl.Finish()
			scheduler := gamemock.NewMockScheduler(ctrl)
			player1 := gamemock.NewMockPlayer(ctrl)

			g := NewGame(scheduler, &types.GameSetting{})
			test.setup(g.(*game), player1)

			p := g.KillPlayer(test.playerID, test.isExited)

			if test.expectedIsNilPlayer {
				gs.Nil(p)
			} else {
				gs.NotNil(p)
			}
		})
	}
}
