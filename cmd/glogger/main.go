package main

import (
	"fmt"
	"glogger/internal"
	"glogger/pkg/cli"
	"glogger/pkg/commondata"
	"glogger/pkg/utils"
	"os"

	"github.com/fatih/color"
)

func getArgs() (string, string, error) {

	if err := cli.ParseCli(); err != true {

		//Unbale to read args
		return "", "", fmt.Errorf("Error: %s", "Please provide a file Path Empty")
	}

	//fmt.Println(cli.File, cli.OutputFolder)
	//Check if the input is valid file paths
	if utils.IsValidFilePath(cli.File) == false {
		fmt.Printf("%s\n : is not a valid file Path", cli.File)
		return "", "", fmt.Errorf("Error: %s", "Not a valid File Path")
	}

	if utils.IsValidFilePath(cli.OutputFolder) == false {
		fmt.Printf("%s\n : is not a valid file Path", cli.OutputFolder)
		return "", "", fmt.Errorf("Error: %s", "Not a valid Folder Path")
	}

	//Check if the url file exists
	if utils.FileExists(cli.File) == false {
		fmt.Printf("file does not exist : %s \n", cli.File)
		return "", "", fmt.Errorf("Error: %s", "url file does not exist")
	}

	if utils.FolderExists(cli.OutputFolder) == false {
		fmt.Printf("folder does not exist : %s \n", cli.OutputFolder)
		return "", "", fmt.Errorf("Error: %s", "folder does not exist")
	}

	return cli.File, cli.OutputFolder, nil
}

func main() {

	//1.Try to fetch the Arguments
	urlFilePath, outputFolder, err := getArgs()

	if err != nil {
		fmt.Printf("MAIN: Unable to Read input: (%s) \n", err)
		return
	}

	cCy := color.New(color.FgCyan) //Normal
	cGB := color.New(color.FgGreen).Add(color.BgBlack).Add(color.Bold).SprintFunc()
	cMag := color.New(color.FgMagenta).Add(color.Bold)

	//Print to Terminals
	fmt.Printf("Path to the Urls file: %s\n", cGB(urlFilePath))
	fmt.Printf("Path to the output folder: %s\n", cGB(outputFolder))

	//2.Parse the urls from input file
	var targetUrls commondata.TargetUrls
	totalSize, err := utils.ReadUrls(&urlFilePath, &targetUrls)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Total : %s url(s)  (%s valid + %s invalid url(s) )\n",
		cGB(len(targetUrls.InvalidUrls)+len(targetUrls.ValidUrls)),
		cCy.SprintFunc()(len(targetUrls.ValidUrls)),
		cCy.SprintFunc()(len(targetUrls.InvalidUrls)),
	)

	cCy.Printf("Total size to fetch : [size %.2f MB]\n", float32(totalSize)/(1024*1024))

	//3. Wait for user Input
	ans, err := utils.UserInputWait()
	if err != nil {
		cMag.Printf("ERROR: %s\n", err.Error())
		return
	} else {

		if ans == false {
			cMag.Printf("Exiting ...\n")
			os.Exit(0)
		} else {
			cCy.Printf("continuing ...\n")
			internal.Run(&totalSize, &targetUrls, &outputFolder)
		}
	}
}
