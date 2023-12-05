package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/3l-d1abl0/goProgressBar"
	"github.com/fatih/color"
)

var pathToFile *string
var sliceUrls []string
var wg sync.WaitGroup
var nOReq int
var countReq int
var mu sync.Mutex
var previousMessageType int

func checkNilErr(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchUrl(target_url string, doneCh chan<- string, msgCh chan<- string) {

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

		file, err := os.Create(filepath.Join("output", fileName))
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

func slogger(doneCh chan<- string, msgCh chan<- string) {

	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err := os.Mkdir("output", os.ModePerm); err != nil {
			panic(err)
		}
	}

	//Set waitgroup size
	//wg.Add(nOReq)

	color.Cyan("Fetching urls ...")
	for _, target_url := range sliceUrls {
		go fetchUrl(target_url, doneCh, msgCh)
		//go fetchUrl1(target_url, doneCh, msgCh)
		//fmt.Println(target_url)
	}

	//wg.Wait()

}

func receiver(doneCh <-chan string, msgCh <-chan string, quitCh chan<- bool, bar *goProgressBar.ProgressBar) {
	/*	doneCh - bidirectional Channel
		msgCh - receive from channel
		quitCh <-send to channel
	*/

	nOReq = len(sliceUrls)
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

func init() {
	fmt.Println("init")

	pathToFile = flag.String("file", "", "file path to urls")
	flag.Parse()

	if *pathToFile == "" {
		os.Exit(1)
	}

}

func main() {

	color.Cyan("Path to file: %s", *pathToFile)
	color.Cyan("Parsing file ...")

	//Read line by line
	file, err := os.Open(*pathToFile)
	checkNilErr(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	boldYellow := color.New(color.FgYellow, color.Bold)
	for scanner.Scan() {

		target_url := scanner.Text()
		_, err := url.ParseRequestURI(target_url)
		if err != nil {
			boldYellow.Printf("Skipping: %s\n", target_url)
			//fmt.Println(msg)
		} else {
			color.Green("OK: %s", target_url)
			sliceUrls = append(sliceUrls, target_url)
		}
	}
	color.Cyan("Scanner ends !")

	//Channel to transfer send messages
	doneCh := make(chan string)
	//Channel to transfer other messages
	msgCh := make(chan string)
	//Bidirection Channel, used as send Quit message
	quitCh := make(chan bool)

	//Setting up Downloader Settings
	var N int = len(sliceUrls)
	var barSize int64 = 70
	var barSymbol string = "#"
	bar := goProgressBar.GetNewBar(int64(N), 0, barSymbol, barSize)

	//boldGreen := color.New(color.FgGreen, color.Bold)
	//boldGreen.Printf(" %d/%d Started ...\n", countReq, nOReq)

	//Setting up the recievers
	//Start Downloading
	slogger(doneCh, msgCh)

	go receiver(doneCh, msgCh, quitCh, &bar)

	//Wait for the Quit Signal
	println(<-quitCh)
}
