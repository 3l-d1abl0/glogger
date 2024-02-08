package utils

import (
	"os"
	"path/filepath"
)

func IsValidFilePath(path string) bool {

	// Clean and normalize the path
	path = filepath.Clean(path)

	// Check if the path is absolute
	return filepath.IsAbs(path)
}

// Check if the input file exist
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Check if the input folder exists
func FolderExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {

		// Directory does not exist
		if os.IsNotExist(err) {
			return false
		}
		// Other error,
		//fmt.Printf("Error checking directory: %v\n", err)
		return false
	}
	return fileInfo.IsDir()
}
