package main

import (
	"3sync/internal/auth"
	"3sync/internal/server"
	"fmt"
)

const (
	clientID     = "YOUR_CLIENT_ID"
	clientSecret = "YOUR_CLIENT_SECRET"
	refreshToken = "YOUR_REFRESH_TOKEN"
)

func main() {
	// Instantiate gothic
	auth.NewAuth()

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
