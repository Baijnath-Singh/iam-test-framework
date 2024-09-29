package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iam-test-framework/config"
	"io/ioutil"
	"net/http"
	"os"
)

// LoginUser handles user login.
func LoginUser() {
	var credentials map[string]string
	fmt.Println("Enter your login credentials (in JSON format):")
	json.NewDecoder(os.Stdin).Decode(&credentials)

	iamConfig, err := config.GetIAMConfig("config/config.json") // Fetch the relevant config for the IAM solution
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

	// Prepare the login URL
	loginURL := fmt.Sprintf("%s%s",
		solutionConfig.Domain, solutionConfig.LoginEndpoint)

	// Make the HTTP request to login the user
	jsonData, err := json.Marshal(credentials) // Marshal the credentials into JSON
	if err != nil {
		fmt.Printf("Error during login: %v\n", err)
		return
	}

	response, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData)) // Use bytes.NewBuffer for the POST body
	if err != nil {
		fmt.Printf("Error during login: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Handle response
	if response.StatusCode == http.StatusOK {
		var loginResponse map[string]interface{} // Adjust this according to your response structure
		body, _ := ioutil.ReadAll(response.Body) // Read the body
		json.Unmarshal(body, &loginResponse)     // Unmarshal into the response struct

		fmt.Println("Login successful:", loginResponse)
	} else {
		fmt.Printf("Failed to login: %s\n", response.Status)
	}
}
