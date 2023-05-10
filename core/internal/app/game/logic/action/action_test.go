package action

import (
	"errors"
	"testing"
	"uwwolf/internal/app/game/logic/types"

	"github.com/stretchr/testify/suite"
)

type ActionSuite struct {
	suite.Suite
	id types.ActionId
}

func TestActionSuiteSuite(t *testing.T) {
	suite.Run(t, new(ActionSuite))
}
func (as *ActionSuite) SetupSuite() {
	as.id = 1
}

func (as ActionSuite) TestId() {
	act := action{
		id: as.id,
	}

	as.Equal(as.id, act.Id())
}

func (as ActionSuite) TestValidate() {
	act := action{}
	err := act.validate(&types.ActionRequest{})

	as.Equal(errors.New("Validation is required!"), err)
}

func (as ActionSuite) TestSkip() {
	act := action{}
	expectedRes := types.ActionResponse{
		Ok: true,
		ActionRequest: types.ActionRequest{
			IsSkipped: true,
		},
		Message: "Skipped!",
	}

	as.Equal(expectedRes, act.skip(&types.ActionRequest{}))
}

func (as ActionSuite) TestPerform() {
	act := action{}
	expectedRes := types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}

	as.Equal(expectedRes, act.perform(&types.ActionRequest{}))
}

func (as ActionSuite) TestExecuteSchema() {

	tests := []struct {
		name         string
		skipValidate bool
		req          *types.ActionRequest
		expectedRes  types.ActionResponse
	}{
		{
			name: "Failure (Nil request)",
			req:  nil,
			expectedRes: types.ActionResponse{
				Ok:       false,
				ActionId: as.id,
				Message:  "Action request can not be empty (⊙＿⊙')",
			},
		},
		{
			name: "Ok (Skip)",
			req: &types.ActionRequest{
				IsSkipped: true,
			},
			expectedRes: types.ActionResponse{
				Ok:       true,
				ActionId: as.id,
				ActionRequest: types.ActionRequest{
					IsSkipped: true,
				},
				Message: "Skipped!",
			},
		},
		{
			name: "Failure (Validation failed)",
			req:  &types.ActionRequest{},
			expectedRes: types.ActionResponse{
				Ok:       false,
				ActionId: as.id,
				Message:  "Validation is required!",
			},
		},
		{
			name:         "Ok (Perform)",
			skipValidate: true,
			req:          &types.ActionRequest{},
			expectedRes: types.ActionResponse{
				Ok:            false,
				ActionId:      as.id,
				ActionRequest: types.ActionRequest{},
				Data:          nil,
				Message:       "Nothing to do ¯\\_(ツ)_/¯",
			},
		},
	}

	for _, test := range tests {
		as.Run(test.name, func() {
			act := action{
				skipValidate: test.skipValidate,
			}
			res := act.execute(act, as.id, test.req)

			as.Equal(test.expectedRes, res)
		})
	}
}

func (as ActionSuite) TestExecute() {

	tests := []struct {
		name         string
		skipValidate bool
		req          types.ActionRequest
		expectedRes  types.ActionResponse
	}{
		{
			name: "Ok (Skip)",
			req: types.ActionRequest{
				IsSkipped: true,
			},
			expectedRes: types.ActionResponse{
				Ok:       true,
				ActionId: as.id,
				ActionRequest: types.ActionRequest{
					IsSkipped: true,
				},
				Message: "Skipped!",
			},
		},
		{
			name: "Failure (Validation failed)",
			req:  types.ActionRequest{},
			expectedRes: types.ActionResponse{
				Ok:       false,
				ActionId: as.id,
				Message:  "Validation is required!",
			},
		},
		{
			name:         "Ok (Perform)",
			skipValidate: true,
			req:          types.ActionRequest{},
			expectedRes: types.ActionResponse{
				Ok:            false,
				ActionId:      as.id,
				ActionRequest: types.ActionRequest{},
				Data:          nil,
				Message:       "Nothing to do ¯\\_(ツ)_/¯",
			},
		},
	}

	for _, test := range tests {
		as.Run(test.name, func() {
			act := action{
				id:           as.id,
				skipValidate: test.skipValidate,
			}
			res := act.Execute(test.req)

			as.Equal(test.expectedRes, res)
		})
	}
}
