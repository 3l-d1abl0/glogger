package main

import ("fmt"
		"flag"
		"os"
		"bufio"
		"net/url"
		"sync"
		"path/filepath"
		"path"
		"net/http"
		"io"
		"strconv"
		"time"
		"github.com/fatih/color"
)

var pathToFile *string
var sliceUrls []string
var wg sync.WaitGroup
var nOReq int
var countReq int
var mu sync.Mutex

func checkNilErr(err error){
	if err != nil{
		panic(err)
	}
}

func fetchUrl(target_url string){

	r, err := http.NewRequest("GET", target_url, nil)
	if err != nil {
		fmt.Println("ERR1")
		panic(err)
	}

	fileName := path.Base(r.URL.Path)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("ERR2")
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		downloadSize, _ := strconv.Atoi(res.Header.Get("Content-Length"))
		//fmt.Println("Size: ", int64(downloadSize)/(1024*1024))
		color.Cyan("Size: %.2f MB (%s)", float64(downloadSize)/(1024*1024), target_url)

		file, err := os.Create(filepath.Join("output", fileName))
		if err != nil {
			fmt.Println("ERR3")
			panic(err)
		}

		size, err := io.Copy(file, res.Body)
    	defer file.Close()

		color.Green("Downloaded: %s , size %.2f MB\n", fileName, float32(size)/(1024*1024))
	}else{
		fmt.Println(res)
	}

	//Processed - increase the request counter
	mu.Lock()
	countReq = countReq+1
	//Mutex Unlock
	mu.Unlock()

	wg.Done()
}

func slogger(){
	
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err := os.Mkdir("output", os.ModePerm); err != nil {
			panic(err)
		}
	}

	nOReq = len(sliceUrls)
	//Set waitgroup size
	wg.Add(nOReq)

	//Initialize the counter to Zero
	countReq =0


	color.Cyan("Fetching urls ...")
	for _, target_url := range sliceUrls {
        go fetchUrl(target_url)
    }

	mu.Lock()
	currentCounter := countReq
	mu.Unlock()

	boldGreen := color.New(color.FgGreen, color.Bold)
	for currentCounter < nOReq{
		boldGreen.Printf(" %d/%d Completed ...\n", currentCounter, nOReq)
		currentCounter = countReq
		time.Sleep(time.Second*2)
	}
	wg.Wait()

}

func init(){
	fmt.Println("init")

	pathToFile = flag.String("file", "", "file path to urls")
	flag.Parse()
	
	if *pathToFile == "" {
		os.Exit(1)
	}
	
}

func main(){

	color.Cyan("Path to file: %s", *pathToFile)
	color.Cyan("Parsing file ...")

	//Read line by line
	file, err := os.Open(*pathToFile)
	checkNilErr(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	boldYellow := color.New(color.FgYellow, color.Bold)
	for scanner.Scan(){

		target_url := scanner.Text()
		_ , err := url.ParseRequestURI(target_url)
		if(err!=nil){
			boldYellow.Printf("Skipping: %s\n", target_url)
			//fmt.Println(msg)
		}else{
			color.Green("OK: %s", target_url)
			sliceUrls = append(sliceUrls, target_url)
		}
	}
	color.Cyan("Scanner ends !")

	slogger()
	boldGreen := color.New(color.FgGreen, color.Bold)
	boldGreen.Printf(" %d/%d Completed ...\n", countReq, nOReq)
	
}