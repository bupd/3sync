package local

import (
	"fmt"
	"os"
)

func Delete(fileName string) {
	err := os.Remove(fileName) // remove the file
	if err != nil {
		fmt.Println("Error: ", err) // print the error if file is not removed
	} else {
		fmt.Println("Successfully deleted file: ", fileName) // print success if file is removed
	}
}
