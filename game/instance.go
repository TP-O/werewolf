package game

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"golang.org/x/exp/slices"

	"uwwolf/app/model"
	"uwwolf/contract/typ"
	"uwwolf/database"
	"uwwolf/enum"
	"uwwolf/game/factory"
	"uwwolf/game/stuff"
	"uwwolf/util"
	"uwwolf/validator"
)

type instance struct {
	gameId              string
	capacity            uint
	numberOfWerewolves  uint
	remaining           uint
	remainingWerewolves uint
	socketId2playerId   map[string]uint
	playerId2SocketId   map[uint]string
	playerId2RoleId     map[uint]uint
	roleId2SocketIds    map[uint][]string
	roleId2PlayerIds    map[uint][]uint
	isStarted           bool
	rolePool            []uint
	data                chan string
	phase               *stuff.Phase
}

func NewGameInstance(input *typ.GameInstanceInit) (*instance, error) {
	if !validator.ValidateStruct(input) {
		return nil, errors.New("Invalid game instance!")
	}

	gameInstance := instance{
		gameId:              input.GameId,
		capacity:            input.Capacity,
		numberOfWerewolves:  input.NumberOfWerewolves,
		remaining:           input.Capacity,
		remainingWerewolves: input.NumberOfWerewolves,
		isStarted:           false,
		rolePool:            input.RolePool,
		phase:               &stuff.Phase{},
		data:                make(chan string),
	}

	gameInstance.phase.Init()

	return &gameInstance, nil
}

func (i *instance) IsStarted() bool {
	return i.isStarted
}

func (i *instance) NumberOfVillagers() uint {
	return i.remaining
}

func (i *instance) NumberOfWerewolves() uint {
	return i.remainingWerewolves
}

// Start game instance
func (i *instance) Start() bool {
	if i.isStarted || len(i.socketId2playerId) != int(i.capacity) {
		return false
	}

	i.isStarted = true

	i.assignRoles()
	go i.listen()
	go i.phase.Start()

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

func (i *instance) NextTurn() {
	go i.phase.NextTurn()
}

func (i *instance) Pipe(pub *chan string) {
	*pub = i.data
}

func (i *instance) Do(instruction *typ.ActionInstruction) bool {
	if !i.phase.IsValidPlayer(instruction.Actor) {
		return false
	}

	return i.phase.UseSkill(instruction)
}

func (i *instance) listen() {
	for d := range i.data {
		fmt.Println(d)
	}
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
		i.roleId2PlayerIds[randomRole.ID] = append(i.roleId2PlayerIds[randomRole.ID], playerId)

		if randomRole.ID != enum.VillagerRole {
			i.roleId2SocketIds[enum.VillagerRole] = append(i.roleId2SocketIds[enum.VillagerRole], socketId)
			i.roleId2PlayerIds[enum.VillagerRole] = append(i.roleId2PlayerIds[enum.VillagerRole], playerId)
		}

		if randomRole.ID != enum.WerewolfRole && randomRole.FactionID == enum.WerewolfFaction {
			i.roleId2SocketIds[enum.WerewolfRole] = append(i.roleId2SocketIds[enum.WerewolfRole], socketId)
			i.roleId2PlayerIds[enum.WerewolfRole] = append(i.roleId2PlayerIds[enum.WerewolfRole], playerId)
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
	i.roleId2PlayerIds = make(map[uint][]uint)
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
		i.phase.AddTurn(
			role.PhaseID,
			2*time.Second,
			roleFactory.Create(role.ID, i),
			i.roleId2PlayerIds[role.ID],
		)
	}

	// Test
	i.phase.AddPlayer(enum.NightPhase, 0, 1, 3)

	for i, phases := range i.phase.Data() {
		fmt.Println("Phase: ", i)

		for _, turn := range phases {
			fmt.Println(turn)
		}
	}
}
