package main

import (
	"fmt"
	"iam-test-framework/api/routes" // Change to your actual module name
	"iam-test-framework/config"     // Change to your actual module name
	"iam-test-framework/input"      // Change to your actual module name
)

func main() {
	// Load configuration
	_, err := config.LoadConfig("config/config.json")
	if err != nil {
		fmt.Printf("Error loading IAM config: %v\n", err)
		return
	}

	// Set up logging - Provide the path to the log file
	//logger.NewLogger("logger/app.log") // Adjust the path as necessary

	// User input for selecting the IAM functionality to test
	fmt.Println("Select the IAM functionality to test:")
	fmt.Println("1. User Registration")
	fmt.Println("2. OIDC Authorization")
	fmt.Println("3. Token Issuance")
	fmt.Println("4. User Info Retrieval")

	choice := input.GetUserChoice() // Fetch user input

	switch choice {
	case "1":
		routes.RegisterUser()
	case "2":
		routes.OIDCAuthorization()
	case "3":
		routes.GetToken()
	case "4":
		routes.GetUserInfo()
	default:
		fmt.Println("Invalid choice. Please select a valid option.")
	}
}
