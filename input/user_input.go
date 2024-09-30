package input

import (
	"bufio"
	"fmt"
	"os"
)

// GetUserChoice prompts the user to select an IAM functionality to test.
func GetUserChoice() string {
	var choice string
	fmt.Print("Enter your choice (1-4): ")
	fmt.Scanln(&choice) // Read user input for choice
	return choice
}

// Function to get access token (placeholder)
func GetAccessToken() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your access token: ")
	token, _ := reader.ReadString('\n')
	return token
}

// GetUserInput prompts the user for input values and updates the configuration accordingly.
func GetUserEmail() {
	// User prompts can be added here for other configurations
	fmt.Print("Enter your user email for registration: ")
	var email string
	fmt.Scanln(&email)

	// Further prompts for password, etc.
	// These values can be returned or used directly to pass to the respective functions in the API routes.
}

// Helper function to read user input
func GetUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1] // remove newline character
}
