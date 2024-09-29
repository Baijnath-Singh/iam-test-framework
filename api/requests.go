package api

// UserRegistrationRequest defines the structure for user registration requests.
type UserRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserLoginRequest defines the structure for user login requests.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// OIDCAuthRequest defines the structure for OIDC authorization requests.
type OIDCAuthRequest struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
	Scope        string `json:"scope"`
}

// TokenRequest defines the structure for token requests.
type TokenRequest struct {
	GrantType string `json:"grant_type"`
	ClientID  string `json:"client_id"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

// UserInfoRequest defines the structure for user info requests.
type UserInfoRequest struct {
	AccessToken string `json:"access_token"`
}

// LogoutRequest defines the structure for logout requests.
type LogoutRequest struct {
	AccessToken string `json:"access_token"`
}

// SSORequest defines the structure for SSO requests.
type SSORequest struct {
	State string `json:"state"`
}
