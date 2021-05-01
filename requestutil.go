package moodle

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func doAndMap(client *http.Client, req *http.Request, to interface{}) error {
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

	apiError := APIError{}
	if err := json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.ErrorCode != "" {
		return &apiError
	}

	if err := json.Unmarshal(bodyBytes, to); err != nil {
		return fmt.Errorf("%v, body: %s", err, bodyBytes)
	}

	return nil
}

func strArrayToQueryParams(key string, strs []string) map[string]string {
	queries := make(map[string]string)
	for i, str := range strs {
		queries[fmt.Sprintf("%s[%d]", key, i)] = str
	}
	return queries
}
