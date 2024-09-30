package routes

import (
	"encoding/json"
	"fmt"
	"iam-test-framework/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// TokenResponse represents the structure of the response when retrieving a token.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	// Add other fields as necessary
}

// GetToken retrieves an access token.
func GetToken() {
	// Prompt user to select IAM solution
	var solutionName string
	fmt.Println("NOTE: As of now this feature only generates Client Credential Access Token")
	fmt.Printf("\n\n")
	fmt.Println("Enter the IAM solution (zitadel, keycloak, casdoor):")
	fmt.Scanln(&solutionName)
	token, err := GetClientCredentialAccessToken("zitadel")
	if err != nil {
		fmt.Printf("Error retrieving access token: %v\n", err)
	} else {
		fmt.Printf("Access Token: %s\n", token)
	}
}

// GetClientCredentialAccessToken retrieves an access token using the client credentials grant type.
func GetClientCredentialAccessToken(solutionName string) (string, error) {
	var tokenRequest map[string]string
	iamConfig, err := config.GetIAMConfig("config/config.json")
	if err != nil {
		return "", fmt.Errorf("error loading IAM config: %w", err)
	}

	solutionConfig, exists := iamConfig.IAMSolutions[solutionName]
	if !exists {
		return "", fmt.Errorf("IAM solution %s does not exist", solutionName)
	}

	// Set the grant type to client_credentials
	grantType := "client_credentials"

	// Initialize tokenRequest map
	tokenRequest = make(map[string]string)

	// Check if parameters are pre-configured in the config.json for this grant type
	if params, ok := solutionConfig.TokenRequestParams[grantType]; ok {
		ensureParam("client_id", params, &tokenRequest)
		ensureParam("client_secret", params, &tokenRequest)
		ensureParam("scope", params, &tokenRequest)
	} else {
		return "", fmt.Errorf("no pre-configured parameters found for this grant type. Please enter them manually")
	}

	// Add the grant type to the request
	tokenRequest["grant_type"] = grantType

	// Prepare the token URL
	if solutionConfig.Domain == "" || solutionConfig.TokenEndpoint == "" {
		return "", fmt.Errorf("domain and token endpoint must be configured")
	}

	tokenURL := fmt.Sprintf("%s%s", solutionConfig.Domain, solutionConfig.TokenEndpoint)

	// Encode the parameters
	form := url.Values{}
	for k, v := range tokenRequest {
		form.Set(k, v)
	}

	// Make the HTTP POST request with URL-encoded data
	response, err := http.PostForm(tokenURL, form)
	if err != nil {
		return "", fmt.Errorf("error during token retrieval: %w", err)
	}
	defer response.Body.Close()

	// Handle response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if response.StatusCode == http.StatusOK {
		var tokenResponse TokenResponse
		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			return "", fmt.Errorf("error decoding token response: %w", err)
		}
		return tokenResponse.AccessToken, nil
	}

	return "", fmt.Errorf("token retrieval failed: %s, response body: %s", response.Status, string(body))
}

// ensureParam checks if a parameter is present in the config map, and if not, prompts the user to input it.
func ensureParam(param string, configParams map[string]string, tokenRequest *map[string]string) {
	if val, ok := configParams[param]; ok && val != "" {
		(*tokenRequest)[param] = val
	} else {
		reader := os.Stdin
		readerBuffer := make([]byte, 100)

		fmt.Printf("Enter %s: ", param)
		n, _ := reader.Read(readerBuffer)
		(*tokenRequest)[param] = strings.TrimSpace(string(readerBuffer[:n]))
	}
}
