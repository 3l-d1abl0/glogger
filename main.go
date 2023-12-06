package main

import (
	"fmt"
	"glogger/cmd"
	"glogger/pkg/commondata"
	"glogger/pkg/glogger"
	"glogger/pkg/readUrls"

	"github.com/3l-d1abl0/goProgressBar"
)

func main() {

	//1.Try to fetch the Arguments
	urlFilePath, outputFolder, err := cmd.GetArgs()

	if err != nil {
		fmt.Printf("MAIN: Unable to Read input: (%s) \n", err)
		return
	}

	//2.Parse the urls from input file
	var targetUrls commondata.TargetUrls
	readUrls.ReadUrls(&urlFilePath, &targetUrls)

	//3. Setup channels
	//Channel to transfer send messages
	doneCh := make(chan string)
	//Channel to transfer other messages
	msgCh := make(chan string)
	//Bidirection Channel, used as send Quit message
	quitCh := make(chan bool)

	//4. Set up Progressbar Downloader Settings
	var N int = len(targetUrls.ValidUrls)
	var barSize int64 = 70
	var barSymbol string = "#"
	bar := goProgressBar.GetNewBar(int64(N), 0, barSymbol, barSize)

	//boldGreen := color.New(color.FgGreen, color.Bold)
	//boldGreen.Printf(" %d/%d Started ...\n", countReq, nOReq)

	//5. Start Download Process
	glogger.Glogger(doneCh, msgCh, targetUrls.ValidUrls, outputFolder)

	//6. Setting up the recievers
	go glogger.Receiver(doneCh, msgCh, quitCh, &bar, N)

	//Wait for the Quit Signal
	println(<-quitCh)

}
