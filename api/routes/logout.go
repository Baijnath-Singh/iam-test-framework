package routes

import (
	"fmt"
	"iam-test-framework/config"
	"net/http"
	"os"
)

// LogoutUser handles user logout.
func LogoutUser() {
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

	// Prepare the logout URL
	logoutURL := fmt.Sprintf("%s%s",
		solutionConfig.Domain, solutionConfig.LogoutEndpoint)

	// Include the necessary authorization token
	token := os.Getenv("AUTH_TOKEN") // Get the token from environment variable or user input
	if token == "" {
		fmt.Println("No authorization token found. Please log in first.")
		return
	}

	req, err := http.NewRequest("GET", logoutURL, nil)
	if err != nil {
		fmt.Printf("Error creating logout request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error during logout: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Handle response
	if response.StatusCode == http.StatusOK {
		fmt.Println("Logout successful.")
	} else {
		fmt.Printf("Failed to logout: %s\n", response.Status)
	}
}
