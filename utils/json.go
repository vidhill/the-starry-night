package utils

import "encoding/json"

// For structs manually created so should never fail to marshal into JSON, therefore can safely ignore error
func MarshalIgnoreError(i interface{}) []byte {
	b, _ := json.Marshal(i)
	return b
}

// marshall indented, useful for debug output
func MarshalIndentIgnoreError(i interface{}) []byte {
	b, _ := json.MarshalIndent(i, "", "    ")
	return b
}
