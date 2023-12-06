package readUrls

import (
	"bufio"
	"net/url"
	"os"

	"glogger/pkg/commondata"

	"github.com/fatih/color"
)

func checkNilErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Urls struct {
	flag1      string
	filePath   string
	flag2      string
	folderPath string
	expected   *string
	comments   string
}

func ReadUrls(pathToFile *string, targetUrls *commondata.TargetUrls) {

	if targetUrls.ValidUrls == nil {
		targetUrls.ValidUrls = make([]string, 0)
	}

	if targetUrls.InvalidUrls == nil {
		targetUrls.InvalidUrls = make([]string, 0)
	}

	color.Cyan("Path to file: %s", *pathToFile)
	color.Cyan("Parsing file ...")

	//Read line by line
	file, err := os.Open(*pathToFile)
	checkNilErr(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	//var sliceUrls []string

	boldYellow := color.New(color.FgYellow, color.Bold)
	for scanner.Scan() {

		target_url := scanner.Text()
		_, err := url.ParseRequestURI(target_url)
		if err != nil {
			boldYellow.Printf("Skipping: %s\n", target_url)
			targetUrls.InvalidUrls = append(targetUrls.InvalidUrls, target_url)
		} else {
			//color.Green("OK: %s", target_url)
			targetUrls.ValidUrls = append(targetUrls.ValidUrls, target_url)
		}
	}
	color.Cyan("Scanner ends !")
}
