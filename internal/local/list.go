package local

import (
	"3sync/utils"
	"fmt"
	"os"
	"path/filepath"
)

// List files in local ~/gdrive folder
func List() []string {
	home := utils.GetHomeDir()
	dir := filepath.Join(home, "gdrive")
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	var localFileNames []string

	for _, file := range files {
		fmt.Println(file.Name())
		localFileNames = append(localFileNames, file.Name())
	}

	return localFileNames
}
