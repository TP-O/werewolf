package types

type FactionID uint8

func (f FactionID) IsUnknown() bool {
	return f == 0
}

type RoleID uint8

func (r RoleID) IsUnknown() bool {
	return r == 0
}

type ActivateAbilityRequest struct {
	AbilityIndex uint8
	TargetID     PlayerID
	IsSkipped    bool
}

type Limit int8
