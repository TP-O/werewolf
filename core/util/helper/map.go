package helper

import "encoding/json"

func JsonToMap(jsonStr string) map[string]any {
	var m map[string]any
	json.Unmarshal([]byte(jsonStr), &m) // nolint errcheck
	return m
}
