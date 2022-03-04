package workwithfiles

import (
	"fmt"
	"os"
)

var PathBase = ""

// func MakeFolder(path string) error {

// }

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
func SaveFile(name string, body []byte) bool {
	filePath := PathBase + "/" + name
	if exists, _ := FileExists(filePath); exists {
		fmt.Println("File already exists")
		return false
	} else {
		err := os.WriteFile(filePath, body, 0777)
		return err == nil
	}
}
