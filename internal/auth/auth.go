package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bupd/goth"
	"github.com/bupd/goth/gothic"
	"github.com/bupd/goth/providers/google"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

// Function to read token from file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Function to save token to file
func saveToken(file string, token *oauth2.Token) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// Function to get authenticated HTTP client
func GetClient(config *oauth2.Config) *http.Client {
	// Read token from file
	tokenFile := "token.json"
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		log.Fatalf("Unable to read token file: %v", err)
	}

	// Use token source to ensure it's valid
	tokenSource := config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Unable to refresh access token: %v", err)
	}

	// Save new token if it has changed
	if newToken.AccessToken != token.AccessToken {
		if err := saveToken(tokenFile, newToken); err != nil {
			log.Fatalf("Unable to save updated token: %v", err)
		}
	}

	// Return an HTTP client with the token source
	return oauth2.NewClient(context.Background(), tokenSource)
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
