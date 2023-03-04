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
	TargetID     PlayerID `json:"target_id" validate:"required,min=1,unique,dive,len=20"`
	IsSkipped    bool
}

type Limit int8
