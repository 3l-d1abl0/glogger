package main

import (
	"fmt"
	"glogger/cmd"
	"glogger/pkg/commondata"
	"glogger/pkg/glogger"
	"glogger/pkg/readUrls"

	"github.com/3l-d1abl0/goProgressBar"
	"github.com/fatih/color"
)

func main() {

	//1.Try to fetch the Arguments
	urlFilePath, outputFolder, err := cmd.GetArgs()

	if err != nil {
		fmt.Printf("MAIN: Unable to Read input: (%s) \n", err)
		return
	}

	cCy := color.New(color.FgCyan) //Normal
	cGB := color.New(color.FgGreen).Add(color.BgBlack).Add(color.Bold).SprintFunc()
	fmt.Printf("Path to the Urls file: %s\n", cGB(urlFilePath))
	fmt.Printf("Path to the output folder: %s\n", cGB(outputFolder))

	//2.Parse the urls from input file
	var targetUrls commondata.TargetUrls
	readUrls.ReadUrls(&urlFilePath, &targetUrls)

	fmt.Printf("Total : %s url(s)  (%s valid + %s invalid url(s) )\n",
		cGB(len(targetUrls.InvalidUrls)+len(targetUrls.ValidUrls)),
		cCy.SprintFunc()(len(targetUrls.ValidUrls)),
		cCy.SprintFunc()(len(targetUrls.InvalidUrls)),
	)

	//3. Setup channels
	//Channel to transfer send messages
	doneCh := make(chan string)
	//Channel to transfer any other messages
	msgCh := make(chan string)
	//Bidirection Channel, used as send Quit message
	quitCh := make(chan bool)

	//4. Set up Progressbar Downloader Settings
	var N int = len(targetUrls.ValidUrls)
	var barSize int64 = 70
	var barSymbol string = "#"
	bar := goProgressBar.GetNewBar(int64(N), 0, barSymbol, barSize)

	//5. Start Download Process
	glogger.Glogger(doneCh, msgCh, quitCh, targetUrls.ValidUrls, outputFolder, &bar)

	//6. Setting up the recievers
	//go glogger.Receiver(doneCh, msgCh, quitCh, &bar, N)

	//Wait for the Quit Signal
	if true == <-quitCh {
		boldMag := color.New(color.FgMagenta).Add(color.Bold)
		boldMag.Printf("\nExiting ! ")
	}

}
