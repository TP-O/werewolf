package game

import (
	"errors"
	"fmt"
	"sort"

	"golang.org/x/exp/slices"

	"uwwolf/app/model"
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/database"
	"uwwolf/enum"
	"uwwolf/game/factory"
	"uwwolf/util"
	"uwwolf/validator"
)

type turn struct {
	players []string
	role    itf.IRole
}

type instance struct {
	gameId             string
	capacity           uint
	numberOfWerewolves uint
	socketId2playerId  map[string]uint
	playerId2SocketId  map[uint]string
	playerId2RoleId    map[uint]uint
	roleId2SocketIds   map[uint][]string
	isStarted          bool
	rolePool           []uint
	currentPhase       uint
	nextTurn           *turn
	turns              map[uint][]*turn
}

func NewGameInstance(input *typ.GameInstanceInit) (*instance, error) {
	if !validator.ValidateStruct(input) {
		return nil, errors.New("Invalid game instance!")
	}

	gameInstance := instance{
		gameId:             input.GameId,
		capacity:           input.Capacity,
		numberOfWerewolves: input.NumberOfWerewolves,
		isStarted:          false,
		rolePool:           input.RolePool,
		currentPhase:       enum.NightPhase,
		turns:              make(map[uint][]*turn),
	}

	return &gameInstance, nil
}

// Start game instance
func (i *instance) Start() bool {
	if i.isStarted || len(i.socketId2playerId) != int(i.capacity) {
		return i.isStarted
	}

	i.assignRoles()

	i.isStarted = true

	return i.isStarted
}

// Replace players
func (i *instance) AddPlayers(socketIds []string, playerIds []uint) bool {
	if i.isStarted ||
		len(socketIds) != int(i.capacity) ||
		len(playerIds) != int(i.capacity) ||
		len(socketIds) != len(playerIds) ||
		!validator.ValidateVar(socketIds, "unique,dive,len=20,alphanum") ||
		!validator.ValidateVar(playerIds, "unique,dive,min=1") {

		return false
	}

	i.resetPlayers()

	for index, socketId := range socketIds {
		i.socketId2playerId[socketId] = playerIds[index]
		i.playerId2SocketId[playerIds[index]] = socketId
	}

	return true
}

// Remove a player
func (i *instance) RemovePlayer(socketId string) bool {
	if i.socketId2playerId[socketId] == 0 {
		return false
	}

	delete(i.playerId2SocketId, i.socketId2playerId[socketId])
	delete(i.socketId2playerId, socketId)

	return true
}

func (i *instance) Do(instruction *typ.ActionInstruction) bool {
	if !slices.Contains(i.nextTurn.players, instruction.Actor) {
		return false
	}

	i.nextTurn.role.UseSkill(instruction)

	return true
}

func (i *instance) NextTurn() {
	fmt.Println("Next turn!!!!!!!!!!!")
}

// Assign roles to players randomly
func (i *instance) assignRoles() {
	var specialRoles []model.Role
	database.DB().Find(&specialRoles, i.rolePool)

	randomRoles := i.pickUpRoles(specialRoles)
	cloneRandomRoles := slices.Clone(randomRoles)

	for socketId, playerId := range i.socketId2playerId {
		randIndex := util.Intn(len(randomRoles))
		randomRole := randomRoles[randIndex]

		i.playerId2RoleId[playerId] = randomRole.ID
		i.roleId2SocketIds[randomRole.ID] = append(i.roleId2SocketIds[randomRole.ID], socketId)

		if randomRole.ID != enum.VillagerRole {
			i.roleId2SocketIds[enum.VillagerRole] = append(i.roleId2SocketIds[enum.VillagerRole], i.playerId2SocketId[playerId])
		}

		if randomRole.ID != enum.WerewolfRole && randomRole.FactionID == enum.WerewolfFaction {
			i.roleId2SocketIds[enum.WerewolfRole] = append(i.roleId2SocketIds[enum.WerewolfRole], i.playerId2SocketId[playerId])
		}

		randomRoles = append(randomRoles[:randIndex], randomRoles[randIndex+1:]...)
	}

	i.setUpTurns(cloneRandomRoles)
}

// Pick up number of roles corresponding to game capacity randomly
func (i *instance) pickUpRoles(roles []model.Role) []model.Role {
	var defaultRoles []model.Role
	database.DB().Find(&defaultRoles, []uint{enum.VillagerRole, enum.WerewolfRole})

	randomRoles := make([]model.Role, i.capacity)

	werewolfRoles, remainingRoles := i.classifyRoles(roles)

	for j := uint(0); j < i.capacity; j++ {
		var currentRoles *[]model.Role

		if j < i.numberOfWerewolves {
			currentRoles = &werewolfRoles
		} else {
			currentRoles = &remainingRoles
		}

		// Fill in with default role if all special roles are taken
		if len(*currentRoles) == 0 {
			if currentRoles == &werewolfRoles {
				randomRoles[j] = defaultRoles[1]
			} else {
				randomRoles[j] = defaultRoles[0]
			}

			continue
		}

		randIndex := util.Intn(len(*currentRoles))
		randomRoles[j] = (*currentRoles)[randIndex]

		(*currentRoles)[randIndex].Quantity--

		// Delete role if its quantity runs out
		if (*currentRoles)[randIndex].Quantity == 0 {
			*currentRoles = append((*currentRoles)[:randIndex], (*currentRoles)[randIndex+1:]...)
		}
	}

	return randomRoles
}

// Classify roles into 2 factions
func (i *instance) classifyRoles(roles []model.Role) ([]model.Role, []model.Role) {
	var werewolfRoles []model.Role
	var remainingRoles []model.Role

	for _, role := range roles {
		if role.FactionID == enum.WerewolfFaction {
			werewolfRoles = append(werewolfRoles, role)
		} else {
			remainingRoles = append(remainingRoles, role)
		}
	}

	return werewolfRoles, remainingRoles
}

// Reset player-related data
func (i *instance) resetPlayers() {
	i.socketId2playerId = make(map[string]uint)
	i.playerId2SocketId = make(map[uint]string)
	i.playerId2RoleId = make(map[uint]uint)
	i.roleId2SocketIds = make(map[uint][]string)
}

// Prepare turn for special roles
func (i *instance) setUpTurns(roles []model.Role) {
	roleFactory := factory.GetRoleFactory()
	roles = util.RemoveDuplicate(roles)

	// Order by priority
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Priority < roles[j].Priority
	})

	for _, role := range roles {
		i.turns[role.PhaseID] = append(i.turns[role.PhaseID], &turn{
			players: i.roleId2SocketIds[role.ID],
			role:    roleFactory.Create(role.ID, i),
		})
	}

	// Test
	i.turns[enum.NightPhase][0].players = append(i.turns[enum.NightPhase][0].players, "11111111111111111111", "11111111111111111113")

	for i, phases := range i.turns {
		fmt.Println("Phase: ", i)

		for _, turn := range phases {
			fmt.Println(turn.players)
		}
	}

	i.nextTurn = i.turns[i.currentPhase][0]
}
