package types

type FactionId uint

func (f FactionId) IsUnknown() bool {
	return f == 0
}
