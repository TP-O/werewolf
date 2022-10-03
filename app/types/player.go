package types

type PlayerId string

func (pId PlayerId) IsUnknown() bool {
	return pId == ""
}
