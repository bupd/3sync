package local

import (
	"3sync/utils"
	"os"
	"path/filepath"
)

// create ~/gdrive folder
func CreateFolder() {
	home := utils.GetHomeDir()
	dir := filepath.Join(home, "gdrive")
	os.Mkdir(dir, 0700)
}
