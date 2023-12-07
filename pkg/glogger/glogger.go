package glogger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/3l-d1abl0/goProgressBar"
	"github.com/fatih/color"
)

var totalReqs int
var reqCount int
var mu sync.Mutex
var previousMessageType int

func fetchUrl(target_url string, doneCh chan<- string, msgCh chan<- string, outputFolder string) {

	r, err := http.NewRequest("GET", target_url, nil)
	if err != nil {
		fmt.Println("ERR1")
		msgCh <- "ERR1"
		panic(err)
	}

	fileName := path.Base(r.URL.Path)

	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		msgSplit := strings.Split(err.Error(), ": ")
		msgCh <- fmt.Sprintf("%s [%s] [%s]", msgSplit[3], r.URL.Host, fileName)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		file, err := os.Create(filepath.Join(outputFolder, fileName))
		if err != nil {
			fmt.Println("ERR3")
			panic(err)
		}

		size, err := io.Copy(file, res.Body)
		defer file.Close()

		doneCh <- fmt.Sprintf("Downloaded: [%s] [size %.2f MB]", fileName, float32(size)/(1024*1024))

	} else {
		msgCh <- fmt.Sprintf("%s [%s]", res.Status, fileName)
	}

}

func Glogger(doneCh chan string, msgCh chan string, quitCh chan<- bool, validUrls []string, outputFolder string, bar *goProgressBar.ProgressBar) {

	color.Cyan("Fetching valid urls ...")

	//Setting up the recievers
	go receiver(doneCh, msgCh, quitCh, bar, len(validUrls))

	for _, target_url := range validUrls {

		//Fetch the urls
		go fetchUrl(target_url, doneCh, msgCh, outputFolder)

	}

}

func receiver(doneCh <-chan string, msgCh <-chan string, quitCh chan<- bool, bar *goProgressBar.ProgressBar, ValidUrlCount int) {
	/*	doneCh - bidirectional Channel
		msgCh - receive from channel
		quitCh <-send to channel
	*/

	//Initialize the counter to Zero
	reqCount = 0
	//Total numbe rof request to be fulfilled
	totalReqs = ValidUrlCount

	//Default timeout of 5 Mins (Assuming network speed is enough to wrap it up in 5 mins)
	defaultTimeoutSeconds := 600
	timeout := time.After(time.Second * time.Duration(defaultTimeoutSeconds))

	//Color Convention
	boldGreen := color.New(color.FgGreen, color.Bold) //for done channel
	boldRed := color.New(color.FgRed, color.Bold)     //for msg channel (unsuccessful msgs)
	boldMag := color.New(color.FgMagenta, color.Bold) //for timeouts

	//That wild infi Loop
	for {

		select {

		//check the doneChannel, if any goroutine have finished
		case msg := <-doneCh:
			/* If previous was from progressbar, override it
			pbar has message type 2, rest 1
			*/
			if previousMessageType == 2 {
				boldGreen.Printf("\r%-100s", msg)
			} else {
				boldGreen.Printf("\n\r%-100s", msg)
			}

			//update the message Type
			previousMessageType = 1

			//Processed a url - increase the request counter
			mu.Lock() //Mutex Lock
			reqCount = reqCount + 1
			mu.Unlock() //Mutex Unlock

		//check msgChannel for any message
		case msg := <-msgCh:

			if previousMessageType == 2 {
				boldRed.Printf("\rERROR: %-100s", msg)
			} else {
				boldRed.Printf("\n\rERROR: %-100s", msg)
			}

			//update the message Type
			previousMessageType = 1

			//Processed - increase the request counter
			mu.Lock()
			reqCount = reqCount + 1
			mu.Unlock()

		case <-timeout:
			boldMag.Printf("\n\rTimeout of %d seconds reached! Quitting", defaultTimeoutSeconds)

			//Signalling the quit Channel
			quitCh <- true
			return

		default:

			mu.Lock()
			currentCounter := reqCount
			mu.Unlock()

			//The previous print message was from other channels
			if previousMessageType != 2 {
				fmt.Println()
			}
			//update the message Type
			previousMessageType = 2

			//Update the Progres bar
			bar.Display(int64(currentCounter))
			if currentCounter == totalReqs {
				//Progress bar Ends
				bar.End()
				//Signalling the quit Channel
				quitCh <- true

				//boldGreen.Add(color.BgBlack)
				boldGreen.Printf("\nProcessed : %d url(s)", totalReqs)
				return
			}

			time.Sleep(time.Second * 1)
		}
	}
}
