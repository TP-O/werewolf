package factory

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

var roleFactoryInstance *roleFactory

func init() {
	roleFactoryInstance = &roleFactory{}
}

func Role(id types.RoleId, game contract.Game, setting *types.RoleSetting) contract.Role {
	return roleFactoryInstance.Create(id, game, setting)
}
