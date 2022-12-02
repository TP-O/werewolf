package enum

type PlayerStatus uint

const (
	OnlineStatus PlayerStatus = iota + 1
	BusyStatus
	InGameStatus
)

type PlayerID string

func (p PlayerID) IsUnknown() bool {
	return p == ""
}
