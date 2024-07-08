package main

import (
	"3sync/internal/gdrive"
	"3sync/utils"
)

const (
	clientID     = "YOUR_CLIENT_ID"
	clientSecret = "YOUR_CLIENT_SECRET"
	refreshToken = "YOUR_REFRESH_TOKEN"
)

func main() {
	// Instantiate gothic
	// auth.NewAuth()

	utils.LoadEnv()
	gdrive.UploadFile()

	// server := server.NewServer()
	//
	// err := server.ListenAndServe()
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot start server: %s", err))
	// }
}
