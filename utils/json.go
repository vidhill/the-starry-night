package utils

import "encoding/json"

func MarshalIgnoreError(i interface{}) []byte {
	b, _ := json.Marshal(i)
	return b
}
