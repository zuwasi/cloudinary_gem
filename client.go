package cloudinary_gem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Body       map[string]any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("cloudinary api error: status %d", e.StatusCode)
}

func (c *Cloudinary) doJSON(req *http.Request) (map[string]any, error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return out, &APIError{StatusCode: res.StatusCode, Body: out}
	}
	return out, nil
}
