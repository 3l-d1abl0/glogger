package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	File         string
	OutputFolder string
)

var rootCmd = &cobra.Command{
	Use:   "glogger",
	Short: "A simple CLI application for multiple concurrent download",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//File, outputFolder, err := parseArguments(cmd)

		// Store the parsed arguments in the context
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			//Check if flag "File" is setup
			fmt.Println(err)
			os.Exit(1)
		}
		File = filePath
		fmt.Println(File)

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			//Check if flag "output" is setup
			fmt.Println(err)
			os.Exit(1)
		}
		OutputFolder = output
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if File == "" {
			return fmt.Errorf("Error: %s", "Please provide a File Path")
		}

		if OutputFolder == "" {
			return fmt.Errorf("Error: %s", "Please provide an output folder")
		}

		return nil
	},
}

func init() {

	//Setting up Flags
	rootCmd.Flags().StringP("file", "f", "", "File path with url links")
	rootCmd.Flags().StringP("output", "o", "", "output folder path")

	//Marking mandatory Flag
	rootCmd.MarkFlagRequired("file")
	rootCmd.MarkFlagRequired("output")
}

func ParseCli() bool {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println("CLI: ", err)
		return false
	}
	return true
}
