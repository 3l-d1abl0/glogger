package cmd

import (
	"fmt"
	"glogger/pkg/cli"
	"os"
	"path/filepath"
)

func isValidFilePath(path string) bool {

	// Clean and normalize the path
	path = filepath.Clean(path)

	// Check if the path is absolute
	return filepath.IsAbs(path)
}

// Check if the input file exist
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Check if the input folder exists
func folderExists(path string) bool {
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

func GetArgs() (string, string, error) {

	if err := cli.ParseCli(); err != true {

		//Unbale to read args
		return "", "", fmt.Errorf("Error: %s", "Please provide a file Path Empty")
	}

	//fmt.Println(cli.File, cli.OutputFolder)
	//Check if the input is valid file paths
	if isValidFilePath(cli.File) == false {
		fmt.Printf("%s\n : is not a valid file Path", cli.File)
		return "", "", fmt.Errorf("Error: %s", "Not a valid File Path")
	}

	if isValidFilePath(cli.OutputFolder) == false {
		fmt.Printf("%s\n : is not a valid file Path", cli.OutputFolder)
		return "", "", fmt.Errorf("Error: %s", "Not a valid Folder Path")
	}

	//Check if the url file exists
	if fileExists(cli.File) == false {
		fmt.Printf("file does not exist : %s \n", cli.File)
		return "", "", fmt.Errorf("Error: %s", "url file does not exist")
	}

	if folderExists(cli.OutputFolder) == false {
		fmt.Printf("folder does not exist : %s \n", cli.OutputFolder)
		return "", "", fmt.Errorf("Error: %s", "folder does not exist")
	}

	return cli.File, cli.OutputFolder, nil
}
