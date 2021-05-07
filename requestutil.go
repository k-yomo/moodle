package moodle

import (
	"encoding/json"
	"fmt"
)

func mapResponseBodyToStruct(body []byte, to interface{}) error {

	apiError := APIError{}
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.ErrorCode != "" {
		return &apiError
	}

	if err := json.Unmarshal(body, to); err != nil {
		return fmt.Errorf("%v, body: %s", err, body)
	}

	return nil
}

func mapStrArrayToQueryParams(key string, strs []string) map[string]string {
	queries := make(map[string]string)
	for i, str := range strs {
		queries[fmt.Sprintf("%s[%d]", key, i)] = str
	}
	return queries
}

// func mapIntArrayToQueryParams(key string, ints []int) map[string]string {
// 	queries := make(map[string]string)
// 	for i, num := range ints {
// 		queries[fmt.Sprintf("%s[%d]", key, i)] = strconv.Itoa(num)
// 	}
// 	return queries
// }

// moodle takes bool as 1(true) or 0(false)
func mapBoolToBitStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func mapBitToBool(b int) bool {
	return b == 1
}
