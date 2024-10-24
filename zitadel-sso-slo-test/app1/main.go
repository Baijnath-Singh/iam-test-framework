package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var oauth2Config = &oauth2.Config{
	ClientID:     "290577641713500163",                                               // Replace with your App1 Client ID
	ClientSecret: "PYStJXKbxkG47WLigwJebFnHrasUHUl3Qv1rNIiNTnLLwbxwCkPUNbrgfdj474Dx", // Replace with your App1 Client Secret
	RedirectURL:  "http://localhost:3001/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://localhost:8080/oauth/v2/authorize",
		TokenURL: "http://localhost:8080/oauth/v2/token",
	},
	Scopes: []string{"openid", "profile", "email"},
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("App1 is running on http://localhost:3001")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		logFatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<a href="/login">Login with ZITADEL (App1)</a>`)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Unable to exchange code for token", http.StatusInternalServerError)
		logFatal(err)
		return
	}

	client := oauth2Config.Client(ctx, token)
	resp, err := client.Get("http://localhost:8080/oidc/v1/userinfo")
	if err != nil {
		http.Error(w, "Unable to get user info", http.StatusInternalServerError)
		logFatal(err)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Unable to decode user info", http.StatusInternalServerError)
		logFatal(err)
		return
	}

	// Set the content type to text/html
	w.Header().Set("Content-Type", "text/html")

	// Display user info and a logout link
	fmt.Fprintf(w, "Logged in as: %v<br>", userInfo)
	fmt.Fprintln(w, `<a href="/logout">Logout</a>`)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/oidc/v1/end_session?post_logout_redirect_uri=http://localhost:3001", http.StatusTemporaryRedirect)
}

func logFatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
