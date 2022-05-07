package util

import (
	"encoding/json"
)

func ToRecordJson(s interface{}) string {
	str, _ := json.Marshal(s)
	return string(str)
}

