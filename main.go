package main

import (
	"fmt"
	"glogger/cmd"
)

func main() {

	//filePath, outputFolder, error := cmd.GetArgs()
	_, _, err := cmd.GetArgs()

	if err != nil {
		//fmt.Println(error)
		fmt.Printf("MAIN: Unable to Read input: (%s) \n", error)
	}
}
