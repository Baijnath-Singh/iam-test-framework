package api

// RegisterResponse represents the response from the registration API.
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

// LoginResponse represents the response from the login API.
type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	Error        string `json:"error,omitempty"`
}

// AuthorizeResponse represents the response from the authorization API.
type AuthorizeResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

// TokenResponse represents the response from the token API.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}

// UserInfoResponse represents the response from the user info API.
type UserInfoResponse struct {
	Sub     string `json:"sub"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Picture string `json:"picture,omitempty"`
	Error   string `json:"error,omitempty"`
}

// LogoutResponse represents the response from the logout API.
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
