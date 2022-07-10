package contract

type ActionInstruction struct {
	GameId  string   `validate:"required,len=20,alphanum"`
	Actor   string   `validate:"required,len=20,alphanum"`
	Targets []string `validate:"required_if=Skipped false,min=1,dive,len=20,alphanum"`
	Skipped bool     `validate:""`
	Payload []int    `validate:""`
}

type Action interface {
	GetName() string
	Perform(instruction ActionInstruction) bool
}
