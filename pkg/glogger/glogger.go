package glogger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/3l-d1abl0/goProgressBar"
	"github.com/fatih/color"
)

var nOReq int
var countReq int
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
	//fmt.Println(res)
	if err != nil {
		msgCh <- err.Error()
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		//downloadSize, _ := strconv.Atoi(res.Header.Get("Content-Length"))
		//fmt.Println("Size: ", int64(downloadSize)/(1024*1024))
		//color.Cyan("Size: %.2f MB (%s)", float64(downloadSize)/(1024*1024), target_url)

		file, err := os.Create(filepath.Join(outputFolder, fileName))
		if err != nil {
			fmt.Println("ERR3")
			panic(err)
		}

		size, err := io.Copy(file, res.Body)
		defer file.Close()

		doneCh <- fmt.Sprintf("Downloaded: %s [size %.2f MB]", fileName, float32(size)/(1024*1024))
		//color.Green()

	} else {
		//fmt.Println(res)
		msgCh <- fmt.Sprintf("%s [%s]", res.Status, fileName)
	}

	//wg.Done()
}

func Glogger(doneCh chan<- string, msgCh chan<- string, validUrls []string, outputFolder string) {

	/*if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err := os.Mkdir("output", os.ModePerm); err != nil {
			panic(err)
		}
	}*/

	//Set waitgroup size
	//wg.Add(nOReq)

	color.Cyan("Fetching urls ...")
	for _, target_url := range validUrls {
		go fetchUrl(target_url, doneCh, msgCh, outputFolder)
		//go fetchUrl1(target_url, doneCh, msgCh)
		//fmt.Println(target_url)
	}

	//wg.Wait()

}

func Receiver(doneCh <-chan string, msgCh <-chan string, quitCh chan<- bool, bar *goProgressBar.ProgressBar, ValidUrlCount int) {
	/*	doneCh - bidirectional Channel
		msgCh - receive from channel
		quitCh <-send to channel
	*/

	nOReq = ValidUrlCount
	//Initialize the counter to Zero
	countReq = 0

	//Default timeout of 5 Mins
	defaultTimeoutSeconds := 600
	timeout := time.After(time.Second * time.Duration(defaultTimeoutSeconds))

	for {

		select {

		case msg := <-doneCh:
			/* If previous was from progressbar, override it
			 */
			boldGreen := color.New(color.FgGreen, color.Bold)

			if previousMessageType == 2 {
				boldGreen.Printf("\r%-100s", msg)
			} else {
				boldGreen.Printf("\n\r%-100s", msg)
			}

			//update the message Type
			previousMessageType = 1

			//Processed - increase the request counter
			mu.Lock()
			countReq = countReq + 1
			//Mutex Unlock
			mu.Unlock()

		case msg := <-msgCh:

			boldRed := color.New(color.FgRed, color.Bold)

			if previousMessageType == 2 {
				boldRed.Printf("\rERROR: %-100s", msg)
			} else {
				boldRed.Printf("\n\rERROR: %-100s", msg)
			}

			//update the message Type
			previousMessageType = 1

			//Processed - increase the request counter
			mu.Lock()
			countReq = countReq + 1
			//Mutex Unlock
			mu.Unlock()

		case <-timeout:
			boldRed := color.New(color.FgRed, color.Bold)
			boldRed.Printf("\n\rTimeout of %d seconds reached! Quitting", defaultTimeoutSeconds)
			quitCh <- true
			return

		default:

			mu.Lock()
			currentCounter := countReq
			mu.Unlock()
			//fmt.Println(currentCounter, nOReq)
			//boldGreen := color.New(color.FgGreen, color.Bold)

			if previousMessageType != 2 {
				fmt.Println()
			}
			//update the message Type
			previousMessageType = 2

			bar.Display(int64(currentCounter))
			/*if currentCounter < nOReq {
				boldGreen.Printf(" %d/%d Completed ...\n", currentCounter, nOReq)
				currentCounter = countReq
				time.Sleep(time.Second * 2)
			} else*/
			if currentCounter == nOReq {
				bar.End()
				quitCh <- true
				fmt.Printf("\nDone Quitting!")
				return
			}

			time.Sleep(time.Second * 1)
		}
	}
}
