package routes

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"iam-test-framework/config"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// UserProfile represents the profile data required for registration.
type UserProfile struct {
	GivenName  string `json:"givenName"`  // Required field
	FamilyName string `json:"familyName"` // Required field
}

// UserEmail represents the email data required for registration.
type UserEmail struct {
	Email           string `json:"email"`           // Required field
	IsEmailVerified bool   `json:"isEmailVerified"` // Optional field
}

// HumanUser represents the data structure required for Zitadel user registration.
type HumanUser struct {
	Profile UserProfile `json:"profile"`
	Email   UserEmail   `json:"email"`
}

// RegisterUser handles user registration.
func RegisterUser() {
	fmt.Println("NOTE: This feature uses service account and access token")
	fmt.Printf("\n\n")
	// Prompt user to select IAM solution
	var solutionName string
	fmt.Println("Enter the IAM solution (zitadel, keycloak, casdoor):")
	fmt.Scanln(&solutionName)

	// Collect user profile data
	var firstName, lastName, email string
	var isEmailVerified bool

	fmt.Println("Enter first name:")
	fmt.Scanln(&firstName)
	fmt.Println("Enter last name:")
	fmt.Scanln(&lastName)
	fmt.Println("Enter email:")
	fmt.Scanln(&email)
	fmt.Println("Is the email verified? (true/false):")
	fmt.Scanln(&isEmailVerified)

	// Prepare the HumanUser structure
	userData := HumanUser{
		Profile: UserProfile{
			GivenName:  firstName,
			FamilyName: lastName,
		},
		Email: UserEmail{
			Email:           email,
			IsEmailVerified: isEmailVerified,
		},
	}

	// Load IAM configuration
	iamConfig, err := config.GetIAMConfig("config/config.json")
	if err != nil {
		fmt.Printf("Error loading IAM config: %v\n", err)
		return
	}

	solutionConfig, exists := iamConfig.IAMSolutions[solutionName]
	if !exists {
		fmt.Printf("IAM solution %s does not exist.\n", solutionName)
		return
	}

	// Get Client ID and Client Secret
	//clientID, clientSecret, err := GetClientCredentials()
	if err != nil {
		fmt.Printf("Error getting client credentials: %v\n", err)
		return
	}

	// Ensure the domain uses HTTP instead of HTTPS
	domain := solutionConfig.Domain
	if strings.HasPrefix(domain, "https://") {
		domain = strings.Replace(domain, "https://", "http://", 1)
	}
	// Combine domain with the token endpoint
	//tokenEndpoint := fmt.Sprintf("%s%s", domain, solutionConfig.TokenEndpoint) // Full token endpoint

	// Obtain access token
	//accessToken, err := ObtainAccessToken(tokenEndpoint, clientID, clientSecret)
	accessToken, err := GetClientCredentialAccessToken(solutionName)
	if err != nil {
		fmt.Printf("Error obtaining access token: %v\n", err)
		return
	}

	// Prepare the registration URL
	registrationURL := fmt.Sprintf("%s%s", solutionConfig.Domain, solutionConfig.RegistrationEndpoint)

	// Marshal the user data into JSON
	jsonData, err := json.Marshal(userData)
	if err != nil {
		fmt.Printf("Error during registration: %v\n", err)
		return
	}

	// Create a new HTTP request for user registration
	req, err := http.NewRequest("POST", registrationURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set the Content-Type header and Authorization header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the HTTP request to register the user
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error during registration: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Handle the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	if response.StatusCode == http.StatusCreated {
		// Registration successful
		fmt.Println("User registered successfully:", string(body))
	} else {
		// Registration failed
		fmt.Printf("Failed to register user: %s\n", response.Status)
		fmt.Println("Response Body:", string(body))
	}

}

// GetClientCredentials checks for client credentials in environment variables or prompts the user.
func GetClientCredentials() (string, string, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		reader := bufio.NewReader(os.Stdin)

		// If CLIENT_ID is missing, prompt the user
		if clientID == "" {
			fmt.Print("Enter Client ID: ")
			clientID, _ = reader.ReadString('\n')
			clientID = strings.TrimSpace(clientID)
			// Set the CLIENT_ID environment variable
			os.Setenv("CLIENT_ID", clientID)
		}

		// If CLIENT_SECRET is missing, prompt the user
		if clientSecret == "" {
			fmt.Print("Enter Client Secret: ")
			clientSecret, _ = reader.ReadString('\n')
			clientSecret = strings.TrimSpace(clientSecret)
			// Set the CLIENT_SECRET environment variable
			os.Setenv("CLIENT_SECRET", clientSecret)
		}

		// Ensure both values are provided
		if clientID == "" || clientSecret == "" {
			return "", "", errors.New("client ID or client secret missing")
		}
	}

	return clientID, clientSecret, nil
}

// ObtainAccessToken performs Client Credentials Flow to get an access token.
func ObtainAccessToken(tokenEndpoint, clientID, clientSecret string) (string, error) {
	// Prepare the form data for client credentials flow
	data := fmt.Sprintf("grant_type=client_credentials&scope=openid profile urn:zitadel:iam:org:project:id:zitadel:aud")

	// Create a new request with basic authentication
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	// Send the request to get the access token
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check if the token request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get access token: %s", resp.Status)
	}

	// Parse the JSON response to extract the access token
	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", errors.New("access token not found in response")
	}
	fmt.Printf("Service user %s received the accessToken: %s\n", clientID, accessToken)
	return accessToken, nil
}
