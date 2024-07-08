package auth

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bupd/goth"
	"github.com/bupd/goth/gothic"
	"github.com/bupd/goth/providers/google"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v2"
)

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

func GetClient() *http.Client {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{drive.DriveFileScope}, // Full access to Google Drive
	}

	// Create a token with the refresh token
	token := &oauth2.Token{RefreshToken: refreshToken}

	// Use the token source to refresh the token
	tokenSource := config.TokenSource(context.Background(), token)

	// Get the new access token
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Unable to get access token: %v", err)
	}

	// Return an HTTP client with the refreshed token
	return config.Client(context.Background(), newToken)
}

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error in loading env file: %v", err)
	}

	// Retrieve environment variables
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")

	if googleClientID == "" || googleClientSecret == "" || redirectURI == "" {
		log.Fatalf("Google OAuth environment variables not set")
	}

	// Initialize session store
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // Prevent client-side scripts from accessing the cookie
	store.Options.Secure = IsProd // Secure cookies are only sent over HTTPS

	// Assign the session store to gothic
	gothic.Store = store

	_ = []string{
		"https://www.googleapis.com/auth/drive",
	}

	// Initialize the Google provider
	goth.UseProviders(
		google.New(
			googleClientID,
			googleClientSecret,
			redirectURI,
		),
	)
}
