package action

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"

	"github.com/stretchr/testify/assert"
)

func TestActionID(t *testing.T) {
	id := config.KillActionID
	action := action[any]{
		id: id,
	}

	assert.Equal(t, id, action.ID())
}

func TestActionState(t *testing.T) {
	type strt struct {
		key string
	}

	tests := []struct {
		name          string
		expectedState any
	}{
		{
			name:          "Primitive state",
			expectedState: 0,
		},
		{
			name:          "Slice state",
			expectedState: []int{},
		},
		{
			name:          "Map state",
			expectedState: make(map[string]int),
		},
		{
			name:          "Struct state",
			expectedState: strt{},
		},
		{
			name:          "Reference state",
			expectedState: &strt{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			action := action[any]{
				state: test.expectedState,
			}

			assert.Equal(t, test.expectedState, action.State())
			assert.IsType(t, test.expectedState, action.State())
		})
	}
}

func TestActionValidate(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedErr string
	}{
		{
			name:        "Empty request",
			req:         nil,
			expectedErr: "Action request can not be empty (⊙＿⊙')",
		},
		{
			name: "Non-empty request",
			req:  &types.ActionRequest{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			action := action[any]{}
			err := action.Validate(test.req)

			if test.expectedErr == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestActionSkip(t *testing.T) {
	action := action[any]{}
	expectedRes := &types.ActionResponse{
		Ok:        true,
		IsSkipped: true,
	}

	assert.Equal(t, expectedRes, action.Skip(&types.ActionRequest{}))
}

func TestActionPerform(t *testing.T) {
	action := action[any]{}
	expectedRes := &types.ActionResponse{
		Ok:      false,
		Message: "Nothing to do ¯\\_(ツ)_/¯",
	}

	assert.Equal(t, expectedRes, action.Perform(&types.ActionRequest{}))
}

func TestActionExecute(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.ActionRequest
		expectedRes *types.ActionResponse
	}{
		{
			name: "Skip",
			req: &types.ActionRequest{
				IsSkipped: true,
			},
			expectedRes: &types.ActionResponse{
				Ok:        true,
				IsSkipped: true,
				Data:      nil,
				Message:   "",
			},
		},
		{
			name: "Validation failed",
			req:  nil,
			expectedRes: &types.ActionResponse{
				Ok:        false,
				IsSkipped: false,
				Data:      nil,
				Message:   "Action request can not be empty (⊙＿⊙')",
			},
		},
		{
			name: "Perform",
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
		t.Run(test.name, func(t *testing.T) {
			action := action[any]{}
			res := action.Execute(test.req)

			assert.Equal(t, test.expectedRes, res)
		})
	}
}
