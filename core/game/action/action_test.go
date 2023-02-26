package action

// import (
// 	"testing"
// 	"uwwolf/game/enum"
// 	"uwwolf/game/types"

// 	"github.com/stretchr/testify/suite"
// )

// type ActionSuite struct {
// 	suite.Suite
// }

// func TestActionSuiteSuite(t *testing.T) {
// 	suite.Run(t, new(ActionSuite))
// }

// func (as *ActionSuite) TestActionID() {
// 	id := enum.KillActionID
// 	action := action[any]{
// 		id: id,
// 	}

// 	as.Equal(id, action.ID())
// }

// func (as *ActionSuite) TestActionState() {
// 	type strt struct {
// 		key string
// 	}

// 	tests := []struct {
// 		name          string
// 		expectedState any
// 	}{
// 		{
// 			name:          "Primitive state",
// 			expectedState: 0,
// 		},
// 		{
// 			name:          "Slice state",
// 			expectedState: []int{},
// 		},
// 		{
// 			name:          "Map state",
// 			expectedState: make(map[string]int),
// 		},
// 		{
// 			name:          "Struct state",
// 			expectedState: strt{},
// 		},
// 		{
// 			name:          "Reference state",
// 			expectedState: &strt{},
// 		},
// 	}

// 	for _, test := range tests {
// 		as.Run(test.name, func() {
// 			action := action[any]{
// 				state: test.expectedState,
// 			}

// 			as.Equal(test.expectedState, action.State())
// 			as.IsType(test.expectedState, action.State())
// 		})
// 	}
// }

// func (as *ActionSuite) TestActionValidate() {
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedErr string
// 	}{
// 		{
// 			name:        "Failure (Empty request)",
// 			req:         nil,
// 			expectedErr: "Action request can not be empty (⊙＿⊙')",
// 		},
// 		{
// 			name: "Ok (Non-empty request)",
// 			req:  &types.ActionRequest{},
// 		},
// 	}

// 	for _, test := range tests {
// 		as.Run(test.name, func() {
// 			action := action[any]{}
// 			err := action.Validate(test.req)

// 			if test.expectedErr == "" {
// 				as.Nil(err)
// 			} else {
// 				as.Equal(test.expectedErr, err.Error())
// 			}
// 		})
// 	}
// }

// func (as *ActionSuite) TestActionSkip() {
// 	action := action[any]{}
// 	expectedRes := &types.ActionResponse{
// 		Ok:        true,
// 		IsSkipped: true,
// 	}

// 	as.Equal(expectedRes, action.Skip(&types.ActionRequest{}))
// }

// func (as *ActionSuite) TestActionPerform() {
// 	action := action[any]{}
// 	expectedRes := &types.ActionResponse{
// 		Ok:      false,
// 		Message: "Nothing to do ¯\\_(ツ)_/¯",
// 	}

// 	as.Equal(expectedRes, action.Perform(&types.ActionRequest{}))
// }

// func (as *ActionSuite) TestActionExecute() {
// 	tests := []struct {
// 		name        string
// 		req         *types.ActionRequest
// 		expectedRes *types.ActionResponse
// 	}{
// 		{
// 			name: "Ok (Skip)",
// 			req: &types.ActionRequest{
// 				IsSkipped: true,
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        true,
// 				IsSkipped: true,
// 				Data:      nil,
// 				Message:   "",
// 			},
// 		},
// 		{
// 			name: "Failure (Validation failed)",
// 			req:  nil,
// 			expectedRes: &types.ActionResponse{
// 				Ok:        false,
// 				IsSkipped: false,
// 				Data:      nil,
// 				Message:   "Action request can not be empty (⊙＿⊙')",
// 			},
// 		},
// 		{
// 			name: "Ok (Perform)",
// 			req:  &types.ActionRequest{
// 				//
// 			},
// 			expectedRes: &types.ActionResponse{
// 				Ok:        false,
// 				IsSkipped: false,
// 				Data:      nil,
// 				Message:   "Nothing to do ¯\\_(ツ)_/¯",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		as.Run(test.name, func() {
// 			action := action[any]{}
// 			res := action.Execute(test.req)

// 			as.Equal(test.expectedRes, res)
// 		})
// 	}
// }
