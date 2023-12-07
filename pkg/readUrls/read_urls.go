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

// fetches the size of a url
func getSize(url string) (int64, error) {

	response, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	// Check if the response status code is in the 2xx range
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		//Head "https://static.git.ir/sub/cdn15.git.ir/02/Udemy%20Rust%20Programming%20The%20Complete%20Guide-git.ir/088-Channels-git.ir.fa.srt": dial tcp: lookup static.git.ir: no such host

		urlSplit := strings.Split(url, "/")
		return 0, fmt.Errorf("GET %s: dail tcp: %s: %d", url, urlSplit[2], response.StatusCode)
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
		targetUrls.ValidUrls = make([]commondata.UrlObject, 0)
	}

	if targetUrls.InvalidUrls == nil {
		targetUrls.InvalidUrls = make([]commondata.UrlObject, 0)
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
		fileName := path.Base(target_url)

		_, err := url.ParseRequestURI(target_url)
		if err != nil {

			fmt.Printf("%s: %s\n", cCy.SprintFunc()("Skipping"), cYel.SprintFunc()(target_url))

			newUrlObject := commondata.UrlObject{
				Url:      target_url,
				Filename: fileName,
				Size:     0,
			}

			targetUrls.InvalidUrls = append(targetUrls.InvalidUrls, newUrlObject)
		} else {

			// Fetch the size of the URL
			size, err := getSize(target_url)

			//new UrlObject
			var newUrlObject commondata.UrlObject

			//Update the total Size
			totalSize += size
			if err != nil {

				msgSplit := strings.Split(err.Error(), ": ")
				cCy.Printf("INFO: [%s] [%s] [%s]", msgSplit[2], fileName, msgSplit[3])
				cRe.Printf(" [cannot fetch size]\n")

				newUrlObject = commondata.UrlObject{
					Url:      target_url,
					Filename: fileName,
					Size:     0,
				}
			} else {
				cCy.Printf("INFO :[%s] [size %.2f MB]\n", fileName, float32(size)/(1024*1024))
				newUrlObject = commondata.UrlObject{
					Url:      target_url,
					Filename: fileName,
					Size:     size,
				}
			}

			targetUrls.ValidUrls = append(targetUrls.ValidUrls, newUrlObject)
		}
	} //for scanner

	return totalSize, nil
}
