package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// APIError defines the structure for API error responses.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	return fmt.Sprintf("APIError %d: %s", e.Code, e.Message)
}

// HandleError is a utility function to handle errors in API responses.
func HandleError(w http.ResponseWriter, err error, statusCode int) {
	// Log the error for debugging purposes
	log.Printf("Error: %v", err)

	// Set the status code
	w.WriteHeader(statusCode)

	// Create an API error response
	apiError := APIError{
		Code:    statusCode,
		Message: err.Error(),
	}

	// Send the error response as JSON
	if jsonErr := json.NewEncoder(w).Encode(apiError); jsonErr != nil {
		log.Printf("Failed to encode error response: %v", jsonErr)
	}
}

// HandleResponse checks the response for errors and handles them appropriately.
func HandleResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return err
		}
		return &APIError{
			Code:    apiError.Code,
			Message: apiError.Message,
		}
	}
	return nil
}
