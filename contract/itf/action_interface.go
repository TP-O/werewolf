package itf

import "uwwolf/contract/typ"

type IAction interface {
	GetName() string
	Perform(game IGame, instruction *typ.ActionInstruction) bool
}
