package types

type Limit int

type Priority int

type RoleID uint

func (r RoleID) IsUnknown() bool {
	return r == 0
}

type FactionID uint

func (f FactionID) IsUnknown() bool {
	return f == 0
}

type UseRoleRequest struct {
	ActionID  ActionID
	TargetIDs []PlayerID
	IsSkipped bool
}
