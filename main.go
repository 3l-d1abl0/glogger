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

	"github.com/fatih/color"
)

var pathToFile *string
var sliceUrls []string
var wg sync.WaitGroup
var nOReq int
var countReq int
var mu sync.Mutex

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
	if err != nil {
		fmt.Println("ERR2")
		msgCh <- "ERR2"
		panic(err)
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

		doneCh <- fmt.Sprintf("Downloaded: %s [size %.2f MB]\n", fileName, float32(size)/(1024*1024))
		//color.Green()

	} else {
		//fmt.Println(res)
		msgCh <- fmt.Sprintf("%s [%s]\n", res.Status, fileName)
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

func receiver(doneCh <-chan string, msgCh <-chan string, quitCh chan<- bool) {
	/*	doneCh - bidirectional Channel
		msgCh - receive from channel
		quitCh <-send to channel
	*/

	nOReq = len(sliceUrls)
	//Initialize the counter to Zero
	countReq = 0

	timeout := time.After(time.Second * 20)

	for {

		select {

		case msg := <-doneCh:
			println(": " + msg)
			//Processed - increase the request counter
			mu.Lock()
			countReq = countReq + 1
			//Mutex Unlock
			mu.Unlock()

		case msg := <-msgCh:
			println("::" + msg)
			//Processed - increase the request counter
			mu.Lock()
			countReq = countReq + 1
			//Mutex Unlock
			mu.Unlock()

		case <-timeout:
			println("Nothing received in 20 seconds. Exiting")
			quitCh <- true
			return

		default:

			mu.Lock()
			currentCounter := countReq
			mu.Unlock()
			fmt.Println(currentCounter, nOReq)
			boldGreen := color.New(color.FgGreen, color.Bold)
			if currentCounter < nOReq {
				boldGreen.Printf(" %d/%d Completed ...\n", currentCounter, nOReq)
				currentCounter = countReq
				time.Sleep(time.Second * 2)
			} else if currentCounter == nOReq {
				quitCh <- true
				fmt.Println("Done Quitting")
				return
			}

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

	slogger(doneCh, msgCh)
	boldGreen := color.New(color.FgGreen, color.Bold)
	boldGreen.Printf(" %d/%d Started ...\n", countReq, nOReq)

	//Setting up the recievers
	go receiver(doneCh, msgCh, quitCh)

	println("Waiting")
	println(<-quitCh)
}
