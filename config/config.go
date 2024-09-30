package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// IAMConfig holds the configuration for IAM solutions.
type IAMConfig struct {
	IAMSolutions    map[string]SolutionConfig `json:"iam_solutions"`
	DefaultSolution string                    `json:"default_solution"` // Field for default solution
}

// SolutionConfig holds the configuration for a specific IAM solution.
type SolutionConfig struct {
	Domain               string                       `json:"domain"`
	RegistrationEndpoint string                       `json:"registration_endpoint"`
	LoginEndpoint        string                       `json:"login_endpoint"`
	AuthorizeEndpoint    string                       `json:"authorize_endpoint"`
	TokenEndpoint        string                       `json:"token_endpoint"`
	UserinfoEndpoint     string                       `json:"userinfo_endpoint"`
	LogoutEndpoint       string                       `json:"logout_endpoint"`
	AppClientID          string                       `json:"app_client_id"`
	AppScope             string                       `json:"app_scope"`
	TokenRequestParams   map[string]map[string]string `json:"token_request_params,omitempty"` // Field for token request params
}

// LoadConfig loads the configuration from the JSON file and allows for user overrides.
func LoadConfig(filename string) (IAMConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return IAMConfig{}, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	var config IAMConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return IAMConfig{}, fmt.Errorf("could not decode config file: %v", err)
	}

	return config, nil
}

// OverrideConfig allows the user to override any configuration settings.
func OverrideConfig(config *IAMConfig) {
	// Example: Override domain for the default solution
	var input string
	fmt.Printf("Enter the %s domain (leave blank to keep default): ", config.DefaultSolution)
	fmt.Scanln(&input)
	if input != "" {
		solutionConfig := config.IAMSolutions[config.DefaultSolution]
		solutionConfig.Domain = input                                // Modify the local copy
		config.IAMSolutions[config.DefaultSolution] = solutionConfig // Reassign it back to the map
	}
}

// GetIAMConfig retrieves the IAM configuration after loading it from the specified file.
func GetIAMConfig(filename string) (IAMConfig, error) {
	config, err := LoadConfig(filename)
	if err != nil {
		return IAMConfig{}, err
	}
	return config, nil
}
