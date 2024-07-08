package gdrive

import (
	"3sync/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Upload a file to the gdrive
func List() {
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

	// List the file
	// TO-DO: if problems occur change the insert to io.reader as expected.
	fileList, err := driveService.Files.List().Q("trashed=false").Do()
	if err != nil {
		if googleapi.IsNotModified(err) {
			log.Fatalf("Unable to upload file: %v", err)
		}
	}

	// Convert the file list to a JSON string for pretty printing
	listJson, err := json.MarshalIndent(fileList, "", "  ")
	if err != nil {
		log.Fatal("Unable to convert file list to JSON", err)
	}

	// Print the pretty-printed JSON
	fmt.Println("\nFile list JSON:\n", string(listJson))

	// Print the individual file names
	fmt.Println("\nFile names:")
	for _, file := range fileList.Items {
		fmt.Println(file.Title)
	}

	fmt.Printf("File list fetched successfully: %v\n", fileList)
}
