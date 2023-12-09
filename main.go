package main

import (
	"fmt"
	"glogger/cmd"
	"glogger/pkg/commondata"
	"glogger/pkg/glogger"
	"glogger/pkg/progressbar"
	"glogger/pkg/readUrls"

	"github.com/fatih/color"
)

func userInputWait() (bool, error) {

	var userInput string

	cCy := color.New(color.FgCyan)
	cMeg := color.New(color.FgMagenta).Add(color.Bold)
	for {
		cCy.Printf("Do you want to proced ? (Y/n): ")

		_, err := fmt.Scan(&userInput)
		if err != nil {
			fmt.Printf("\n Error reading input: %s \n", err.Error())
			return false, err
		}

		//Check user Input
		if userInput == "Y" {
			cCy.Printf("\nYou entered : %s\n", userInput)
			return true, nil
		} else if userInput == "n" {
			cCy.Printf("\nYou entered : %s\n", userInput)
			return false, nil
		} else {
			cMeg.Printf("\nInvalid input. Please enter Y or n.\n")
		}
	}
}

func main() {

	//1.Try to fetch the Arguments
	urlFilePath, outputFolder, err := cmd.GetArgs()

	if err != nil {
		fmt.Printf("MAIN: Unable to Read input: (%s) \n", err)
		return
	}

	cCy := color.New(color.FgCyan) //Normal
	cGB := color.New(color.FgGreen).Add(color.BgBlack).Add(color.Bold).SprintFunc()
	cMag := color.New(color.FgMagenta).Add(color.Bold)
	fmt.Printf("Path to the Urls file: %s\n", cGB(urlFilePath))
	fmt.Printf("Path to the output folder: %s\n", cGB(outputFolder))

	//2.Parse the urls from input file
	var targetUrls commondata.TargetUrls
	totalSize, err := readUrls.ReadUrls(&urlFilePath, &targetUrls)
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
	ans, err := userInputWait()
	if err != nil {
		cMag.Printf("ERROR: %s\n", err.Error())
		return
	} else {

		if ans == false {
			cMag.Printf("Exiting ...\n")
			return
		} else {
			cCy.Printf("Continuing ...\n")
		}
	}

	//4. Setup channels
	//Channel to transfer send messages
	doneCh := make(chan commondata.Message)
	//Channel to transfer any other messages
	msgCh := make(chan string)
	//Bidirection Channel, used as send Quit message
	quitCh := make(chan bool)

	//5. Set up Progressbar Downloader Settings
	var N int = len(targetUrls.ValidUrls)
	var barSize int64 = 50
	var barSymbol string = ">"
	bar := progressbar.GetNewBar(int64(N), 0, barSymbol, barSize, 0, totalSize)

	//6. Start Download Process
	glogger.Glogger(doneCh, msgCh, quitCh, targetUrls.ValidUrls, outputFolder, &bar)

	//Wait for the Quit Signal
	if true == <-quitCh {
		boldMag := color.New(color.FgMagenta).Add(color.Bold)
		boldMag.Printf("\nExiting ! ")
	}

}
