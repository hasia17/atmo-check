package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadResponse[T any](url string) ([]T, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %v", url, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body for request %u. Error: %v", url, err)
	}

	var results []T
	if err = json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("unmarshalling response body for request %u failed: %v", url, err)
	}
	return results, nil
}
