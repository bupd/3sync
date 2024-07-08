package auth

import (
	"log"
	"os"

	"github.com/bupd/goth"
	"github.com/bupd/goth/gothic"
	"github.com/bupd/goth/providers/google"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

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
