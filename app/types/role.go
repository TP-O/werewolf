package types

type RoleId uint

type RoleActionSetting struct {
	ActionId   ActionId
	Expiration Expiration
	Payload    any
}
