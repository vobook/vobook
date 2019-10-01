package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StructToQueryString converts flat struct to query string
func StructToQueryString(st interface{}) (qs string, err error) {
	jsonBytes, err := json.Marshal(st)
	if err != nil {
		return
	}

	convert := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &convert)
	if err != nil {
		return
	}

	qsParams := make([]string, len(convert))
	i := 0
	for key, val := range convert {
		qsParams[i] = fmt.Sprintf("%s=%v", key, val)
		i++
	}

	qs = strings.Join(qsParams, "&")
	return
}
