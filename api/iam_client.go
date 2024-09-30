package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IAMClient is a struct that holds the base URL and HTTP client.
type IAMClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewIAMClient initializes a new IAMClient with the given base URL.
func NewIAMClient(baseURL string) *IAMClient {
	return &IAMClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// PostRequest makes a POST request to the given endpoint with the provided payload.
func (client *IAMClient) PostRequest(endpoint string, payload interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", client.BaseURL, endpoint)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %v", err)
	}

	return resp, nil
}

// GetRequest makes a GET request to the given endpoint.
func (client *IAMClient) GetRequest(endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", client.BaseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}

	return resp, nil
}

// ReadResponseBody reads the response body and returns it as a string.
func ReadResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close() // Ensure the response body is closed after reading
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	return string(body), nil
}
