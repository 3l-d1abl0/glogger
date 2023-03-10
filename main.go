package main

import ("fmt"
		"flag"
		"os"
		"bufio"
		"net/url"
		"sync"
		"time"
)

var path *string
var sliceUrls []string
var wg sync.WaitGroup

func checkNilErr(err error){
	if err != nil{
		panic(err)
	}
}

func fetchUrl(target_url string){
	
	fmt.Println("Trying to fetch: ...", target_url)
	time.Sleep(5 * time.Second)
	fmt.Println("Fetched !")

	wg.Done()
}

func slogger(){
	
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err := os.Mkdir("output", os.ModePerm); err != nil {
			panic(err)
		}
	}

	urlSz := len(sliceUrls)
	//Set waitgroup size
	wg.Add(urlSz)

	for _, target_url := range sliceUrls {
        go fetchUrl(target_url)
    }

	wg.Wait()

}

func init(){
	fmt.Println("init")

	path = flag.String("file", "", "file path to urls")
	flag.Parse()
	
	if *path == "" {
		os.Exit(1)
	}
	
}

func main(){

	fmt.Println("path to file ", *path)
	fmt.Println("main")

	//Read line by line
	file, err := os.Open(*path)
	checkNilErr(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan(){

		target_url := scanner.Text()
		_ , err := url.ParseRequestURI(target_url)
		if(err!=nil){
			fmt.Println("Skipping: ", target_url)
			//fmt.Println(msg)
		}else{
			fmt.Println("OK: ", target_url)
			sliceUrls = append(sliceUrls, target_url)
		}
	}
	fmt.Println("Scanner end")

	slogger()

	fmt.Println("END")
}