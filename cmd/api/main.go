package main

import (
	"3sync/internal/auth"
	"3sync/internal/server"
	"fmt"
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
