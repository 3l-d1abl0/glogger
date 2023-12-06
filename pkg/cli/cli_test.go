package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {

	filePath := "/media/murphy/P3/groundzero/glogger/urls/local-testing.txt"
	outputFolder := "/media/murphy/P3/groundzero/glogger/output"

	//Create the cli arguments
	args := []string{
		"-f", filePath,
		"-o", outputFolder,
	}

	//Set the arguments for cli Context
	rootCmd.SetArgs(args)

	//parse the cli arguments
	cliStatus := ParseCli()

	//Should be true
	assert.Equal(t, cliStatus, true)
}
