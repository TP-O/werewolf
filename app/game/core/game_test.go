package core_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"uwwolf/app/game/core"
// 	"uwwolf/types"
// )

// var playerIds = []types.PlayerId{
// 	"1111111111111111111111111111",
// 	"1111111111111111111111111112",
// 	"1111111111111111111111111113",
// 	"1111111111111111111111111114",
// 	"1111111111111111111111111115",
// }

// func TestGameRound(t *testing.T) {
// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds: playerIds,
// 	})

// 	assert.NotNil(t, g.Round())

// 	//=============================================================
// 	// Before role assignment
// 	assert.True(t, g.Round().IsEmpty())

// 	//=============================================================
// 	// After role assignment
// 	g.Start()

// 	assert.False(t, g.Round().IsEmpty())

// }

// func TestGamePoll(t *testing.T) {
// 	g := core.NewGame(&types.GameSetting{})

// 	assert.NotNil(t, g.Poll(types.VillagerFaction))
// 	assert.NotNil(t, g.Poll(types.WerewolfFaction))
// }

// func TestGameIsStarted(t *testing.T) {
// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds: playerIds,
// 		RolePool: []types.RoleId{
// 			types.HunterRole,
// 			types.SeerRole,
// 		},
// 		NumberOfWerewolves: 2,
// 	})

// 	//=============================================================
// 	// Not start yet
// 	assert.False(t, g.IsStarted())

// 	//=============================================================
// 	// Started
// 	g.Start()

// 	assert.True(t, g.IsStarted())
// }

// func TestGameStart(t *testing.T) {
// 	numberOfWerewolves := 2
// 	rolePool := []types.RoleId{
// 		types.HunterRole,
// 		types.SeerRole,
// 	}

// 	g1 := core.NewGame(&types.GameSetting{
// 		PlayerIds:          playerIds,
// 		RolePool:           rolePool,
// 		NumberOfWerewolves: numberOfWerewolves,
// 	})
// 	g2 := core.NewGame(&types.GameSetting{
// 		PlayerIds:          playerIds,
// 		RolePool:           rolePool,
// 		NumberOfWerewolves: numberOfWerewolves,
// 	})

// 	selectedRoleIds := append(rolePool, types.VillagerRole, types.WerewolfRole)

// 	//=============================================================
// 	// Not start yet
// 	assert.False(t, g1.IsStarted())
// 	assert.Empty(t, g1.PlayerIdsWithFaction(types.VillagerFaction))
// 	assert.Empty(t, g1.PlayerIdsWithFaction(types.WerewolfFaction))
// 	assert.Empty(t, g1.PlayerIdsWithFaction(types.VillagerFaction))
// 	assert.Empty(t, g1.PlayerIdsWithFaction(types.WerewolfFaction))

// 	//=============================================================
// 	// Start successfully
// 	g1.Start()
// 	assert.True(t, g1.IsStarted())
// 	assert.Len(t, g1.PlayerIdsWithFaction(types.VillagerFaction), len(playerIds)-numberOfWerewolves)
// 	assert.Len(t, g1.PlayerIdsWithFaction(types.WerewolfFaction), numberOfWerewolves)

// 	for _, pId := range playerIds {
// 		assert.NotNil(t, g1.Player(pId))
// 		assert.Equal(t, pId, g1.Player(pId).Id())
// 		assert.NotEmpty(t, g1.Player(pId).RoleIds())
// 		assert.Contains(t, g1.PlayerIdsWithRole(types.VillagerRole), pId)

// 		for _, rId := range g1.Player(pId).RoleIds() {
// 			assert.Contains(t, selectedRoleIds, rId)
// 		}
// 	}

// 	assert.Len(t, g1.Poll(types.VillagerFaction).ElectorIds, len(playerIds))
// 	assert.Len(t, g1.Poll(types.WerewolfFaction).ElectorIds, numberOfWerewolves)

// 	// These assertions may fail because role assignment is random.
// 	// Run test again if it happens
// 	assert.NotEqual(t, g1.PlayerIdsWithRole(types.VillagerRole), g2.PlayerIdsWithRole(types.VillagerRole))
// 	assert.NotEqual(t, g1.PlayerIdsWithRole(types.WerewolfRole), g2.PlayerIdsWithRole(types.WerewolfRole))
// 	assert.NotEqual(t, g1.PlayerIdsWithRole(types.HunterRole), g2.PlayerIdsWithRole(types.HunterRole))
// 	assert.NotEqual(t, g1.PlayerIdsWithRole(types.SeerRole), g2.PlayerIdsWithRole(types.SeerRole))
// }

// func TestGamePlayer(t *testing.T) {
// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds: playerIds,
// 	})

// 	for _, pId := range playerIds {
// 		assert.NotNil(t, g.Player(pId))
// 		assert.Equal(t, pId, g.Player(pId).Id())
// 	}
// }

// func TestGamePlayerIdsWithRole(t *testing.T) {
// 	numberOfWerewolves := 1
// 	rolePool := []types.RoleId{
// 		types.HunterRole,
// 		types.SeerRole,
// 	}
// 	selectedRoleIds := append(rolePool, types.VillagerRole, types.WerewolfRole)

// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds:          playerIds,
// 		RolePool:           rolePool,
// 		NumberOfWerewolves: numberOfWerewolves,
// 	})
// 	g.Start()

// 	for _, rId := range selectedRoleIds {
// 		assert.NotEmpty(t, g.PlayerIdsWithRole(rId))
// 	}
// }

// func TestGamePlayerIdsWithFaction(t *testing.T) {
// 	numberOfWerewolves := 1

// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds:          playerIds,
// 		NumberOfWerewolves: numberOfWerewolves,
// 	})
// 	g.Start()

// 	assert.NotEmpty(t, g.PlayerIdsWithFaction(types.VillagerFaction))
// 	assert.NotEmpty(t, g.PlayerIdsWithFaction(types.WerewolfFaction))
// }

// func TestGameKillPlayer(t *testing.T) {
// 	numberOfWerewolves := 2
// 	rolePool := []types.RoleId{
// 		types.HunterRole,
// 		types.SeerRole,
// 	}

// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds:          playerIds,
// 		RolePool:           rolePool,
// 		NumberOfWerewolves: numberOfWerewolves,
// 	})

// 	g.Start()

// 	//=============================================================
// 	// Kill non-exist player
// 	assert.Nil(t, g.KillPlayer(types.PlayerId("99")))

// 	//=============================================================
// 	// Kill player successfully
// 	killedPlayer := g.KillPlayer(playerIds[0])

// 	assert.NotNil(t, killedPlayer)
// 	assert.Equal(t, playerIds[0], killedPlayer.Id())
// 	assert.NotContains(t, g.Poll(types.VillagerFaction).ElectorIds, playerIds[0])
// 	assert.NotContains(t, g.Poll(types.WerewolfFaction).ElectorIds, playerIds[0])

// 	// Loop through all turns
// 	for i := 0; i < 4; i++ {
// 		assert.NotContains(t, g.Round().CurrentTurn().PlayerIds(), playerIds[0])

// 		g.Round().NextTurn()
// 	}
// }

// func TestGameRequestAction(t *testing.T) {
// 	numberOfWerewolves := 2
// 	rolePool := []types.RoleId{
// 		types.HunterRole,
// 		types.SeerRole,
// 	}

// 	g := core.NewGame(&types.GameSetting{
// 		PlayerIds:          playerIds,
// 		RolePool:           rolePool,
// 		NumberOfWerewolves: numberOfWerewolves,
// 	})

// 	g.Start()
// 	g.KillPlayer(playerIds[4])
// 	g.Poll(types.VillagerFaction).Open()

// 	// Move to villager turn
// 	g.Round().NextTurn()
// 	g.Round().NextTurn()

// 	//=============================================================
// 	// Invalid input
// 	res := g.RequestAction(&types.ActionRequest{
// 		ActorId:   playerIds[3],
// 		TargetIds: []types.PlayerId{playerIds[0]},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.InvalidInputErrorTag, res.Error.Tag)

// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   playerIds[3],
// 		TargetIds: []types.PlayerId{},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.InvalidInputErrorTag, res.Error.Tag)

// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   playerIds[3],
// 		TargetIds: []types.PlayerId{playerIds[0], playerIds[0]},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.InvalidInputErrorTag, res.Error.Tag)

// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   playerIds[3],
// 		TargetIds: []types.PlayerId{types.UnknownPlayer},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.InvalidInputErrorTag, res.Error.Tag)

// 	//=============================================================
// 	// Player is dead
// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   playerIds[4],
// 		TargetIds: []types.PlayerId{playerIds[0]},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.UnauthorizedErrorTag, res.Error.Tag)

// 	//=============================================================
// 	// Non-exist player
// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   types.PlayerId("0000000000000000000000000000"),
// 		TargetIds: []types.PlayerId{playerIds[0]},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.UnauthorizedErrorTag, res.Error.Tag)

// 	//=============================================================
// 	// Wrong turn
// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   types.PlayerId("1111111111111111111111111116"),
// 		TargetIds: []types.PlayerId{playerIds[0]},
// 	})

// 	assert.False(t, res.Ok)
// 	assert.Equal(t, types.UnauthorizedErrorTag, res.Error.Tag)

// 	//=============================================================
// 	// Request successfully
// 	res = g.RequestAction(&types.ActionRequest{
// 		GameId:    1,
// 		ActorId:   playerIds[0],
// 		TargetIds: []types.PlayerId{playerIds[0]},
// 	})

// 	assert.True(t, res.Ok)
// }
