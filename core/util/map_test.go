package util_test

import (
	"testing"
	"uwwolf/util"

	"github.com/stretchr/testify/require"
)

func TestJsonToMap(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectedMap map[string]any
	}{
		{
			name:        "Invalid json",
			json:        "{Invalid JSON}",
			expectedMap: nil,
		},
		{
			name: "Ok",
			json: `{
                "f1": 1,
                "f2": "2",
                "f3": true,
                "f4": [1,2],
                "f5": {
                    "f5.1": 1
                }
            }`,
			expectedMap: map[string]any{
				"f1": 1.0,
				"f2": "2",
				"f3": true,
				"f4": []interface{}{1.0, 2.0},
				"f5": map[string]interface{}{
					"f5.1": 1.0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := util.JsonToMap(test.json)

			require.Equal(t, test.expectedMap, m)
		})
	}
}
