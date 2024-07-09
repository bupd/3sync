package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load the .env file to the program environment
func LoadEnv() {
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

func GetHomeDir() string {
	// get home Directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	fmt.Println("Home Directory:", home)
	return home
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("%s File exists\n", path)
		return true
	} else {
		fmt.Printf("%s File does not exist\n", path)
		return false
	}
}
