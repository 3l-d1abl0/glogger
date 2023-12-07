package readUrls

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

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

// fetches the size of a url
func getSize(url string) (int64, error) {

	response, err := http.Head(url)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer response.Body.Close()

	// Check if the response status code is in the 2xx range
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return 0, fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	size := response.ContentLength
	return size, nil
}

func ReadUrls(pathToFile *string, targetUrls *commondata.TargetUrls) (int64, error) {

	//Cyan for normal print
	cCy := color.New(color.FgCyan)
	cRe := color.New(color.FgRed)
	cYel := color.New(color.FgYellow, color.Bold)

	if targetUrls.ValidUrls == nil {
		targetUrls.ValidUrls = make([]string, 0)
	}

	if targetUrls.InvalidUrls == nil {
		targetUrls.InvalidUrls = make([]string, 0)
	}

	cCy.Printf("Parsing file : %s\n", *pathToFile)

	//Read line by line
	file, err := os.Open(*pathToFile)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var totalSize int64 = 0
	for scanner.Scan() {

		target_url := strings.TrimSpace(scanner.Text())
		_, err := url.ParseRequestURI(target_url)
		if err != nil {

			fmt.Printf("%s: %s\n", cCy.SprintFunc()("Skipping"), cYel.SprintFunc()(target_url))
			targetUrls.InvalidUrls = append(targetUrls.InvalidUrls, target_url)
		} else {
			targetUrls.ValidUrls = append(targetUrls.ValidUrls, target_url)

			fileName := path.Base(target_url)
			// Fetch the size of the URL
			size, err := getSize(target_url)

			//Update the total Size
			totalSize += size
			if err != nil {
				cCy.Printf("INFO: [%s]", fileName)
				cRe.Printf(" [cannot fetch size]\n")
			} else {
				cCy.Printf("INFO :[%s] [size %.2f MB]\n", fileName, float32(size)/(1024*1024))
			}
		}
	} //for scanner

	return totalSize, nil
}
