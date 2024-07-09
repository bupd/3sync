package gdrive

import (
	"3sync/internal/auth"
	"context"
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
func Delete(fileID string) {
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
	// https://developers.google.com/drive/api/guides/search-files -- link to doc
	files, err := driveService.Files.Trash(fileID).Do()
	if err != nil {
		if googleapi.IsNotModified(err) {
			log.Fatalf("Unable to upload file: %v", err)
		}
	}

	fmt.Printf("File Deleted successfully: %v\n", files)
}
