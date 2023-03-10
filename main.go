package main

import ("fmt"
		"flag"
		"os"
		"bufio"
		"net/url"
)

var path *string
var sliceUrls []string

func checkNilErr(err error){
	if err != nil{
		panic(err)
	}
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
}