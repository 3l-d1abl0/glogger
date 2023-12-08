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

func TestAllInputCases(t *testing.T) {

	var testArgs = []struct {
		flag1      string
		filePath   string
		flag2      string
		folderPath string
		expected   bool
		comments   string
	}{
		{"", "", "", "", false, "No Arguments"},
		{"-f", "", "-o", "/media/murphy/P3/groundzero/glogger/output", false, "No file path input"},
		{"--file", "/media/murphy/P3/groundzero/glogger/urls/local-testing.txt", "--output", "", false, "No outputfolder path"},
		{"-f", "/media/murphy/P3/groundzero/glogger/urls/local-testing2.txt", "-o", "/media/murphy/P3/groundzero/glogger/input", true, "Both input file and output folder does not exist, but is a valid path"},

		{"-f", "/media/murphy/P3/groundzero/glogger/urls/local-testing.txt", "-o", "/media/murphy/P3/groundzero/glogger/output", true, "Both valid"},
		{"--file", "/media/murphy/P3/groundzero/glogger/urls/local-testing.txt", "--output", "/media/murphy/P3/groundzero/glogger/output", true, "Both valid"},
	}

	for _, test := range testArgs {

		//Set the arguments for cli Context
		rootCmd.SetArgs([]string{test.flag1, test.filePath, test.flag2, test.folderPath})
		assert.Equal(t, test.expected, ParseCli())
	}

}
