package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchData[T any](ctx context.Context, url string) ([]T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request %s failed: %v", url, err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %v", url, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body for request %s. Error: %v", url, err)
	}

	var results []T
	if err = json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("unmarshalling response body for request %s failed: %v", url, err)
	}
	return results, nil
}
