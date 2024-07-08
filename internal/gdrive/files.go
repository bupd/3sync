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
func UploadFile() {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
	}

	client := auth.GetClient(config)
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Drive client: %v", err)
	}

	// Open the file to be uploaded
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current Directory: %v", err)
	}

	filePath := currDir + "/hires.jpg"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	fileName := file.Name()
	log.Println(fileName)

	// File metadata
	fileMetadata := &drive.File{
		Title: "hires.jpg",
	}

	// Upload the file
	// TO-DO: if problems occur change the insert to io.reader as expected.
	fileUpload, err := driveService.Files.Insert(fileMetadata).Media(file).Do()
	if err != nil {
		if googleapi.IsNotModified(err) {
			log.Fatalf("Unable to upload file: %v", err)
		}
	}

	fmt.Printf("File uploaded successfully: %v\n", fileUpload)
}
