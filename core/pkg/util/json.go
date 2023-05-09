package util

import "encoding/json"

func JsonMarshal(src any, dest *string) error {
	if b, err := json.Marshal(src); err != nil {
		return err
	} else {
		*dest = string(b)
		return nil
	}
}

func JsonUnmarshal(src string, dest any) error {
	if err := json.Unmarshal([]byte(src), dest); err != nil {
		return err
	} else {
		return nil
	}
}
