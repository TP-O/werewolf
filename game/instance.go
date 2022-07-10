package game

import (
	"errors"
	"math/rand"
	"time"

	"uwwolf/database"
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/model"
	"uwwolf/validator"
)

type instance struct {
	gameId             string
	capacity           uint
	numberOfWerewolves uint
	socketId2playerId  map[string]uint
	playerId2SocketId  map[uint]string
	playerId2RoleId    map[uint]uint
	isStarted          bool
	rolePool           []uint
}

func NewGameInstance(input *contract.GameInstanceInit) (*instance, error) {
	if !validator.ValidateStruct(input) {
		return nil, errors.New("Invalid game instance!")
	}

	gameInstance := instance{
		gameId:             input.GameId,
		capacity:           input.Capacity,
		numberOfWerewolves: input.NumberOfWerewolves,
		isStarted:          false,
		rolePool:           input.RolePool,
	}

	return &gameInstance, nil
}

func (i *instance) Start() bool {
	if i.isStarted || len(i.socketId2playerId) != int(i.capacity) {
		return i.isStarted
	}

	i.assignRoles()

	i.isStarted = true

	return i.isStarted
}

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

func (i *instance) RemoveId(socketId string) bool {
	if i.socketId2playerId[socketId] == 0 {
		return false
	}

	delete(i.playerId2SocketId, i.socketId2playerId[socketId])
	delete(i.socketId2playerId, socketId)

	return true
}

func (i *instance) assignRoles() {
	var roles []model.Role
	database.DB().Find(&roles, i.rolePool)
	randomRoleIds := i.randomRoleIds(roles)

	// Assign roles to players randomly
	for _, playerId := range i.socketId2playerId {
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(randomRoleIds))
		i.playerId2RoleId[playerId] = randomRoleIds[randIndex]

		randomRoleIds = append(randomRoleIds[:randIndex], randomRoleIds[randIndex+1:]...)
	}
}

func (i *instance) randomRoleIds(roles []model.Role) []uint {
	werewolfCounter := uint(0)
	randomRoleIds := []uint{}
	werewolfRoles, remainingRoles := i.splitRoles(roles)

	for len(randomRoleIds) < int(i.capacity) {
		var randIndex int
		var currentRoles *[]model.Role

		// Get random roles belonging to werewolf faction first
		if werewolfCounter < i.numberOfWerewolves {
			currentRoles = &werewolfRoles

			werewolfCounter++
		} else {
			currentRoles = &remainingRoles
		}

		// Put default role to array if no roles are available
		if len(*currentRoles) == 0 {
			if currentRoles == &werewolfRoles {
				randomRoleIds = append(randomRoleIds, enum.WerewolfRole)
			} else {
				randomRoleIds = append(randomRoleIds, enum.VillagerRole)
			}

			continue
		}

		rand.Seed(time.Now().UnixNano())
		randIndex = rand.Intn(len(*currentRoles))
		randomRoleIds = append(randomRoleIds, (*currentRoles)[randIndex].ID)

		(*currentRoles)[randIndex].Quantity--

		// Delete role if its quantity runs out
		if (*currentRoles)[randIndex].Quantity == 0 {
			*currentRoles = append((*currentRoles)[:randIndex], (*currentRoles)[randIndex+1:]...)
		}
	}

	return randomRoleIds
}

// Split roles into 2 factions
func (i *instance) splitRoles(roles []model.Role) ([]model.Role, []model.Role) {
	var werewolfRoles []model.Role
	var remainingRoles []model.Role

	for _, role := range roles {
		if role.TeamID == enum.WerewolfFaction {
			werewolfRoles = append(werewolfRoles, role)
		} else {
			remainingRoles = append(remainingRoles, role)
		}
	}

	return werewolfRoles, remainingRoles
}

func (i *instance) resetPlayers() {
	i.socketId2playerId = make(map[string]uint)
	i.playerId2SocketId = make(map[uint]string)
	i.playerId2RoleId = make(map[uint]uint)
}
