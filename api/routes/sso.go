package routes

import (
	"fmt"
	"iam-test-framework/config"
)

// SSO handles Single Sign-On.
func SSO() {
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

	// Prepare the SSO URL
	ssoURL := fmt.Sprintf("%s%s", solutionConfig.Domain, solutionConfig.AuthorizeEndpoint)

	// Redirect the user to the SSO URL
	fmt.Printf("Redirecting to SSO at:\n%s\n", ssoURL)
}
