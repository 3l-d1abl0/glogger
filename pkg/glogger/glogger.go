package glogger

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

var totalReqs int
var reqCount int
var mu sync.Mutex
var currentSize int64
var previousMessageType int

func fetchUrl(idx int, urlObj commondata.UrlObject, doneCh chan<- commondata.Message, msgCh chan<- string, outputFolder string) {

	//fmt.Println("<< ", idx, urlObj.Filename, urlObj.Size)

	r, err := http.NewRequest("GET", urlObj.Url, nil)
	if err != nil {
		fmt.Println("ERR1")
		msgCh <- "ERR1"
		panic(err)
	}

	//fileName := path.Base(r.URL.Path)

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

		//contentLength := res.Header.Get("Content-Length")
		//	fmt.Println("Content Length:", contentLength, urlObj.Filename)

		file, err := os.Create(filepath.Join(outputFolder, urlObj.Filename))
		if err != nil {
			fmt.Println("ERR3")
			panic(err)
		}

		size, err := io.Copy(file, res.Body)
		defer file.Close()

		doneCh <- commondata.Message{Idx: idx, Size: size}
		//doneCh <- fmt.Sprintf("Downloaded: [%s] [size %.2f MB]", fileName, float32(size)/(1024*1024))

	} else {
		msgCh <- fmt.Sprintf("%s [%s]", res.Status, urlObj.Filename)
	}

}

func Glogger(doneCh chan commondata.Message, msgCh chan string, quitCh chan<- bool, validUrls []commondata.UrlObject, outputFolder string, bar *progressbar.ProgressBar) {

	color.Cyan("Fetching valid urls ...")

	//Setting up the recievers
	go receiver(doneCh, msgCh, quitCh, bar, validUrls)

	for idx, urlObj := range validUrls {

		//fmt.Println(">> ", idx, urlObj.Filename, urlObj.Size)

		//Fetch the urls
		go fetchUrl(idx, urlObj, doneCh, msgCh, outputFolder)

	}

}

func receiver(doneCh <-chan commondata.Message, msgCh <-chan string, quitCh chan<- bool, bar *progressbar.ProgressBar, validUrls []commondata.UrlObject) {
	/*	doneCh - bidirectional Channel
		msgCh - receive from channel
		quitCh <-send to channel
	*/

	//Initialize the counter to Zero
	reqCount = 0
	//Total numbe rof request to be fulfilled
	totalReqs = len(validUrls)

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

			//fmt.Println(msg)
			validUrls[msg.Idx].Size = msg.Size
			fileSize := validUrls[msg.Idx].Size
			fileName := validUrls[msg.Idx].Filename

			var printMsg string = fmt.Sprintf("Downloaded: [%s] [size %.2f MB]", fileName, float32(fileSize)/(1024*1024))
			if previousMessageType == 2 {
				boldGreen.Printf("\r%-120s", printMsg)
			} else {
				boldGreen.Printf("\n\r%-120s", printMsg)
			}

			//update the message Type
			previousMessageType = 1

			//Processed a url - increase the request counter and the downloaded Size
			mu.Lock() //Mutex Lock
			reqCount = reqCount + 1
			currentSize += fileSize
			mu.Unlock() //Mutex Unlock

		//check msgChannel for any message
		case msg := <-msgCh:

			if previousMessageType == 2 {
				boldRed.Printf("\rERROR: %-120s", msg)
			} else {
				boldRed.Printf("\n\rERROR: %-120s", msg)
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
			cSize := currentSize
			mu.Unlock()

			//The previous print message was from other channels
			if previousMessageType != 2 {
				fmt.Println()
			}
			//update the message Type
			previousMessageType = 2

			//Update the Progres bar
			bar.Display(int64(currentCounter), cSize)
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
