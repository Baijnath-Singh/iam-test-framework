package routes

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"iam-test-framework/config"
	"iam-test-framework/input"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time" // Import time package for timeout
)

// GenerateCodeVerifier creates a secure random string for the code verifier (43-128 characters long).
func GenerateCodeVerifier() (string, error) {
	verifierBytes := make([]byte, 32)
	_, err := rand.Read(verifierBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(verifierBytes), nil
}

// GenerateCodeChallenge creates the SHA-256 hash of the code verifier and returns the base64-encoded code challenge.
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// OpenBrowser opens the given URL in the default browser, cross-platform.
func OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		if isWSL() {
			// Open the URL using the Windows command prompt
			err = exec.Command("cmd.exe", "/C", "start", url).Start()
		} else {
			err = exec.Command("xdg-open", url).Start()
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

// isWSL checks if the code is running in WSL (Windows Subsystem for Linux).
func isWSL() bool {
	// Check for the presence of "WSL" in the kernel version
	return strings.Contains(strings.ToLower(string(os.Getenv("WSL_DISTRO_NAME"))), "wsl") ||
		strings.Contains(strings.ToLower(string(os.Getenv("TERM"))), "wsl")
}

// OIDCAuthorization handles the OIDC authorization flow with PKCE support and a local callback server.
func OIDCAuthorization() {

	fmt.Println("NOTE: This feature implements OIDC Authorization with PKCE")
	fmt.Printf("\n\n")

	// Prompt user to select IAM solution
	var solutionName string
	fmt.Println("Enter the IAM solution (zitadel, keycloak, casdoor):")
	fmt.Scanln(&solutionName)

	iamConfig, err := config.GetIAMConfig("config/config.json") // Fetch the relevant config for the IAM solution
	if err != nil {
		fmt.Printf("Error loading IAM config: %v\n", err)
		return
	}

	solutionConfig, exists := iamConfig.IAMSolutions[solutionName]
	if !exists {
		fmt.Printf("IAM solution %s does not exist.\n", solutionName)
		return
	}

	// Get client ID, redirect URI, and scope from environment variables or prompt user
	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		fmt.Print("Enter the Client ID: ")
		clientID = strings.TrimSpace(input.GetUserInput()) // Trim any extra spaces
	}

	redirectURI := "http://localhost:8081/callback" // Local callback server

	scope := os.Getenv("SCOPE")
	if scope == "" {
		fmt.Print("Enter the Scope (e.g., openid profile email): ")
		scope = strings.TrimSpace(input.GetUserInput()) // Trim any extra spaces
	}

	if clientID == "" || redirectURI == "" || scope == "" {
		fmt.Println("Error: CLIENT_ID, REDIRECT_URI, and SCOPE must be set.")
		return
	}

	// Add the required custom scope for Zitadel
	scope += " urn:zitadel:iam:user:metadata"

	// Generate PKCE code verifier and code challenge
	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		fmt.Printf("Error generating code verifier: %v\n", err)
		return
	}
	codeChallenge := GenerateCodeChallenge(codeVerifier)

	// URL-encode the redirect URI and scope (use %20 for spaces instead of +)
	encodedRedirectURI := url.QueryEscape(redirectURI)
	encodedScope := strings.ReplaceAll(url.QueryEscape(scope), "+", "%20")

	// Ensure the domain uses HTTP instead of HTTPS
	domain := solutionConfig.Domain
	if strings.HasPrefix(domain, "https://") {
		domain = strings.Replace(domain, "https://", "http://", 1)
	}

	// Combine domain with the token endpoint
	tokenEndpoint := fmt.Sprintf("%s%s", domain, solutionConfig.TokenEndpoint) // Full token endpoint

	// Prepare the authorization URL with encoded parameters and PKCE
	authURL := fmt.Sprintf("%s%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&code_challenge=%s&code_challenge_method=S256",
		domain, solutionConfig.AuthorizeEndpoint, clientID, encodedRedirectURI, encodedScope, codeChallenge)

	// Store the code verifier for use in the token exchange step
	fmt.Printf("Save this code verifier for the token exchange step: %s\n", codeVerifier)

	// Start the local callback server in a goroutine
	var wg sync.WaitGroup // Create a WaitGroup
	wg.Add(1)             // Add a count to the WaitGroup

	go func() {
		defer wg.Done() // Signal that this goroutine is done when it exits
		startCallbackServer(tokenEndpoint, clientID, redirectURI, codeVerifier)
	}()

	// Automatically open the authorization URL in the user's browser
	err = OpenBrowser(authURL)
	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	} else {
		fmt.Println("Authorization URL opened in your default browser.")
	}

	// Wait for the callback server to complete
	wg.Wait() // Wait for the callback server goroutine to finish
}

// startCallbackServer starts the local callback server on a specified port.
func startCallbackServer(tokenEndpoint, clientID, redirectURI, codeVerifier string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Log the entire request for debugging
		fmt.Printf("Received request: %s\n", r.URL.String())

		// Extract authorization code from query parameters
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Fprintln(w, "Authorization code not found.")
			fmt.Println("Error: Authorization code not found in the request.")
			return
		}

		// Exchange authorization code for tokens
		tokens, err := exchangeCodeForTokens(tokenEndpoint, clientID, redirectURI, code, codeVerifier)
		if err != nil {
			fmt.Fprintf(w, "Error exchanging code for tokens: %v\n", err)
			fmt.Println("Error exchanging code for tokens:", err)
			return
		}

		// Display tokens in browser
		response := fmt.Sprintf("Access Token: %s\nID Token: %s\nRefresh Token: %s\n", tokens["access_token"], tokens["id_token"], tokens["refresh_token"])
		fmt.Fprintln(w, response)

		// Display tokens on terminal
		fmt.Println("Received Tokens:")
		fmt.Println("Access Token:", tokens["access_token"])
		fmt.Println("ID Token:", tokens["id_token"])
		fmt.Println("Refresh Token:", tokens["refresh_token"])

		// Terminate the server after handling the callback
		go func() {
			time.Sleep(1 * time.Second) // Delay to ensure the response is sent before shutting down
			os.Exit(0)                  // Exit the program cleanly
		}()
	})

	// Start the server on localhost:8081
	fmt.Println("Starting callback server on http://localhost:8081...") // Log before starting
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Println("Error starting the callback server:", err)
		return // Exit the goroutine if it fails to start
	}
	fmt.Println("Callback server is running on http://localhost:8081") // Log after starting
}

// exchangeCodeForTokens exchanges the authorization code for tokens by making a POST request.
func exchangeCodeForTokens(tokenEndpoint, clientID, redirectURI, code, codeVerifier string) (map[string]interface{}, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)
	data.Set("code_verifier", codeVerifier)

	resp, err := http.PostForm(tokenEndpoint, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokens map[string]interface{}
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
