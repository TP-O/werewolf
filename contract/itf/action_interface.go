package itf

import "uwwolf/contract/typ"

type IAction interface {
	GetName() string
	Perform(instruction *typ.ActionInstruction) bool
}
