package gdrive

import (
	"3sync/internal/auth"
	"3sync/utils"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

// Download a file from gdrive
func Download(fileID string) {
	// TO-DO: remove the config and get from the main
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
	}

	// TO-DO: spawn this in the separate function
	client := auth.GetClient(config)
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Drive client: %v", err)
	}

	// Download the file
	// TO-DO: if problems occur change the insert to io.reader as expected.
	// https://developers.google.com/drive/api/guides/search-files -- link to doc

	// Get the file metadata
	file, err := driveService.Files.Get(fileID).Do()
	if err != nil {
		log.Fatalf("Unable to get file: %v", err)
	}

	// Download the file content
	response, err := driveService.Files.Get(fileID).Download()
	if err != nil {
		log.Fatalf("Unable to download file: %v", err)
	}
	defer response.Body.Close()

	// Create the directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get home directory: %v", err)
	}

	downloadDir := filepath.Join(homeDir, "gdrive")
	fmt.Println("downloadDir: ", downloadDir)

	if !utils.Exists(downloadDir) {
		err = os.MkdirAll(downloadDir, 0755)
		if err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
	}

	// Create the local file
	localFilePath := filepath.Join(downloadDir, file.Title)
	localFile, err := os.Create(localFilePath)
	if err != nil {
		log.Fatalf("Unable to create local file: %v", err)
	}
	defer localFile.Close()

	// Copy the content to the local file
	_, err = io.Copy(localFile, response.Body)
	if err != nil {
		log.Fatalf("Unable to save file: %v", err)
	}

	fmt.Printf("File downloaded successfully: %s\n", localFilePath)
}
