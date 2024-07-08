package gdrive

import (
	"3sync/internal/auth"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/drive/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func uploadFile() {
	client := auth.GetClient()
	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Drive client: %v", err)
	}

	filePath := "../../hires.jpg"
	fileExtension := ".jpg"
	// Open the file to be uploaded
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current Directory: %v", err)
	}

	file, err := os.Open(currDir + filePath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer file.Close()

	fileName := file.Name()

	// File metadata
	fileMetadata := &drive.File{
		Title: fileName + fileExtension,
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
