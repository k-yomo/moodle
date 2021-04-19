package moodleapi

import (
	"encoding/json"
	"io"
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
	err := json.NewDecoder(body).Decode(to)
	if err == nil {
		return nil, nil
	}

	apiError := APIError{}
	if err := json.NewDecoder(body).Decode(&apiError); err == nil {
		return &apiError, nil
	}

	return nil, err
}
