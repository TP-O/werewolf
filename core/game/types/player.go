package types

type PlayerID string

func (p PlayerID) IsUnknown() bool {
	return p == ""
}
