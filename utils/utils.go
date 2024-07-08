package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
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
}
