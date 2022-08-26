package types

type PlayerId uint

func (pId PlayerId) IsUnknown() bool {
	return pId == UnknownPlayer
}

const (
	UnknownPlayer PlayerId = iota
)
