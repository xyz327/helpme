package utils

import (
	"encoding/json"
)

func JsonMarshalToString(src interface{}) string {
	j, err := json.Marshal(src)
	if err != nil {
		return ""
	}
	return string(j)
}
