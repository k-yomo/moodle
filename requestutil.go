package moodle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func getAndUnmarshal(ctx context.Context, client *http.Client, u *url.URL, to interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := mapResponseBodyToStruct(resp.Body, to); err != nil {
		return err
	}
	return nil
}

func mapResponseBodyToStruct(body io.ReadCloser, to interface{}) error {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	fmt.Println("**************************")
	fmt.Println(string(bodyBytes))
	fmt.Println("**************************")

	apiError := APIError{}
	if err := json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.ErrorCode != "" {
		return &apiError
	}

	if err := json.Unmarshal(bodyBytes, to); err != nil {
		return fmt.Errorf("%v, body: %s", err, bodyBytes)
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

func mapIntArrayToQueryParams(key string, ints []int) map[string]string {
	queries := make(map[string]string)
	for i, num := range ints {
		queries[fmt.Sprintf("%s[%d]", key, i)] = strconv.Itoa(num)
	}
	return queries
}

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
