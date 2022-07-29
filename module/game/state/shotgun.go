package state

import "uwwolf/types"

type Shotgun struct {
	isShot bool
	target types.PlayerId
}

func NewShotgun() Shotgun {
	return Shotgun{}
}

func (s *Shotgun) IsShot() bool {
	return s.target != 0
}

func (s *Shotgun) GetTarget() types.PlayerId {
	return s.target
}

func (s *Shotgun) Shoot(target types.PlayerId) bool {
	if target == 0 || s.target != 0 {
		return false
	}

	s.target = target

	return true
}
