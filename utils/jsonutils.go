package utils

import (
	"encoding/json"
)

func SafeMarshal(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(jsonData)
}
