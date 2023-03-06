package action

import (
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"

	"github.com/stretchr/testify/suite"
)

type ActionSuite struct {
	suite.Suite
}

func TestActionSuiteSuite(t *testing.T) {
	suite.Run(t, new(ActionSuite))
}

func (as ActionSuite) TestActionID() {
	id := vars.KillActionID
	act := action{
		id: id,
	}

	as.Equal(id, act.ID())
}

func (as ActionSuite) TestActionValidate() {
	act := action{}
	err := act.validate(&types.ActionRequest{})

	as.Nil(err)
}

func (as ActionSuite) TestActionSkip() {
	act := action{}
	expectedRes := &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
		Message:   "Skipped!",
	}

	as.Equal(expectedRes, act.skip(&types.ActionRequest{}))
}

func (as ActionSuite) TestActionPerform() {
	act := action{}
	expectedRes := &types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}

	as.Equal(expectedRes, act.perform(&types.ActionRequest{}))
}

func (as ActionSuite) TestActionExecute() {

	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
	}{
		{
			name: "Failure (Nil request)",
			req:  nil,
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Action request can not be empty (⊙＿⊙')",
			},
		},
		{
			name: "Ok (Skip)",
			req: &types.ActionRequest{
				IsSkipped: true,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: true,
				Message:   "Skipped!",
			},
		},
		{
			name: "Failure (Validation failed)",
			req:  nil,
			expectedRes: &types.ActionResponse{
				Ok:      false,
				Message: "Action request can not be empty (⊙＿⊙')",
			},
		},
		{
			name: "Ok (Perform)",
			req:  &types.ActionRequest{
				//
			},
			expectedRes: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "Nothing to do ¯\\_(ツ)_/¯",
			},
		},
	}

	for _, test := range tests {
		as.Run(test.name, func() {
			act := action{}
			res := act.execute(act, test.req)

			as.Equal(test.expectedRes, res)
		})
	}
}
