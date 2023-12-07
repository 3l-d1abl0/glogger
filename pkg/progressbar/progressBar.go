package progressbar

import (
	"fmt"
	"math"
	"strings"
)

type count struct {
	current int64 //current count
	total   int64 //total count
}

type ProgressBar struct {
	status1        count     //status1
	status2        count     //status2
	percentage     float64   //progress % based on status1
	char           string    //character/Symbol for ProgressBar
	totalBarSize   int64     //total no of character to depict 100%
	currentBarSize int64     //current number of Characters to represent bars
	spinner        [4]string //characters to show a spinner
	pulse          int64     //an increasing number to be used by spinner
}

func GetNewBar(totalElements int64, currentElements int64, symbol string, totalBarSize int64, currentCount int64, totalCount int64) ProgressBar {

	CurrentBar := int64(math.Ceil(float64(currentElements/totalElements) * 50))
	return ProgressBar{
		status1:        count{currentElements, totalElements},
		status2:        count{currentCount, totalCount},
		char:           symbol,
		percentage:     float64(currentElements/totalElements) * 100,
		totalBarSize:   totalBarSize,
		currentBarSize: CurrentBar,
		spinner:        [4]string{"|", "/", "—", "\\"},
		pulse:          0,
	}

}

func (bar *ProgressBar) Display(currentStatus1 int64, currentStatus2 int64) {

	//update the current status
	bar.status1.current = currentStatus1

	//calculate the new current percentage
	bar.percentage = (float64(bar.status1.current) * 100.0) / float64(bar.status1.total)

	//calculate the number of Symbols to display
	bar.currentBarSize = int64(bar.percentage * float64(bar.totalBarSize) / 100)

	//Generate the ProgressBar string
	currentBar := strings.Repeat(bar.char, int(bar.currentBarSize))

	//Generate the Format String
	fmtString := fmt.Sprintf("\r{%%s} [%%-%ds] %%3.2f%%%% [%%2.2fMB/%%2.2fMB] %%5d/%%d", bar.totalBarSize)
	//fmtString := fmt.Sprintf("\r[%%-%ds] %%3.2f%%%% %%8d/%%d {%%s}", bar.totalBarSize)
	bar.pulse = (bar.pulse + 1) % 4

	//fmt.Println(fmtString)
	//[%-70s] %3.2f%% %2.2fMB/%2.2fMB %8d/%d {%s}
	var currentVal float32 = float32(currentStatus2) / (1024 * 1024)
	var totalVal float32 = float32(bar.status2.total) / (1024 * 1024)
	fmt.Printf(fmtString, bar.spinner[bar.pulse], currentBar, bar.percentage, currentVal, totalVal, bar.status1.current, bar.status1.total)
}

func (bar *ProgressBar) End() {
	fmt.Println()
}
