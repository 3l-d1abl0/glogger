package main

import ("fmt"
		"flag"
		"os"
)

var path *string

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
	
}