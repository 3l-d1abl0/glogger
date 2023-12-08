package cmd

import (
	"glogger/pkg/cli"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Ideally a setUp and tearDown
function should create a file and output folder

*/

func TestValidArgs(t *testing.T) {

	filePath := "/media/murphy/P3/groundzero/glogger/urls/local-testing.txt"
	outputFolder := "/media/murphy/P3/groundzero/glogger/output"

	//Create the cli arguments
	args := []string{
		"-f", filePath,
		"-o", outputFolder,
	}

	//Set the arguments for cli Context
	rootCmd := cli.GetCli()
	rootCmd.SetArgs(args)

	//parse the arguments
	fileArg, folderArg, err := GetArgs()

	//Should be true
	assert.True(t, err == nil && filePath == fileArg && outputFolder == folderArg, "All valid")
}

func TestInvalidArg(t *testing.T) {

	var testArgs = []struct {
		flag1      string
		filePath   string
		flag2      string
		folderPath string
		expected   *string
		comments   string
	}{
		{"--file", "/media/murphy/P3/groundzero/glogger/urls/local-testing1.txt", "-o", "/media/murphy/P3/groundzero/glogger/output", nil, "File path Invalid"},
		{"--file", "/media/murphy/P3/groundzero/glogger/urls/local-testing.txt", "--output", "/media/murphy/P3/groundzero/glogger/input", nil, "Folder Path Invalid"},
		{"-f", "/media/murphy/P3/groundzero/glogger/urls/local-testing2.txt", "-o", "/media/murphy/P3/groundzero/glogger/input", nil, "Both input file and output folder does not exist"},
	}

	//Get the cli Context
	rootCmd := cli.GetCli()

	for _, test := range testArgs {

		//Set the arguments for cli Context
		rootCmd.SetArgs([]string{test.flag1, test.filePath, test.flag2, test.folderPath})

		//parse the arguments
		_, _, err := GetArgs()

		//Should be true
		assert.True(t, err != nil, err)
	}

}
