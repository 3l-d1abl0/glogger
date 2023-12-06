package cmd

import (
	"fmt"
	"glogger/pkg/cli"
)

func GetArgs() (string, string, error) {

	if err := cli.ParseCli(); err != true {

		//Unbale to read args
		return "", "", fmt.Errorf("Error: %s", "Please provide a file Path Empty")
	}

	return cli.File, cli.OutputFolder, nil
}
