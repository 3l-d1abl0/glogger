package readUrls

import (
	"bufio"
	"fmt"
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

	//Cyan for normal print
	cCy := color.New(color.FgCyan)

	if targetUrls.ValidUrls == nil {
		targetUrls.ValidUrls = make([]string, 0)
	}

	if targetUrls.InvalidUrls == nil {
		targetUrls.InvalidUrls = make([]string, 0)
	}

	cCy.Printf("Parsing file : %s\n", *pathToFile)

	//Read line by line
	file, err := os.Open(*pathToFile)
	checkNilErr(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	cYel := color.New(color.FgYellow, color.Bold)
	for scanner.Scan() {

		target_url := scanner.Text()
		_, err := url.ParseRequestURI(target_url)
		if err != nil {

			fmt.Printf("%s: %s\n", cCy.SprintFunc()("Skipping"), cYel.SprintFunc()(target_url))
			targetUrls.InvalidUrls = append(targetUrls.InvalidUrls, target_url)
		} else {
			targetUrls.ValidUrls = append(targetUrls.ValidUrls, target_url)
		}
	}
}
