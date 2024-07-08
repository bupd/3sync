package main

import (
	"3sync/internal/auth"
	"3sync/internal/gdrive"
	"3sync/utils"
	"fmt"
	"os"
)

func main() {
	utils.LoadEnv()

	// check if token file exists
	if !fileExists() {
		doAuth()
	}

  // List File
  gdrive.List()
	// UploadFile
	// gdrive.UploadFile()
}

// gets the refreshToken and creates the oauth token
func doAuth() {
	config := auth.GetOAuthConfig()
	refreshToken := auth.GetRefreshToken(config)
	fmt.Printf("Your refresh token: %s\n", refreshToken)
}

// Check if token file exists or not
func fileExists() bool {
	if _, err := os.Stat("token.json"); err == nil {
		fmt.Printf("File exists\n")
		return true
	} else {
		fmt.Printf("Token File does not exist\n")
		return false
	}
}
