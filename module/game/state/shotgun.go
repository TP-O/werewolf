package state

import (
	"uwwolf/types"
)

type Shotgun struct {
	TargetId types.PlayerId `json:"target"`
}

func NewShotgun() *Shotgun {
	return &Shotgun{}
}

func (s *Shotgun) IsShot() bool {
	return s.TargetId != 0
}

func (s *Shotgun) Shoot(targetId types.PlayerId) bool {
	if targetId.IsUnknown() || s.IsShot() {
		return false
	}

	s.TargetId = targetId

	return true
}
