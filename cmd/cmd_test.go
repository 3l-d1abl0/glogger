package cmd

import (
	"glogger/pkg/cli"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgs(t *testing.T) {

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
	//assert.Equal(t, cliStatus, true)

	assert.True(t, err == nil && filePath == fileArg && outputFolder == folderArg, "All valid")
}
