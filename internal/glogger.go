package internal

import (
	"fmt"
	"glogger/pkg/commondata"
	"glogger/pkg/progressbar"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type terminalMessage struct {
	id      int
	message string
}

var totalReqs int
var reqCount int
var mu sync.Mutex
var currentSize int64
var previousMessageType int

func flushToTerminal(msgs *[]terminalMessage) {

	//Color Convention
	//1, boldGreen//for done channel
	//2, boldRed //for msg channel (unsuccessful msgs)
	//3, boldMag //for timeouts
	//4, cyan - normal msgs

	for _, str := range *msgs {

		switch str.id {

		case 1:
			boldGreen := color.New(color.FgGreen, color.Bold)
			boldGreen.Printf(str.message)

		case 2:
			boldRed := color.New(color.FgRed, color.Bold)
			boldRed.Printf(str.message)

		case 3:
			boldMag := color.New(color.FgMagenta, color.Bold)
			boldMag.Printf(str.message)

		default:

			normalCyan := color.New(color.FgCyan)
			normalCyan.Printf(str.message)

		}

	}

	//Clear the buffer
	*msgs = nil

}

func fetchUrl(idx int, urlObj commondata.UrlObject, speedCh chan<- int, doneCh chan<- commondata.Message, msgCh chan<- string, outputFolder string) {

	r, err := http.NewRequest("GET", urlObj.Url, nil)
	if err != nil {
		msgCh <- fmt.Sprintf("Unable to make GET request: %s", urlObj.Url)
		return
	}

	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		msgSplit := strings.Split(err.Error(), ": ")
		if len(msgSplit) < 4 {
			msgCh <- fmt.Sprintf("%s [%s] [%s]", msgSplit[2], r.URL.Host, urlObj.Filename)
		} else {
			msgCh <- fmt.Sprintf("%s [%s] [%s]", msgSplit[3], r.URL.Host, urlObj.Filename)
		}
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		file, err := os.Create(filepath.Join(outputFolder, urlObj.Filename))
		if err != nil {
			msgCh <- fmt.Sprintf("Unbale to create file : %s/%s", outputFolder, urlObj.Filename)
			return
		}

		//size, err := io.Copy(file, res.Body)
		defer file.Close()

		var chunkSize int64 = 1024 * 1024 // 1 MB chunks
		var totalSize int64
		buffer := make([]byte, chunkSize)

		for {
			n, err := res.Body.Read(buffer)
			if err != nil && err != io.EOF {
				fmt.Println("\n Error reading response:", err)
				msgCh <- fmt.Sprintf("Unable to make read response: %s", urlObj.Url)
				return
			}

			if n == 0 {
				//completed
				break
			}

			n, writeErr := file.Write(buffer[:n])
			if writeErr != nil {
				fmt.Println("Error writing to file:", writeErr)
				msgCh <- fmt.Sprintf("Unbale to create file : %s/%s [%s]", outputFolder, urlObj.Filename, writeErr)
				return
			}

			//fmt.Printf("[%.3f MB]\n", float32(n)/(1024*1024))
			totalSize += int64(n)
			//fmt.Printf("Downloaded %d bytes [%.3f MB]\n", totalSize, float32(totalSize)/(1024*1024))
			speedCh <- n
		}

		doneCh <- commondata.Message{Idx: idx, Size: totalSize}

	} else {
		msgCh <- fmt.Sprintf("%s [%s]", res.Status, urlObj.Filename)
	}

}

func glogger(doneCh chan commondata.Message, msgCh chan string, quitCh chan<- bool, validUrls []commondata.UrlObject, outputFolder string, bar *progressbar.ProgressBar) {

	//For recieving per second Info
	speedCh := make(chan int, 50)

	color.Cyan("Fetching valid urls ...")
	//Setting up the recievers
	go receiver(doneCh, msgCh, quitCh, speedCh, bar, validUrls)

	for idx, urlObj := range validUrls {

		//fmt.Println(">> ", idx, urlObj.Filename, urlObj.Size)

		//Fetch the urls
		go fetchUrl(idx, urlObj, speedCh, doneCh, msgCh, outputFolder)

	}

}

func receiver(doneCh <-chan commondata.Message, msgCh <-chan string, quitCh chan<- bool, speedCh <-chan int, bar *progressbar.ProgressBar, validUrls []commondata.UrlObject) {
	/*	doneCh - bidirectional Channel
		msgCh - receive from channel
		quitCh <-send to channel
	*/

	//Initialize the counter to Zero
	reqCount = 0
	//Total numbe rof request to be fulfilled
	totalReqs = len(validUrls)

	//cumulative byte size
	var totalBytes int

	//Default timeout of 5 Mins (Assuming network speed is enough to wrap it up in 5 mins)
	defaultTimeoutSeconds := 120000000

	//TimeoutChannel
	timeout := time.After(time.Second * time.Duration(defaultTimeoutSeconds))

	//Create a new ticker for 1 second
	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	var tMessages []terminalMessage

	//That wild infi Loop
	for {

		select {

		//Save the size of bytes sent by reader Go routines
		case bps := <-speedCh:
			totalBytes += bps

		//check the doneChannel, if any goroutine have finished
		case msg := <-doneCh:

			//Fetch the filename and Size
			validUrls[msg.Idx].Size = msg.Size
			fileSize := validUrls[msg.Idx].Size
			fileName := validUrls[msg.Idx].Filename

			//Prepare Success Message
			var printMsg string = fmt.Sprintf("Downloaded: [%s] [size %.2f MB]", fileName, float32(fileSize)/(1024*1024))
			tMessages = append(tMessages, terminalMessage{
				id:      1,
				message: fmt.Sprintf("\r%-120s\n", printMsg),
			})

			reqCount = reqCount + 1
			//currentSize += fileSize

		//check msgChannel for any message
		case msg := <-msgCh:

			tMessages = append(tMessages, terminalMessage{
				id:      2,
				message: fmt.Sprintf("\rERROR: %-120s\n", msg),
			})

			reqCount = reqCount + 1

		case <-timeout:

			tMessages = append(tMessages, terminalMessage{
				id:      3,
				message: fmt.Sprintf("\rTimeout of %d seconds reached! Quitting\n", defaultTimeoutSeconds),
			})

			//Signalling the quit Channel
			quitCh <- true
			return

		//1sec
		case <-timer.C:

			currentCounter := reqCount
			//cSize := currentSize
			currentSize += int64(totalBytes)

			//Print Messages for Terminal
			flushToTerminal(&tMessages)

			//Update the Progres bar
			bar.Display(int64(currentCounter), currentSize, totalBytes)

			//reset the total bytes so far to zero
			totalBytes = 0

			if currentCounter == totalReqs {
				//Progress bar Ends
				bar.End()
				//Signalling the quit Channel
				quitCh <- true

				boldGreen := color.New(color.FgGreen, color.BgBlack, color.Bold)
				boldGreen.Printf("\nProcessed : %d url(s)", totalReqs)
				return
			}

		} //select

	} //for
}

func Run(totalSize *int64, targetUrls *commondata.TargetUrls, outputFolder *string) {

	//Setup channels
	//Channel to transfer send messages
	doneCh := make(chan commondata.Message)
	//Channel to transfer any other messages
	msgCh := make(chan string)
	//Bidirection Channel, used as send Quit message
	quitCh := make(chan bool)

	// Set up Progressbar Downloader Settings
	var N int = len(targetUrls.ValidUrls)
	var barSize int64 = 50
	var barSymbol string = "█"
	bar := progressbar.GetNewBar(int64(N), 0, barSymbol, barSize, 0, *totalSize)

	//Start Download Process
	glogger(doneCh, msgCh, quitCh, targetUrls.ValidUrls, *outputFolder, &bar)

	//Wait for the Quit Signal
	if true == <-quitCh {
		boldMag := color.New(color.FgMagenta).Add(color.Bold)
		boldMag.Printf("\nExiting ! ")
		return
	}
}
