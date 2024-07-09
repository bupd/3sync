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
func GetID(file string) string {
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
	fileList, err := driveService.Files.List().Q(file).Q("trashed=false").Do()
	if err != nil {
		if googleapi.IsNotModified(err) {
			log.Fatalf("Unable to fetch & get file: %v", err)
		}
	}

	var IDList []string
	var NameList []string
	// Print the individual file names
	fmt.Println("\nFile names:")
	for _, file := range fileList.Items {
		fmt.Println(file.Title, file.Id)
		IDList = append(IDList, file.Id)
		NameList = append(NameList, file.Title)
	}
	if len(IDList) > 0 {
		return IDList[0]
	}

	return ""
}
