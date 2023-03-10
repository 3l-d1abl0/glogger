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
)

var pathToFile *string
var sliceUrls []string
var wg sync.WaitGroup

func checkNilErr(err error){
	if err != nil{
		panic(err)
	}
}

func fetchUrl(target_url string){

	r, err := http.NewRequest("GET", target_url, nil)
	if err != nil {
		panic(err)
	}

	fileName := path.Base(r.URL.Path)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		file, err := os.Create(filepath.Join("output", fileName))
		if err != nil {
			panic(err)
		}

		size, err := io.Copy(file, res.Body)
    	defer file.Close()

		fmt.Printf("Downloaded: %s , size %.2f MB\n", fileName, float32(size)/(1024*1024))
	}

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

	pathToFile = flag.String("file", "", "file path to urls")
	flag.Parse()
	
	if *pathToFile == "" {
		os.Exit(1)
	}
	
}

func main(){

	fmt.Println("path to file ", *pathToFile)
	fmt.Println("main")

	//Read line by line
	file, err := os.Open(*pathToFile)
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