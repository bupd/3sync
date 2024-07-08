package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	drive "google.golang.org/api/drive/v3"
)

// Create OAuth2 config without using credentials.json
func getOAuthConfig() *oauth2.Config {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := "urn:ietf:wg:oauth:2.0:oob" // Out-of-band, for command-line apps

	return &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       []string{drive.DriveScope}, // Modify scopes as required
	}
}

// Retrieve a token, save it, and return the refresh token.
func getRefreshToken(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	// Save the token to a file
	saveToken("token.json", tok)

	// Return the refresh token
	return tok.RefreshToken
}

// Save the token to a file
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving token file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Token retrieval and saving
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error in loading env file: %v", err)
	}

	// Retrieve environment variables

	config := getOAuthConfig()

	// Get and print the refresh token
	refreshToken := getRefreshToken(config)
	fmt.Printf("Your refresh token: %s\n", refreshToken)

	// If you need to use the token to make a request, you can do so like this:
	// client := config.Client(context.Background(), &oauth2.Token{RefreshToken: refreshToken})
	// Use `client` for making requests to the Google API.
}
