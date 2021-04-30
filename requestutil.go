package moodle

import (
	"encoding/json"
	"errors"
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

	apiError, err := mapResponseBodyToStruct(resp.Body, to)
	if apiError != nil {
		return apiError
	} else if err != nil {
		return err
	}
	return nil
}

func mapResponseBodyToStruct(body io.ReadCloser, to interface{}) (*APIError, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, to)
	if err == nil {
		return nil, nil
	}

	apiError := APIError{}
	if err := json.Unmarshal(bodyBytes, &apiError); err == nil {
		return &apiError, nil
	}

	return nil, errors.New(fmt.Sprintf("%v, body: %s", err, bodyBytes))
}

func strArrayToQueryParams(key string, strs []string) map[string]string {
	queries := make(map[string]string)
	for i, str := range strs {
		queries[fmt.Sprintf("%s[%d]", key, i)] = str
	}
	return queries
}
