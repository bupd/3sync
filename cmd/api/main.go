package main

import (
	"3sync/internal/auth"
	"3sync/internal/gdrive"
	"3sync/internal/local"
	"3sync/utils"
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	gIDList   []string
	gNamelist []string
)

func main() {
	utils.LoadEnv()

	// check if token file exists
	if !utils.Exists("token.json") {
		doAuth()
	}

	// declare the map
	var localMap map[string]string
	// Initialize the map
	localMap = make(map[string]string)
	// list local files in folder
	localList := local.List()
	if localList == nil {
		// List File
		gIDList, gNamelist = gdrive.List()
		// download all items in the list
		for _, item := range gIDList {
			gdrive.Download(item)
		}

		return
	}
	// add items to local map
	for i := range localList {
		localMap[localList[i]] = ""
	}

	// declare the map
	var gListMap map[string]string
	// Initialize the map
	gListMap = make(map[string]string)

	// List File
	gIDList, gNamelist = gdrive.List()

	// check if all local files exist on gdrive
	for _, item := range localMap {
		value, exists := gListMap[item]
		if exists == true {
			delete(gListMap, item)
			delete(localMap, item)
		} else {
      gdrive.Download(gListMap[item])
    }
		// In case when key is not present in map variable exists will be false.
		fmt.Printf("key exists in map: %t, value: %v \n", exists, value)
	}
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
