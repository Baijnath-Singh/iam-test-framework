package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iam-test-framework/config"
	"io/ioutil"
	"net/http"
)

// UserInfoResponse represents the structure of a single user info response.
type UserInfoResponse struct {
	UserID     string   `json:"userId"`
	Username   string   `json:"username"`
	LoginNames []string `json:"loginNames"`
	State      string   `json:"state"`
}

// UserListResponse represents the structure of the list of user info responses.
type UserListResponse struct {
	Result []UserInfoResponse `json:"result"`
}

// GetUserInfo retrieves user information.
func GetUserInfo() {
	fmt.Println("NOTE: This feature uses service account and access token")
	fmt.Printf("\n\n")
	iamConfig, err := config.GetIAMConfig("config/config.json")
	if err != nil {
		fmt.Printf("Error loading IAM config: %v\n", err)
		return
	}

	// Prompt user to select IAM solution
	var solutionName string
	fmt.Println("Enter the IAM solution (zitadel, keycloak, casdoor):")
	fmt.Scanln(&solutionName)

	solutionConfig, exists := iamConfig.IAMSolutions[solutionName]
	if !exists {
		fmt.Printf("IAM solution %s does not exist.\n", solutionName)
		return
	}

	// Check if the necessary endpoints are configured
	if solutionConfig.Domain == "" || solutionConfig.UserinfoEndpoint == "" {
		fmt.Println("Error: Domain and UserinfoEndpoint must be configured.")
		return
	}

	// Prepare the user info URL
	userInfoURL := fmt.Sprintf("%s%s", solutionConfig.Domain, solutionConfig.UserinfoEndpoint)
	fmt.Printf("User Info URL: %s\n", userInfoURL)

	// If token is not present, use the function to retrieve it
	token, err := GetClientCredentialAccessToken(solutionName)
	if err != nil {
		fmt.Printf("Error retrieving access token: %v\n", err)
		return
	}

	// Create the request body
	body := map[string]interface{}{
		"offset": 0,
		"limit":  10, // Adjust the limit as needed
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error marshalling request body: %v\n", err)
		return
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", userInfoURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Printf("Error creating user info request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error retrieving user info: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Read response body as []byte
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	if response.StatusCode == http.StatusOK {
		var userListResponse UserListResponse // Expecting an object that contains an array of users
		if err := json.Unmarshal(responseBody, &userListResponse); err != nil {
			fmt.Printf("Error decoding user info response: %v\n", err)
			return
		}

		// Display user info
		for _, user := range userListResponse.Result {
			fmt.Printf("User retrieved successfully. ID: %s, Username: %s, State: %s\n", user.UserID, user.Username, user.State)
		}
	} else {
		fmt.Printf("Failed to retrieve user info: %s\n", response.Status)
		fmt.Printf("Response body: %s\n", string(responseBody))
	}
}
