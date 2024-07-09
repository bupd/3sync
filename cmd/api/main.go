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

	// send all keys in glist map to download
	// // Write to file
	// err := writeLines(gNamelist, "prev.txt")
	// if err != nil {
	// 	log.Fatalf("Error in writing drive list to the file: %v", err)
	// }

	// Read from file
	gOldList, err := readLines("prev.txt")
	if err != nil {
		log.Fatalf("Error in reading drive list from the file: %v", err)
	}

	// Write to file
	err = writeLines(gNamelist, "prev.txt")
	if err != nil {
		log.Fatalf("Error in writing drive list to the file: %v", err)
	}

	for i := range gOldList {
		_, exists := localMap[gOldList[i]]
		if exists {
			delete(localMap, gOldList[i])
		} else {
			id := gdrive.GetID(gOldList[i])
			gdrive.Delete(id)
		}
	}
	// UploadFile leftover local map
	for _, item := range localMap {
		gdrive.UploadFile(item)
	}

	for i := range localList {
		_, exists := gListMap[localList[i]]
		if exists {
			delete(gListMap, localList[i])
		} else {
			local.Delete(localList[i])
		}
	}

	// UploadFile leftover local map
	for _, item := range gListMap {
		gdrive.Download(gListMap[item])
	}

	// send all glist to previous.txt

	// delete all items in the list
	//
	// for _, item := range glist {
	// 	gdrive.Delete(item)
	// }
}

// gets the refreshToken and creates the oauth token
func doAuth() {
	config := auth.GetOAuthConfig()
	refreshToken := auth.GetRefreshToken(config)
	fmt.Printf("Your refresh token: %s\n", refreshToken)
}

// readLines reads a file into a slice of strings
func readLines(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

// writeLines writes a slice of strings to a file
func writeLines(lines []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}
