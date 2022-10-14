package types

type ActionId = uint

type ActionRequest struct {
	GameId    GameId     `validate:"required,number,gt=0"`
	ActionId  ActionId   `validate:"required,number,gt=0"`
	ActorId   PlayerId   `validate:"required,len=28"`
	TargetIds []PlayerId `validate:"required,min=1,unique,dive,len=28"`
	IsSkipped bool       ``
}

type ActionResponse struct {
	Ok              bool
	Data            any
	PerformError    string
	ValidationError ValidationError
}

type VoteActionSetting struct {
	FactionId FactionId
	PlayerId  PlayerId
	Weight    uint
}
