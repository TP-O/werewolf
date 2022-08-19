package core

import (
	"time"

	"uwwolf/module/game/contract"
	"uwwolf/module/game/state"
	"uwwolf/types"

	"golang.org/x/exp/slices"
)

// import (
// 	"errors"
// 	"fmt"
// 	"sort"
// 	"time"

// 	"golang.org/x/exp/slices"

// 	"uwwolf/app/model"
// 	"uwwolf/contract/typ"
// 	"uwwolf/database"
// 	"uwwolf/enum"
// 	"uwwolf/game/factory"
// 	"uwwolf/game/stuff"
// 	"uwwolf/util"
// 	"uwwolf/validator"
// )

// type game struct {
// 	gameId              string
// 	capacity            int
// 	numberOfWerewolves  int
// 	remaining           int
// 	remainingWerewolves int
// 	socketId2playerId   map[string]int
// 	playerId2SocketId   map[int]string
// 	playerId2RoleId     map[int]int
// 	roleId2SocketIds    map[int][]string
// 	roleId2PlayerIds    map[int][]int
// 	isStarted           bool
// 	rolePool            []int
// 	data                chan string
// 	phase               *stuff.Phase
// }

// func NewGameInstance(input *typ.GameInstanceInit) (*instance, error) {
// 	if !validator.ValidateStruct(input) {
// 		return nil, errors.New("Invalid game instance!")
// 	}

// 	gameInstance := instance{
// 		gameId:              input.GameId,
// 		capacity:            input.Capacity,
// 		numberOfWerewolves:  input.NumberOfWerewolves,
// 		remaining:           input.Capacity,
// 		remainingWerewolves: input.NumberOfWerewolves,
// 		isStarted:           false,
// 		rolePool:            input.RolePool,
// 		phase:               &stuff.Phase{},
// 		data:                make(chan string),
// 	}

// 	gameInstance.phase.Init()

// 	return &gameInstance, nil
// }

// func (g *game) IsStarted() bool {
// 	return g.isStarted
// }

// func (g *game) NumberOfVillagers() int {
// 	return g.remaining
// }

// func (g *game) NumberOfWerewolves() int {
// 	return g.remainingWerewolves
// }

// // Start game instance
// func (g *game) Start() bool {
// 	if g.isStarted || len(g.socketId2playerId) != int(g.capacity) {
// 		return false
// 	}

// 	g.isStarted = true

// 	g.assignRoles()
// 	go g.listen()
// 	go g.phase.Start()

// 	return g.isStarted
// }

// // Replace players
// func (g *game) AddPlayers(socketIds []string, playerIds []int) bool {
// 	if g.isStarted ||
// 		len(socketIds) != int(g.capacity) ||
// 		len(playerIds) != int(g.capacity) ||
// 		len(socketIds) != len(playerIds) ||
// 		!validator.ValidateVar(socketIds, "unique,dive,len=20,alphanum") ||
// 		!validator.ValidateVar(playerIds, "unique,dive,min=1") {

// 		return false
// 	}

// 	g.resetPlayers()

// 	for index, socketId := range socketIds {
// 		g.socketId2playerId[socketId] = playerIds[index]
// 		g.playerId2SocketId[playerIds[index]] = socketId
// 	}

// 	return true
// }

// // Remove a player
// func (g *game) RemovePlayer(socketId string) bool {
// 	if g.socketId2playerId[socketId] == 0 {
// 		return false
// 	}

// 	delete(g.playerId2SocketId, g.socketId2playerId[socketId])
// 	delete(g.socketId2playerId, socketId)

// 	return true
// }

// func (g *game) NextTurn() {
// 	go g.phase.NextTurn()
// }

// func (g *game) Pipe(pub *chan string) {
// 	*pub = g.data
// }

// func (g *game) Do(instruction *typ.ActionInstruction) bool {
// 	if !g.phase.IsValidPlayer(instruction.Actor) {
// 		return false
// 	}

// 	return g.phase.UseSkill(instruction)
// }

// func (g *game) listen() {
// 	for d := range g.data {
// 		fmt.Println(d)
// 	}
// }

// // Assign roles to players randomly
// func (g *game) assignRoles() {
// 	var specialRoles []model.Role
// 	database.DB().Find(&specialRoles, g.rolePool)

// 	randomRoles := g.pickUpRoles(specialRoles)
// 	cloneRandomRoles := slices.Clone(randomRoles)

// 	for socketId, playerId := range g.socketId2playerId {
// 		randIndex := util.RandomIndex(randomRoles)
// 		randomRole := randomRoles[randIndex]

// 		g.playerId2RoleId[playerId] = randomRole.ID
// 		g.roleId2SocketIds[randomRole.ID] = append(g.roleId2SocketIds[randomRole.ID], socketId)
// 		g.roleId2PlayerIds[randomRole.ID] = append(g.roleId2PlayerIds[randomRole.ID], playerId)

// 		if randomRole.ID != enum.VillagerRole {
// 			g.roleId2SocketIds[enum.VillagerRole] = append(g.roleId2SocketIds[enum.VillagerRole], socketId)
// 			g.roleId2PlayerIds[enum.VillagerRole] = append(g.roleId2PlayerIds[enum.VillagerRole], playerId)
// 		}

// 		if randomRole.ID != enum.WerewolfRole && randomRole.FactionID == enum.WerewolfFaction {
// 			g.roleId2SocketIds[enum.WerewolfRole] = append(g.roleId2SocketIds[enum.WerewolfRole], socketId)
// 			g.roleId2PlayerIds[enum.WerewolfRole] = append(g.roleId2PlayerIds[enum.WerewolfRole], playerId)
// 		}

// 		randomRoles = append(randomRoles[:randIndex], randomRoles[randIndex+1:]...)
// 	}

// 	g.setUpTurns(cloneRandomRoles)
// }

// // Pick up number of roles corresponding to game capacity randomly
// func (g *game) pickUpRoles(roles []model.Role) []model.Role {
// 	var defaultRoles []model.Role
// 	database.DB().Find(&defaultRoles, []int{enum.VillagerRole, enum.WerewolfRole})

// 	randomRoles := make([]model.Role, g.capacity)

// 	werewolfRoles, remainingRoles := g.classifyRoles(roles)

// 	for j := int(0); j < g.capacity; j++ {
// 		var currentRoles *[]model.Role

// 		if j < g.numberOfWerewolves {
// 			currentRoles = &werewolfRoles
// 		} else {
// 			currentRoles = &remainingRoles
// 		}

// 		// Fill in with default role if all special roles are taken
// 		if len(*currentRoles) == 0 {
// 			if currentRoles == &werewolfRoles {
// 				randomRoles[j] = defaultRoles[1]
// 			} else {
// 				randomRoles[j] = defaultRoles[0]
// 			}

// 			continue
// 		}

// 		randIndex := util.RandomIndex(*currentRoles)
// 		randomRoles[j] = (*currentRoles)[randIndex]

// 		(*currentRoles)[randIndex].Quantity--

// 		// Delete role if its quantity runs out
// 		if (*currentRoles)[randIndex].Quantity == 0 {
// 			*currentRoles = append((*currentRoles)[:randIndex], (*currentRoles)[randIndex+1:]...)
// 		}
// 	}

// 	return randomRoles
// }

// // Classify roles into 2 factions
// func (g *game) classifyRoles(roles []model.Role) ([]model.Role, []model.Role) {
// 	var werewolfRoles []model.Role
// 	var remainingRoles []model.Role

// 	for _, role := range roles {
// 		if role.FactionID == enum.WerewolfFaction {
// 			werewolfRoles = append(werewolfRoles, role)
// 		} else {
// 			remainingRoles = append(remainingRoles, role)
// 		}
// 	}

// 	return werewolfRoles, remainingRoles
// }

// // Reset player-related data
// func (g *game) resetPlayers() {
// 	g.socketId2playerId = make(map[string]int)
// 	g.playerId2SocketId = make(map[int]string)
// 	g.playerId2RoleId = make(map[int]int)
// 	g.roleId2SocketIds = make(map[int][]string)
// 	g.roleId2PlayerIds = make(map[int][]int)
// }

// // Prepare turn for special roles
// func (g *game) setUpTurns(roles []model.Role) {
// 	roleFactory := factory.GetRoleFactory()
// 	roles = util.RemoveDuplicate(roles)

// 	// Order by priority
// 	sort.Slice(roles, func(i, j int) bool {
// 		return roles[i].Priority < roles[j].Priority
// 	})

// 	for _, role := range roles {
// 		g.phase.AddTurn(
// 			role.PhaseID,
// 			2*time.Second,
// 			roleFactory.Create(role.ID, i),
// 			g.roleId2PlayerIds[role.ID],
// 		)
// 	}

// 	// Test
// 	g.phase.AddPlayer(enum.NightPhase, 0, 1, 3)

// 	for i, phases := range g.phase.Data() {
// 		fmt.Println("Phase: ", i)

// 		for _, turn := range phases {
// 			fmt.Println(turn)
// 		}
// 	}
// }

type game struct {
	id                 types.GameId
	capacity           uint
	numberOfWerewolves uint
	timeForTurn        time.Duration
	timeForDiscussion  time.Duration
	rolePool           []types.RoleId
	factions           map[types.FactionId][]types.PlayerId
	players            map[types.PlayerId]contract.Player
	deaths             []types.PlayerId
	polls              map[string]*state.Poll
	round              *state.Round
}

func NewGame(setting *types.GameSetting) contract.Game {
	game := game{
		id:                 setting.Id,
		capacity:           uint(len(setting.PlayerIds)),
		numberOfWerewolves: setting.NumberOfWerewolves,
		timeForTurn:        setting.TimeForTurn,
		timeForDiscussion:  setting.TimeForDiscussion,
		rolePool:           setting.RolePool,
		factions:           make(map[types.FactionId][]types.PlayerId),
		players:            make(map[types.PlayerId]contract.Player),
		deaths:             make([]types.PlayerId, len(setting.PlayerIds)),
		polls:              make(map[string]*state.Poll),
		round:              state.NewRound(),
	}

	for _, id := range setting.PlayerIds {
		game.players[id] = NewPlayer(&game, id)
	}

	return &game
}

func (g *game) GetCurrentRoundId() types.RoundId {
	return g.round.GetId()
}

func (g *game) GetCurrentRoleId() types.RoleId {
	return g.round.GetCurrentRoleId()
}

func (g *game) GetCurrentPhaseId() types.PhaseId {
	return g.round.GetPhaseId()
}

func (g *game) GetPlayer(playerId types.PlayerId) contract.Player {
	return g.players[playerId]
}

func (g *game) KillPlayer(playerId types.PlayerId) contract.Player {
	if player := g.players[playerId]; player == nil {
		return nil
	} else {
		g.deaths = append(g.deaths, playerId)

		return player
	}
}

func (g *game) RequestAction(playerId types.PlayerId, req *types.ActionRequest) *types.ActionResponse {
	if playerId != req.Actor ||
		slices.Contains(g.deaths, playerId) ||
		!g.round.IsAllowed(playerId) {

		return &types.ActionResponse{
			Error: &types.ErrorDetail{
				Tag: types.UnauthorizedErrorTag,
				Msg: map[string]string{
					types.AlertErrorField: "Not your turn or you're died!",
				},
			},
		}
	}

	return g.GetPlayer(playerId).UseSkill(req)
}
